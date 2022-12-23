package models

import (
	"time"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

// Workspace is a model for the Workspace database model storing the workspace information.
// It extends the Workspace CRD with additional fields for database storage.
type Workspace struct {
	ID        string    `json:"id" yaml:"id" gorm:"primaryKey,uniqueIndex"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt" grom:"index"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
	UserID    string    `json:"userId" yaml:"userId" gorm:"index"`
	// Name is user facing name of the workspace
	Name string `json:"name" yaml:"name" gorm:"index"`

	Workspace tenancyv1alpha1.Workspace `json:"workspace" yaml:"workspace" gorm:"json"`
}

// User is a model for the User database model storing the user information.
// It extends the User CRD with additional fields for database storage.
type User struct {
	ID        string    `json:"id" yaml:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
	// Email is the email of the user. Must be unique.
	Email string `json:"email" yaml:"email" gorm:"uniqueIndex"`

	User tenancyv1alpha1.User `json:"user" yaml:"user" gorm:"json"`
}

// Membership is a model for the Membership database model storing the membership information.
type Membership struct {
	ID          string    `json:"id" yaml:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" yaml:"updatedAt"`
	UserID      string    `json:"userId" yaml:"userId" gorm:"index"`
	WorkspaceID string    `json:"workspaceId" yaml:"workspaceId" gorm:"index"`
}
