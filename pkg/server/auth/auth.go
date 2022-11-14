package auth

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/golang-jwt/jwt/request"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/kcp-dev/logicalcluster/v2"
	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utiltls "github.com/faroshq/faros-hub/pkg/util/tls"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

const (
	// UserLabel is the label used to identify the user
	UserLabel = "faros.sh/user"
)

// Authenticator authenticator is used to authenticate and handle all authentication related tasks
type Authenticator interface {
	// OIDCLogin will redirect user to OIDC provider
	OIDCLogin(w http.ResponseWriter, r *http.Request)
	// OIDCCallback will handle OIDC callback
	OIDCCallback(w http.ResponseWriter, r *http.Request)
	// Authenticate will authenticate the request if user already exists
	Authenticate(r *http.Request) (authenticated bool, user *tenancyv1alpha1.User, err error)
	// ParseJWTToken will parse the JWT token and return the user
	ParseJWTToken(ctx context.Context, token string) (user *tenancyv1alpha1.User, err error)
}

// Static check
var _ Authenticator = &AuthenticatorImpl{}

type AuthenticatorImpl struct {
	config *config.APIConfig

	oAuthSessions *sessions.CookieStore
	provider      *oidc.Provider
	verifier      *oidc.IDTokenVerifier
	redirectURL   string
	client        *http.Client

	farosClient farosclient.ClusterInterface
	coreClient  kubernetes.ClusterInterface

	cluster logicalcluster.Name
}

func NewAuthenticator(cfg *config.APIConfig, coreClient kubernetes.ClusterInterface, farosClient farosclient.ClusterInterface, callbackURLPrefix string) (*AuthenticatorImpl, error) {
	var client *http.Client
	var err error

	hostingCoreClient, err := kubernetes.NewForConfig(cfg.HostingClusterRestConfig)
	if err != nil {
		return nil, err
	}

	secret, err := hostingCoreClient.CoreV1().Secrets(cfg.HostingClusterNamespace).Get(context.Background(), cfg.OIDCCASecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	crt, ok := secret.Data["tls.crt"]
	if !ok {
		return nil, errors.New("oidc tls.crt not found in secret")
	}
	key, ok := secret.Data["tls.key"]
	if !ok {
		return nil, errors.New("oidc tls.key not found in secret")
	}
	client, err = httpClientForRootCAs(crt, key)
	if err != nil {
		return nil, err
	}

	redirectURL := cfg.ControllerExternalURL + callbackURLPrefix

	ctx := oidc.ClientContext(context.Background(), client)

	provider, err := oidc.NewProvider(ctx, cfg.OIDCIssuerURL)
	if err != nil {
		return nil, err
	}
	// Create an ID token parser, but only trust ID tokens issued to "example-app"
	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.OIDCClientID,
	})

	da := &AuthenticatorImpl{
		config:        cfg,
		farosClient:   farosClient,
		coreClient:    coreClient,
		verifier:      verifier,
		provider:      provider,
		client:        client,
		redirectURL:   redirectURL,
		oAuthSessions: sessions.NewCookieStore([]byte(cfg.OIDCAuthSessionKey)),
		cluster:       logicalcluster.New(cfg.ControllersTenantWorkspace),
	}
	return da, nil
}

func (a *AuthenticatorImpl) OIDCLogin(w http.ResponseWriter, r *http.Request) {
	localRedirect := r.URL.Query().Get("redirect_uri")

	var scopes []string

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Getting the session, it's not an issue if we error here
	session, err := a.oAuthSessions.Get(r, "sess")
	if err != nil {
		// print error
	}

	session.Values["state"] = state
	session.Values["redirect_uri"] = localRedirect
	err = a.oAuthSessions.Save(r, w, session)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed persist state: %q", r.Form), http.StatusBadRequest)
		return
	}

	authCodeURL := ""
	scopes = append(scopes, "openid", "profile", "email")
	if r.FormValue("offline_access") != "yes" {
		authCodeURL = a.oauth2Config(scopes).AuthCodeURL(state)
	} else {
		authCodeURL = a.oauth2Config(scopes).AuthCodeURL(state, oauth2.AccessTypeOffline)
	}

	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}

