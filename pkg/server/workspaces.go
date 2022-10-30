package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/kcp-dev/logicalcluster/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

// workspacesHandler is a http handler for workspaces operations
// GET -  faros.sh/workspaces - list all workspaces for users
// GET -  faros.sh/workspaces/<workspace> - get workspace details
// POST - faros.sh/workspaces - create new workspace
func (s *Service) workspacesHandler() func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, defaultWorkspaceManagement) {
				h.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()

			authenticated, user, err := s.authenticator.Authenticate(r)
			if err != nil {
				klog.Error(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if !authenticated {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			switch r.Method {
			case http.MethodGet:
				parts := strings.Split(r.URL.Path, defaultWorkspaceManagement)
				if len(parts) == 2 && parts[1] == "" { // no workspace name - list all workspaces
					_, err := s.listWorkspaces(ctx, *user)
					if err != nil {
						klog.Error(err)
						http.Error(w, "Internal server error", http.StatusInternalServerError)
					}
					return
				}
				if len(parts) == 2 && parts[1] != "" { // workspace name - get workspace details
					//err := s.getWorkspace(w, r, user, strings.TrimPrefix(parts[1], "/"))
					//////////if err != nil {
					////////////////////	klog.Error(err)
					////////////////////	http.Error(w, "Internal server error", http.StatusInternalServerError)
					//////////}
					//return
				}

			case http.MethodPost:
				spew.Dump(user)

			}
		})
	}
}

func (s *Service) listWorkspaces(ctx context.Context, user tenancyv1alpha1.User) (*tenancyv1alpha1.WorkspaceList, error) {
	cluster := logicalcluster.New(s.config.ControllersTenantWorkspace)
	return s.farosClient.Cluster(cluster).TenancyV1alpha1().Workspaces(user.Name).List(ctx, metav1.ListOptions{})
}
