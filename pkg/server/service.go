package server

import (
	"context"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/coreos/go-oidc"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/gorilla/sessions"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type contextKey int

const (
	contextKeyForwardPath contextKey = iota
	contextKeyResponse
	contextKeyClient
)

const (
	kubeconfigTimeout = time.Hour * 24
)

var _ Interface = &Service{}

type Interface interface {
	// GetHandlers returns a list of handlers that should be added to the HTTP server
	GetHandlers() []func(h http.Handler) http.HandlerFunc
	// SeedClients will inject all api server clients with post-start-hook
	SeedClients(rest *rest.Config) error
}

func New(config *config.ControllerConfig) (*Service, error) {
	var client *http.Client
	var err error
	if config.OIDCCAFile != "" {
		client, err = httpClientForRootCAs(config.OIDCCAFile)
		if err != nil {
			return nil, err
		}
	}

	if client == nil {
		client = http.DefaultClient
	}

	redirectURL := config.ControllerExternalURL + defaultOIDCCallbackPathPrefix

	ctx := oidc.ClientContext(context.Background(), client)

	provider, err := oidc.NewProvider(ctx, config.OIDCIssuerURL)
	if err != nil {
		return nil, err
	}
	// Create an ID token parser, but only trust ID tokens issued to "example-app"
	idTokenVerifier := provider.Verifier(&oidc.Config{
		ClientID: config.OIDCClientID,
	})

	return &Service{
		client:        client,
		provider:      provider,
		verifier:      idTokenVerifier,
		redirectURL:   redirectURL,
		config:        config,
		oAuthSessions: sessions.NewCookieStore([]byte(config.OIDCAuthSessionKey)),
	}, nil
}

type Service struct {
	config *config.ControllerConfig

	// oidc tooling
	oAuthSessions *sessions.CookieStore
	provider      *oidc.Provider
	verifier      *oidc.IDTokenVerifier
	redirectURL   string
	client        *http.Client

	// tunneling tooling
	kcpClient   kcpclient.ClusterInterface
	farosClient farosclient.ClusterInterface
	coreClients kubernetes.ClusterInterface
	proxy       *httputil.ReverseProxy

	// set once above tunneling tooling is seeded
	seeded bool
}

var (
	defaultTunnelsPathPrefix      = "/faros.sh/tunnels"
	defaultOIDCLoginPathPrefix    = "/faros.sh/oidc/login"
	defaultOIDCCallbackPathPrefix = "/faros.sh/oidc/callback"
	defaultWorkspaceManagement    = "/faros.sh/workspace"
)

func (s *Service) GetHandlers() []func(h http.Handler) http.HandlerFunc {
	return []func(h http.Handler) http.HandlerFunc{
		s.oidcCallback(),
		s.oidcLogin(),
		s.customTunnels(),
	}
}
