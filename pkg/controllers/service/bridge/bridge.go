package bridge

import (
	"context"
	"time"

	farosclientset "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	farosclusyrtclientset "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/cluster"
	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/faroshq/faros-hub/pkg/store"
	storesql "github.com/faroshq/faros-hub/pkg/store/sql"
	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/klog/v2"
)

type Bridge interface {
	Run(ctx context.Context)
}

type bridge struct {
	store          store.Store
	config         *config.ControllerConfig
	farosClientSet farosclientset.Interface
}

func New(ctx context.Context, config *config.ControllerConfig) (Bridge, error) {
	store, err := storesql.NewStore(ctx, &config.Database)
	if err != nil {
		return nil, err
	}

	farosClientSet, err := farosclusyrtclientset.NewForConfig(config.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	return &bridge{
		store:          store,
		config:         config,
		farosClientSet: farosClientSet.Cluster(logicalcluster.New(config.ControllersTenantWorkspace)),
	}, nil
}

func (b *bridge) Run(ctx context.Context) {
	logger := klog.FromContext(ctx)

	changesCh := make(chan *models.Event)

	go func() {
		logger.Info("Subscribing to changes")
		defer logger.Info("Unsubscribing from changes")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := b.store.SubscribeChanges(ctx, func(event *models.Event) error {
					changesCh <- event
					return nil
				})
				if err != nil {
					logger.Error(err, "failed to subscribe to changes")
				}
				// Retry to subscribe
				time.Sleep(time.Second)
			}
		}
	}()

	// Start periodically reschedule applications on all devices or individual
	// ones if only minimal changes are required
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-changesCh:
			switch event.Resource {
			case models.EventResourceMembership:
				logger.Info("Membership change detected")
			case models.EventResourceWorkspace:
				switch event.Type {
				case models.EventCreated:
					logger.Info("Workspace create")
					err := b.createWorkspace(ctx, event.ObjectID)
					if err != nil {
						logger.Error(err, "failed to create workspace")
					}
				case models.EventDeleted:
					logger.Info("Workspace delete")
					err := b.deleteWorkspace(ctx, event.ObjectID)
					if err != nil {
						logger.Error(err, "failed to delete workspace")
					}
				case models.EventUpdated:
					logger.Info("Workspace update")
					err := b.updateWorkspace(ctx, event.ObjectID)
					if err != nil {
						logger.Error(err, "failed to update workspace")
					}
				}
			case models.EventResourceUser:
				switch event.Type {
				case models.EventCreated:
					logger.Info("User create")
					err := b.createUser(ctx, event.ObjectID)
					if err != nil {
						logger.Error(err, "failed to create user")
					}
				case models.EventDeleted:
					logger.Info("User delete")
					err := b.deleteUser(ctx, event.ObjectID)
					if err != nil {
						logger.Error(err, "failed to delete user")
					}
				case models.EventUpdated:
					logger.Info("User update")
					err := b.updateUser(ctx, event.ObjectID)
					if err != nil {
						logger.Error(err, "failed to update user")
					}
				}
			}
		}
	}
}
