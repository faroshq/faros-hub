package server

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/kcp-dev/logicalcluster/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog/v2"
	"k8s.io/utils/strings/slices"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

// workspacesHandler is a http handler for workspaces operations
// GET -  faros.sh/workspaces - list all workspaces for users
// GET -  faros.sh/workspaces/<workspace> - get workspace details
// DELETE - faros.sh/workspaces/<workspace> - delete a workspace
// POST - faros.sh/workspaces - create new workspace
func (s *Service) workspacesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := logicalcluster.WithCluster(r.Context(), logicalcluster.New(s.config.ControllersTenantWorkspace))

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
	// list/get
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, path.Join(pathAPIVersion, pathWorkspaces))
		spew.Dump(parts)
		if len(parts) == 2 && parts[1] == "" { // no workspace name - list all workspaces
			workspaces, err := s.listWorkspaces(ctx, *user)
			if err != nil {
				responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
				return
			}
			responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, workspaces)
			return
		} else if len(parts) == 2 && parts[1] != "" { // workspace name - get workspace details
			workspace, err := s.getWorkspace(ctx, *user, strings.TrimPrefix(parts[1], "/"))
			if err != nil {
				responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
				return
			}
			responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, workspace)
			return
		}
		// create
	case http.MethodPost:
		request := &tenancyv1alpha1.Workspace{}
		limitedReader := &io.LimitedReader{R: r.Body, N: limit}
		body, err := ioutil.ReadAll(limitedReader)
		if err != nil {
			responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
			return
		}
		if err := runtime.DecodeInto(codecs.UniversalDecoder(), body, request); err != nil {
			responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
			return
		}

		request.Namespace = user.Name
		if !slices.Contains(request.Spec.Members, user.Spec.Email) {
			request.Spec.Members = append(request.Spec.Members, user.Spec.Email)
		}

		workspace, err := s.farosClient.TenancyV1alpha1().Workspaces(user.Name).Create(ctx, request, metav1.CreateOptions{})
		if err != nil {
			responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
			return
		}
		responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusCreated, workspace)
	case http.MethodDelete:
		parts := strings.Split(r.URL.Path, path.Join(pathAPIVersion, pathWorkspaces))
		if len(parts) == 2 && parts[1] != "" {
			workspace, err := s.getWorkspace(ctx, *user, strings.TrimPrefix(parts[1], "/"))
			if err != nil {
				responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
				return
			}
			err = s.deleteWorkspace(ctx, *user, workspace.Name)
			if err != nil {
				responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
				return
			}
			responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, workspace)
			return
		}
	}
}

func (s *Service) listWorkspaces(ctx context.Context, user tenancyv1alpha1.User) (*tenancyv1alpha1.WorkspaceList, error) {
	return s.farosClient.TenancyV1alpha1().Workspaces(user.Name).List(ctx, metav1.ListOptions{})
}

func (s *Service) getWorkspace(ctx context.Context, user tenancyv1alpha1.User, name string) (*tenancyv1alpha1.Workspace, error) {
	return s.farosClient.TenancyV1alpha1().Workspaces(user.Name).Get(ctx, name, metav1.GetOptions{})
}

func (s *Service) deleteWorkspace(ctx context.Context, user tenancyv1alpha1.User, name string) error {
	return s.farosClient.TenancyV1alpha1().Workspaces(user.Name).Delete(ctx, name, metav1.DeleteOptions{})
}
