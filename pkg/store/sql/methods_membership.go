package storesql

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/faroshq/faros-hub/pkg/store"
)

// GetMembership gets workspaces based on membership ID
func (s *Store) GetMembership(ctx context.Context, p models.Membership) (*models.Membership, error) {
	switch {
	case p.ID != "":
		// OK, getting by ID
	default:
		return nil, store.ErrFailToQuery
	}

	result := models.Membership{}
	if err := s.db.WithContext(ctx).Where(&p).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return &result, nil
}

// CreateMembership creates workspace user membership object
func (s *Store) CreateMembership(ctx context.Context, p models.Membership) (*models.Membership, error) {
	switch {
	case p.UserID != "" && p.WorkspaceID != "":
		// OK, getting by User_ID and Workspace_ID
	default:
		return nil, store.ErrFailToQuery
	}
	p.ID = uuid.New().String()

	err := s.db.WithContext(ctx).Create(&p).Error
	if err != nil {
		return nil, err
	}

	s.notifyUpdatedMembership(ctx, p.ID, models.EventCreated)

	return s.GetMembership(ctx, models.Membership{ID: p.ID})
}
