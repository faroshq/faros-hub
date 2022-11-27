package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	farosclientset "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/cluster"
	edgeinformers "github.com/faroshq/faros-hub/pkg/client/informers/externalversions/edge/v1alpha1"
	edgelisters "github.com/faroshq/faros-hub/pkg/client/listers/edge/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	"github.com/kcp-dev/client-go/kubernetes"
	"github.com/kcp-dev/kcp/pkg/logging"
	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	controllerName = "faros-agent"
	finalizerName  = "agents.edge.faros.sh/finalizer"
)

// Controller watches Faros agents and configures them
type Controller struct {
	config *config.ControllerConfig

	queue workqueue.RateLimitingInterface

	coreClientSet  kubernetes.ClusterInterface
	farosClientSet farosclientset.ClusterInterface
	agentIndexer   cache.Indexer
	agentLister    edgelisters.AgentClusterLister
}

func NewController(
	config *config.ControllerConfig,
	coreClientSet kubernetes.ClusterInterface,
	farosClientSet farosclientset.ClusterInterface,
	agentInformer edgeinformers.AgentClusterInformer,
) (*Controller, error) {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName)

	c := &Controller{
		config:         config,
		queue:          queue,
		coreClientSet:  coreClientSet,
		farosClientSet: farosClientSet,
		agentIndexer:   agentInformer.Informer().GetIndexer(),
		agentLister:    agentInformer.Lister(),
	}

	agentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.enqueue(obj) },
		UpdateFunc: func(_, obj interface{}) { c.enqueue(obj) },
		DeleteFunc: func(obj interface{}) { c.enqueue(obj) },
	})

	return c, nil
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := kcpcache.MetaClusterNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}

	logger := logging.WithQueueKey(logging.WithReconciler(klog.Background(), controllerName), key)
	logger.V(2).Info("queueing Agent")
	c.queue.Add(key)
}

func (c *Controller) Start(ctx context.Context, numThreads int) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	logger := logging.WithReconciler(klog.FromContext(ctx), controllerName)
	ctx = klog.NewContext(ctx, logger)
	logger.Info("Starting controller")
	defer logger.Info("Shutting down controller")

	for i := 0; i < numThreads; i++ {
		go wait.Until(func() { c.startWorker(ctx) }, time.Second, ctx.Done())
	}

	<-ctx.Done()
}

func (c *Controller) startWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	// Wait until there is a new item in the working queue
	k, quit := c.queue.Get()
	if quit {
		return false
	}
	key := k.(string)

	logger := logging.WithQueueKey(klog.FromContext(ctx), key)
	ctx = klog.NewContext(ctx, logger)
	logger.V(1).Info("processing key")

	// No matter what, tell the queue we're done with this key, to unblock
	// other workers.
	defer c.queue.Done(key)

	if requeue, err := c.process(ctx, key); err != nil {
		runtime.HandleError(fmt.Errorf("%q controller failed to sync %q, err: %w", controllerName, key, err))
		c.queue.AddRateLimited(key)
		return true
	} else if requeue {
		// only requeue if we didn't error, but we still want to requeue
		c.queue.Add(key)
	}
	c.queue.Forget(key)
	return true
}

func (c *Controller) process(ctx context.Context, key string) (bool, error) {
	logger := klog.FromContext(ctx)

	cluster, namespace, name, err := kcpcache.SplitMetaClusterNamespaceKey(key)
	if err != nil {
		runtime.HandleError(err)
		return false, nil
	}

	obj, err := c.agentLister.Cluster(cluster).Agents(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil // object deleted before we handled it
		}
		return false, err
	}

	old := obj
	obj = obj.DeepCopy()

	ctx = klog.NewContext(ctx, logger)

	var errs []error
	requeue, err := c.reconcile(ctx, cluster, obj)
	if err != nil {
		errs = append(errs, err)
	}

	// Regardless of whether reconcile returned an error or not, always try to patch status if needed. Return the
	// reconciliation error at the end.

	// If the object being reconciled changed as a result, update it.
	if err := c.patchIfNeeded(ctx, old, obj); err != nil {
		errs = append(errs, err)
	}

	return requeue, utilerrors.NewAggregate(errs)
}

func (c *Controller) patchIfNeeded(ctx context.Context, old, obj *edgev1alpha1.Agent) error {
	specOrObjectMetaChanged := !equality.Semantic.DeepEqual(old.Spec, obj.Spec) || !equality.Semantic.DeepEqual(old.ObjectMeta, obj.ObjectMeta)
	statusChanged := !equality.Semantic.DeepEqual(old.Status, obj.Status)

	if specOrObjectMetaChanged && statusChanged {
		panic("Programmer error: spec and status changed in same reconcile iteration")
	}

	if !specOrObjectMetaChanged && !statusChanged {
		return nil
	}

	clusterAgentForPatch := func(agent *edgev1alpha1.Agent) edgev1alpha1.Agent {
		var ret edgev1alpha1.Agent
		if specOrObjectMetaChanged {
			ret.ObjectMeta = agent.ObjectMeta
			ret.Spec = agent.Spec
		} else {
			ret.Status = agent.Status
		}
		return ret
	}

	clusterName := logicalcluster.From(old)
	name := old.Name

	oldForPatch := clusterAgentForPatch(old)
	// to ensure they appear in the patch as preconditions
	oldForPatch.UID = ""
	oldForPatch.ResourceVersion = ""

	oldData, err := json.Marshal(oldForPatch)
	if err != nil {
		return fmt.Errorf("failed to Marshal old data for User %s|%s: %w", clusterName, name, err)
	}

	newForPatch := clusterAgentForPatch(obj)
	// to ensure they appear in the patch as preconditions
	newForPatch.UID = old.UID
	newForPatch.ResourceVersion = old.ResourceVersion

	newData, err := json.Marshal(newForPatch)
	if err != nil {
		return fmt.Errorf("failed to Marshal new data for User %s|%s: %w", clusterName, name, err)
	}

	patchBytes, err := jsonpatch.CreateMergePatch(oldData, newData)
	if err != nil {
		return fmt.Errorf("failed to create patch for User %s|%s: %w", clusterName, name, err)
	}

	var subresources []string
	if statusChanged {
		subresources = []string{"status"}
	}

	_, err = c.farosClientSet.Cluster(clusterName).EdgeV1alpha1().Agents(obj.Namespace).Patch(ctx, obj.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{}, subresources...)
	return err
}