func (a *AuthenticatorImpl) OIDCCallback(w http.ResponseWriter, r *http.Request) {
	var (
		token *oauth2.Token
	)

	ctx := oidc.ClientContext(r.Context(), a.client)

	var localRedirect string
	oauth2Config := a.oauth2Config(nil)
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

		session, err := a.oAuthSessions.Get(r, "sess")
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

	idToken, err := a.verifier.Verify(r.Context(), rawIDToken)
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

	_, err = a.registerOrUpdateUser(ctx, &tenancyv1alpha1.User{
		Spec: tenancyv1alpha1.UserSpec{
			Email: claims.Email,
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to register user: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		IDToken:       *idToken,
		RawIDToken:    rawIDToken,
		Email:         claims.Email,
		ServerBaseURL: fmt.Sprintf("%s/clusters", a.config.ControllerExternalURL),
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	base64.StdEncoding.EncodeToString(data)

	localRedirect = localRedirect + "?data=" + base64.StdEncoding.EncodeToString(data)
	http.Redirect(w, r, localRedirect, http.StatusSeeOther)

}

func (a *AuthenticatorImpl) Authenticate(r *http.Request) (authenticated bool, user *tenancyv1alpha1.User, err error) {

	// Trying to authenticate via URL query (websocket for SSH/logs, SSE)
	if urlQueryToken := r.URL.Query().Get("_t"); urlQueryToken != "" {
		user, err = a.ParseJWTToken(r.Context(), urlQueryToken)
		if err != nil {
			return false, nil, err
		}

		// authenticated
		return true, user, nil
	}

	if r.Header.Get("Authorization") == "" {
		return false, nil, nil
	}

	// If it's basic auth (service account), it will have 'Basic' instead of
	// 'Bearer'
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer") {
		return false, nil, nil
	}

	token, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
	if err != nil {
		return false, nil, err
	}

	user, err = a.ParseJWTToken(r.Context(), token)
	if err != nil {
		return false, nil, err
	}

	// authenticated
	return true, user, nil
}

// ParseJWTToken validates token's validity and returns models.User that the token belongs to
func (a *AuthenticatorImpl) ParseJWTToken(ctx context.Context, token string) (user *tenancyv1alpha1.User, err error) {
	idToken, err := a.verifier.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	// TODO: extend
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return a.getUser(ctx, claims.Email)
}

// return an HTTP client which trusts the provided root CAs.
func httpClientForRootCAs(crt, key []byte) (*http.Client, error) {
	c, k, err := utiltls.CertificatePairFromBytes(crt, key)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AddCert(c)

	tlsConfig := &tls.Config{
		RootCAs: pool,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{
					crt,
				},
				PrivateKey: k,
			},
		},
		ServerName:         "faros",
		InsecureSkipVerify: true,
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}, nil
}

func (a *AuthenticatorImpl) oauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.config.OIDCClientID,
		ClientSecret: a.config.OIDCClientSecret,
		Endpoint:     a.provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  a.redirectURL,
	}
}

// registerOrUpdateUser will register or update user in the system when user is authenticated
// TODO: This is not quite right place for this
func (a *AuthenticatorImpl) registerOrUpdateUser(ctx context.Context, user *tenancyv1alpha1.User) (*tenancyv1alpha1.User, error) {
	if user.Name == "" {
		user.Name = uuid.New().String()
	}

	// we will be selecting based on labels, but k8s does not allow symbols like '@' in labels
	// so we will replace it with '-at-' checking before if only single @ exists

	current, err := a.getUser(ctx, user.Spec.Email)
	if err != nil && err != errUserNotFound {
		return nil, err
	}

	// TODO: Duplicate with f below
	labelEmail := strings.Replace(user.Spec.Email, "@", "-at-", 1)

	if current != nil {
		current.Spec = user.Spec
		if current.Labels == nil {
			current.Labels = make(map[string]string)
		}

		current.Labels[UserLabel] = labelEmail
		return a.farosClient.Cluster(a.cluster).TenancyV1alpha1().Users().Update(ctx, current, metav1.UpdateOptions{})
	} else {
		user.Labels = map[string]string{
			UserLabel: labelEmail,
		}
		user, err = a.farosClient.Cluster(a.cluster).TenancyV1alpha1().Users().Create(ctx, user, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	}

	// provision user namespace
	_, err = a.coreClient.Cluster(a.cluster).CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: user.Name,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return user, nil
}

var (
	errUserNotFound       = errors.New("user not found")
	errMultipleUsersFound = errors.New("multiple users found")
)

func (a *AuthenticatorImpl) getUser(ctx context.Context, email string) (*tenancyv1alpha1.User, error) {
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email address")
	}

	labelEmail := strings.Replace(email, "@", "-at-", 1)
	users, err := a.farosClient.Cluster(a.cluster).TenancyV1alpha1().Users().List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", UserLabel, labelEmail),
	})
	if err != nil {
		return nil, err
	}

	switch len(users.Items) {
	case 0:
		return nil, errUserNotFound
	case 1:
		return &users.Items[0], nil
	default:
		return nil, errMultipleUsersFound
	}
}
