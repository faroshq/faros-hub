package storesql

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/faroshq/faros-hub/pkg/store"
)

// GetWorkspace gets workspaces based on workspace ID
func (s *Store) GetWorkspace(ctx context.Context, p models.Workspace) (*models.Workspace, error) {
	switch {
	case p.ID != "":
		// OK, getting by ID
	case p.UserID != "" && p.Name != "":
		// OK, getting by User_ID and Name
	default:
		return nil, store.ErrFailToQuery
	}

	result := models.Workspace{}
	if err := s.db.WithContext(ctx).Where(&p).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return &result, nil
}

// CreateWorkspace creates workspace object
func (s *Store) CreateWorkspace(ctx context.Context, p models.Workspace) (*models.Workspace, error) {
	p.ID = uuid.New().String()

	err := s.db.WithContext(ctx).Create(&p).Error
	if err != nil {
		return nil, err
	}

	// create membership for the user
	_, err = s.CreateMembership(ctx, models.Membership{
		UserID:      p.UserID,
		WorkspaceID: p.ID,
	})
	if err != nil {
		return nil, err
	}

	s.notifyUpdatedWorkspace(ctx, p.ID, models.EventCreated)

	return s.GetWorkspace(ctx, models.Workspace{ID: p.ID})
}

// UpdateWorkspace updates workspace based on workspace ID
func (s *Store) UpdateWorkspace(ctx context.Context, p models.Workspace) (*models.Workspace, error) {
	switch {
	case p.ID != "":
		// OK, getting by ID
	case p.UserID != "" && p.Name != "":
		// OK, getting by User_ID and Name
	default:
		return nil, store.ErrFailToQuery
	}

	query := models.Workspace{ID: p.ID}
	err := s.db.WithContext(ctx).Model(&models.Workspace{}).Where(&query).Save(&p).Error
	if err != nil {
		return nil, err
	}

	s.notifyUpdatedWorkspace(ctx, p.ID, models.EventUpdated)

	return s.GetWorkspace(ctx, models.Workspace{ID: p.ID})
}

// DeleteWorkspace deletes workspace based on workspace ID
func (s *Store) DeleteWorkspace(ctx context.Context, p models.Workspace) error {
	switch {
	case p.ID != "":
		// OK, getting by ID
	default:
		return store.ErrFailToQuery
	}

	s.notifyUpdatedWorkspace(ctx, p.ID, models.EventDeleted)

	return s.db.WithContext(ctx).Delete(&p).Error
}

// ListWorkspaces lists workspace based on ID
func (s *Store) ListWorkspaces(ctx context.Context, p models.Workspace) ([]models.Workspace, error) {
	switch {
	case p.UserID != "":
		// OK, listing by User_ID
	default:
		return nil, store.ErrFailToQuery
	}

	results := []models.Workspace{}
	if err := s.db.WithContext(ctx).Where(&p).Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return results, nil
}

func (s *Store) ListWorkspacesByMembership(ctx context.Context, p models.User) ([]models.Workspace, error) {
	switch {
	case p.ID != "":
		// OK, listing by User_ID
	default:
		return nil, store.ErrFailToQuery
	}

	results := []models.Membership{}
	if err := s.db.WithContext(ctx).Where(&models.Membership{UserID: p.ID}).Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	workspaces := []models.Workspace{}
	for _, m := range results {
		w, err := s.GetWorkspace(ctx, models.Workspace{ID: m.WorkspaceID})
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, *w)
	}

	return workspaces, nil
}
