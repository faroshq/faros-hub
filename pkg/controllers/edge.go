package controllers

import (
	"context"
	"fmt"
	"time"

	farosclientset "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/cluster"
	farosinformers "github.com/faroshq/faros-hub/pkg/client/informers/externalversions"
	"github.com/faroshq/faros-hub/pkg/controllers/tenants/edge/agent"
	"github.com/faroshq/faros-hub/pkg/controllers/tenants/edge/registration"
	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/kcp-dev/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

// edge controller is running in edge api virtual workspace context

func (c *controllerManager) runEdge(ctx context.Context, plugins models.PluginsList) error {
	restConfig, err := c.clientFactory.GetWorkspaceRestConfig(ctx, c.config.ControllersWorkspace)
	if err != nil {
		return err
	}

	rootRestConfig, err := c.clientFactory.GetRootRestConfig()
	if err != nil {
		return err
	}

	coreClientSet, err := kubernetes.NewForConfig(rootRestConfig)
	if err != nil {
		return err
	}

	var rest *rest.Config
	// bootstrap rest config for controllers
	if kcpAPIsGroupPresent(restConfig) {
		if err := wait.PollImmediateInfinite(time.Second*5, func() (bool, error) {
			klog.Infof("looking up virtual workspace URL - %s", c.config.ControllersFarosEdgeAPIExportName)
			rest, err = restConfigForAPIExport(ctx, restConfig, c.config.ControllersFarosEdgeAPIExportName)
			if err != nil {
				return false, nil
			}
			return true, nil
		}); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("kcp APIs group not present in cluster. We don't support non kcp clusters yet")
	}

	farosClientSet, err := farosclientset.NewForConfig(rest)
	if err != nil {
		return err
	}

	// Must always follow the order. Otherwise informers are not initialized
	// 1. create shared informer factory
	// 2. get listers and informers out of the factory in controller constructors
	// 3. start the factory
	// 4. wait for the factory to sync.
	informer := farosinformers.NewSharedInformerFactory(farosClientSet, resyncPeriod)

	ctrlRegistration, err := registration.NewController(
		c.config,
		coreClientSet,
		farosClientSet,
		informer.Edge().V1alpha1().Registrations(),
	)

	ctrlAgent, err := agent.NewController(
		c.config,
		coreClientSet,
		farosClientSet,
		informer.Edge().V1alpha1().Agents(),
	)

	informer.Start(ctx.Done())
	informer.WaitForCacheSync(ctx.Done())

	ctrlRegistration.Start(ctx, 2)
	ctrlAgent.Start(ctx, 2)

	<-ctx.Done()
	return nil

}
