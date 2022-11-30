package server

import (
	"net/http"

	"github.com/faroshq/faros-hub/pkg/models"
	"k8s.io/klog/v2"
)

func (s *Service) authenticate(w http.ResponseWriter, r *http.Request) (bool, *models.User, error) {
	authenticated, user, err := s.authenticator.Authenticate(r)
	if err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return false, nil, err
	}

	if !authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false, nil, nil
	}

	return true, user, nil
}
