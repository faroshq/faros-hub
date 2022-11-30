package store

import (
	"context"
	"errors"

	"github.com/faroshq/faros-hub/pkg/models"
)

type Store interface {
	GetWorkspace(context.Context, models.Workspace) (*models.Workspace, error)
	ListWorkspaces(context.Context, models.Workspace) ([]models.Workspace, error)
	ListWorkspacesByMembership(context.Context, models.User) ([]models.Workspace, error)
	DeleteWorkspace(context.Context, models.Workspace) error
	CreateWorkspace(context.Context, models.Workspace) (*models.Workspace, error)
	UpdateWorkspace(context.Context, models.Workspace) (*models.Workspace, error)

	GetUser(context.Context, models.User) (*models.User, error)
	ListUsers(context.Context, models.User) ([]models.User, error)
	DeleteUser(context.Context, models.User) error
	CreateUser(context.Context, models.User) (*models.User, error)
	UpdateUser(context.Context, models.User) (*models.User, error)

	SubscribeChanges(ctx context.Context, callback func(event *models.Event) error) error

	// Status is a health check endpoint
	Status() (interface{}, error)
	RawDB() interface{}
	Close() error
}

var ErrFailToQuery = errors.New("malformed request. failed to query")
var ErrRecordNotFound = errors.New("object not found")
