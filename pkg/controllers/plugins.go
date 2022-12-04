package controllers

import (
	"context"
	"fmt"
	"time"

	farosclientset "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/cluster"
	farosinformers "github.com/faroshq/faros-hub/pkg/client/informers/externalversions"
	pluginbindings "github.com/faroshq/faros-hub/pkg/controllers/service/plugin_bindings"
	pluginrequests "github.com/faroshq/faros-hub/pkg/controllers/service/plugin_requests"
	"github.com/faroshq/faros-hub/pkg/models"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

// runSystemPlugins controller is running in system plugins virtual workspace and is responsible for
// managing plugins
func (c *controllerManager) runSystemPlugins(ctx context.Context, plugins models.PluginsList) error {
	restConfig, err := c.clientFactory.GetWorkspaceRestConfig(ctx, c.config.ControllersWorkspace)
	if err != nil {
		return err
	}

	var rest *rest.Config
	// bootstrap rest config for controllers
	if kcpAPIsGroupPresent(restConfig) {
		if err := wait.PollImmediateInfinite(time.Second*5, func() (bool, error) {
			klog.Infof("looking up virtual workspace URL - %s", c.config.ControllersFarosPluginsAPIExportName)
			rest, err = restConfigForAPIExport(ctx, restConfig, c.config.ControllersFarosPluginsAPIExportName)
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

	ctrlRequests, err := pluginrequests.NewController(
		c.config,
		farosClientSet,
		c.kcpClientSet,
		informer.Plugins().V1alpha1().Requests(),
		plugins,
	)
	if err != nil {
		return err
	}

	ctrlBindings, err := pluginbindings.NewController(
		c.config,
		farosClientSet,
		informer.Plugins().V1alpha1().Bindings(),
	)
	if err != nil {
		return err
	}

	informer.Start(ctx.Done())
	informer.WaitForCacheSync(ctx.Done())

	go ctrlRequests.Start(ctx, 2)
	go ctrlBindings.Start(ctx, 2)

	<-ctx.Done()
	return nil
}
