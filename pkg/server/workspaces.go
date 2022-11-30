package server

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog/v2"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/models"
)

func (s *Service) getWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authenticated, user, err := s.authenticate(w, r)
	if err != nil || !authenticated {
		return
	}

	vars := mux.Vars(r)
	workspaceName := vars["workspace"]
	if workspaceName == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	workspaceRef, err := s.store.GetWorkspace(ctx, models.Workspace{
		UserID: user.ID,
		Name:   workspaceName,
	})
	if err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	workspaces, err := s.farosClient.TenancyV1alpha1().Workspaces().Get(ctx, workspaceRef.Workspace.Name, metav1.GetOptions{})
	if err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, workspaces)
}

func (s *Service) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authenticated, user, err := s.authenticate(w, r)
	if err != nil || !authenticated {
		return
	}

	workspacesRef, err := s.store.ListWorkspaces(ctx, models.Workspace{
		UserID: user.ID,
	})
	if err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result := &tenancyv1alpha1.WorkspaceList{}
	for _, workspaceRef := range workspacesRef {
		workspace, err := s.farosClient.TenancyV1alpha1().Workspaces().Get(ctx, workspaceRef.Workspace.Name, metav1.GetOptions{})
		if err != nil {
			klog.Error(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		result.Items = append(result.Items, *workspace)
	}

	responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, result)
}

func (s *Service) createWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authenticated, user, err := s.authenticate(w, r)
	if err != nil || !authenticated {
		return
	}

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

	// name is unique per user
	workspaces := &models.Workspace{
		Name:   request.Name,
		UserID: user.ID,
	}

	_, err = s.store.GetWorkspace(ctx, *workspaces)
	if err == nil {
		http.Error(w, "Workspace already exists", http.StatusConflict)
		return
	}

	workspace := models.Workspace{
		UserID: user.ID,
		Name:   request.Name,
		Workspace: tenancyv1alpha1.Workspace{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid.New().String(),
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       tenancyv1alpha1.WorkspaceKind,
				APIVersion: tenancyv1alpha1.SchemeGroupVersion.String(),
			},
			Spec: tenancyv1alpha1.WorkspaceSpec{
				Name: request.Name,
				Members: []string{
					user.Email,
				},
			},
		},
	}

	workspaceCreated, err := s.store.CreateWorkspace(ctx, workspace)
	if err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, &workspaceCreated.Workspace)
}

func (s *Service) updateWorkspace(w http.ResponseWriter, r *http.Request) {
	// not implemented
	responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, tenancyv1alpha1.SchemeGroupVersion, w, r, http.StatusNotImplemented, nil)
}

func (s *Service) deleteWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authenticated, user, err := s.authenticate(w, r)
	if err != nil || !authenticated {
		return
	}

	vars := mux.Vars(r)
	workspaceName := vars["workspace"]
	if workspaceName == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	workspace, err := s.store.GetWorkspace(ctx, models.Workspace{
		UserID: user.ID,
		Name:   workspaceName,
	})
	if err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := s.store.DeleteWorkspace(ctx, *workspace); err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
