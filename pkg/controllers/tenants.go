package controllers

import (
	"context"
	"fmt"
	"time"

	farosclientset "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/cluster"
	farosinformers "github.com/faroshq/faros-hub/pkg/client/informers/externalversions"
	"github.com/faroshq/faros-hub/pkg/controllers/service/workspaces"
	"github.com/kcp-dev/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

// runSystem controller is running in system workspace and is responsible for
// managing workspaces and tenants
func (c *controllerManager) runSystem(ctx context.Context) error {
	//cluster := logicalcluster.New(c.config.ControllersWorkspace)

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
			klog.Infof("looking up virtual workspace URL - %s", c.config.ControllersFarosTenancyAPIExportName)
			rest, err = restConfigForAPIExport(ctx, restConfig, c.config.ControllersFarosTenancyAPIExportName)
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

	ctrl, err := workspaces.NewController(
		c.config,
		c.kcpClientSet, // root client to create bindings and workspaces
		coreClientSet,
		farosClientSet, // client to manage workspaces
		informer.Tenancy().V1alpha1().Workspaces(),
	)
	if err != nil {
		return err
	}

	informer.Start(ctx.Done())
	informer.WaitForCacheSync(ctx.Done())

	ctrl.Start(ctx, 2)

	<-ctx.Done()
	return nil
}
