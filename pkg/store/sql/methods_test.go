package storesql_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/models"
	databasetest "github.com/faroshq/faros-hub/test/util/database"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestCascade tests if records deletes cascades
func TestCascade(t *testing.T) {
	if os.Getenv("CI_ONLY") == "" {
		t.Skip("skipping postgres tests in non-CI environment")
		return
	}

	db, err := databasetest.NewPostgresTestingStore(t)
	require.NoError(t, err)

	ctx := context.Background()

	user, err := db.CreateUser(ctx, models.User{
		User: tenancyv1alpha1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: "user.name",
			},
		},
	})
	require.NoError(t, err)

	workspace, err := db.CreateWorkspace(ctx, models.Workspace{
		UserID: user.ID,
		Workspace: tenancyv1alpha1.Workspace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "workspace.name",
			},
		},
	})
	require.NoError(t, err)

	err = db.DeleteUser(ctx, *user)
	require.NoError(t, err)

	_, err = db.GetWorkspace(ctx, models.Workspace{
		ID: workspace.ID,
	})
	require.Error(t, err)
}
