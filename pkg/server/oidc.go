package server

import (
	"net/http"
	"strings"
)

// oidcLogin is a http handler for oidc login
// /faros.sh/oidc/login
func (s *Service) oidcLogin() func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fall through, faros tunnels URL start by /services/faros-tunnels
			if !strings.HasPrefix(r.URL.Path, defaultOIDCLoginPathPrefix) {
				h.ServeHTTP(w, r)
				return
			}

			s.authenticator.OIDCLogin(w, r)
		})
	}
}

// oidcCallback is a http handler for oidc login callback
// /faros.sh/oidc/callback
func (s *Service) oidcCallback() func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fall through, faros tunnels URL start by /services/faros-tunnels
			if !strings.HasPrefix(r.URL.Path, defaultOIDCCallbackPathPrefix) {
				h.ServeHTTP(w, r)
				return
			}

			s.authenticator.OIDCCallback(w, r)

		})

	}
}
