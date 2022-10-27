package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/davecgh/go-spew/spew"
	"github.com/faroshq/faros-hub/pkg/models"
	"golang.org/x/oauth2"
)

// OIDCLogin is a http handler for oidc login
// /services/faros.sh/oidc/login
func (s *server) OIDCLogin() func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fall through, faros tunnels URL start by /services/faros-tunnels
			if !strings.HasPrefix(r.URL.Path, defaultOIDCLoginPathPrefix) {
				h.ServeHTTP(w, r)
				return
			}

			localRedirect := r.URL.Query().Get("redirect_uri")

			var scopes []string

			b := make([]byte, 16)
			rand.Read(b)
			state := base64.URLEncoding.EncodeToString(b)

			// Getting the session, it's not an issue if we error here
			session, err := s.oAuthSessions.Get(r, "sess")
			if err != nil {
				// print error
			}

			session.Values["state"] = state
			session.Values["redirect_uri"] = localRedirect
			err = s.oAuthSessions.Save(r, w, session)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed persist state: %q", r.Form), http.StatusBadRequest)
				return
			}

			redirectURL := s.config.ControllerExternalURL + defaultOIDCCallbackPathPrefix
			spew.Dump(redirectURL)

			authCodeURL := ""
			scopes = append(scopes, "openid", "profile", "email")
			if r.FormValue("offline_access") != "yes" {
				authCodeURL = s.oauth2Config(scopes).AuthCodeURL(state)
			} else {
				authCodeURL = s.oauth2Config(scopes).AuthCodeURL(state, oauth2.AccessTypeOffline)
			}
			spew.Dump(authCodeURL)
			http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
		})
	}
}

// OIDCCallback is a http handler for oidc login callback
// /services/faros.sh/oidc/callback
func (s *server) OIDCCallback() func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fall through, faros tunnels URL start by /services/faros-tunnels
			if !strings.HasPrefix(r.URL.Path, defaultOIDCCallbackPathPrefix) {
				h.ServeHTTP(w, r)
				return
			}

			var (
				token *oauth2.Token
			)

			ctx := oidc.ClientContext(r.Context(), s.client)

			var localRedirect string
			oauth2Config := s.oauth2Config(nil)
			switch r.Method {
			case http.MethodGet:
				// Authorization redirect callback from OAuth2 auth flow.
				if errMsg := r.FormValue("error"); errMsg != "" {
					http.Error(w, errMsg+": "+r.FormValue("error_description"), http.StatusBadRequest)
					return
				}
				code := r.FormValue("code")
				if code == "" {
					http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
					return
				}

				session, err := s.oAuthSessions.Get(r, "sess")
				if err != nil {
					http.Error(w, fmt.Sprintf("no session present: %q", r.Form), http.StatusBadRequest)
					return
				}

				localRedirect = session.Values["redirect_uri"].(string)

				if state := r.FormValue("state"); state != session.Values["state"] {
					http.Error(w, fmt.Sprintf("expected state %q got %q", session.Values["state"], state), http.StatusBadRequest)
					return
				}
				token, err = oauth2Config.Exchange(ctx, code)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
					return
				}
			case http.MethodPost:
				// Form request from frontend to refresh a token.
				refresh := r.FormValue("refresh_token")
				if refresh == "" {
					http.Error(w, fmt.Sprintf("no refresh_token in request: %q", r.Form), http.StatusBadRequest)
					return
				}
				t := &oauth2.Token{
					RefreshToken: refresh,
					Expiry:       time.Now().Add(-time.Hour),
				}
				var err error
				token, err = oauth2Config.TokenSource(ctx, t).Token()
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
					return
				}
			default:
				http.Error(w, fmt.Sprintf("method not implemented: %s", r.Method), http.StatusBadRequest)
				return
			}

			rawIDToken, ok := token.Extra("id_token").(string)
			if !ok {
				http.Error(w, "no id_token in token response", http.StatusInternalServerError)
				return
			}

			idToken, err := s.verifier.Verify(r.Context(), rawIDToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to verify ID token: %v", err), http.StatusInternalServerError)
				return
			}

			// TODO: extend
			var claims struct {
				Email         string `json:"email"`
				EmailVerified bool   `json:"email_verified"`
			}
			if err := idToken.Claims(&claims); err != nil {
				http.Error(w, fmt.Sprintf("failed to parse claim: %v", err), http.StatusInternalServerError)
				return
			}

			response := models.LoginResponse{
				IDToken:    *idToken,
				RawIDToken: rawIDToken,
				Email:      claims.Email,
			}

			data, err := json.Marshal(response)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
				return
			}

			base64.StdEncoding.EncodeToString(data)

			localRedirect = localRedirect + "?data=" + base64.StdEncoding.EncodeToString(data)
			http.Redirect(w, r, localRedirect, http.StatusSeeOther)

		})

	}
}

func (s *server) oauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.OIDCClientID,
		ClientSecret: s.config.OIDCClientSecret,
		Endpoint:     s.provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  s.redirectURL,
	}
}
