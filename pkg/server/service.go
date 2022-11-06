package server

import (
	"net/http"
	"net/http/httputil"
	"time"

	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/server/auth"
	"github.com/faroshq/faros-hub/pkg/util/roundtripper"
	kcptenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	workloadv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/workload/v1alpha1"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(accessv1alpha1.AddToScheme(scheme))
	utilruntime.Must(edgev1alpha1.AddToScheme(scheme))
	utilruntime.Must(workloadv1alpha1.AddToScheme(scheme))
	utilruntime.Must(kcptenancyv1alpha1.AddToScheme(scheme))
	utilruntime.Must(tenancyv1alpha1.AddToScheme(scheme))
}

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
	// Init will inject all api server clients with post-start-hook
	Init(rest *rest.Config) error
}

func New(config *config.ControllerConfig) (*Service, error) {
	return &Service{
		config: config,
	}, nil
}

type Service struct {
	config        *config.ControllerConfig
	authenticator auth.Authenticator

	// oidc tooling
	auth auth.Authenticator

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
	defaultWorkspaceManagement    = "/faros.sh/workspaces"
)

// Init will inject all api server clients with post-start-hook
func (s *Service) Init(rest *rest.Config) error {
	p := newKubeConfigProxy(rest)

	kcpClient, err := kcpclient.NewClusterForConfig(rest)
	if err != nil {
		return err
	}

	farosClient, err := farosclient.NewClusterForConfig(rest)
	if err != nil {
		return err
	}

	coreClient, err := kubernetes.NewClusterForConfig(rest)
	if err != nil {
		return err
	}

	proxy := &httputil.ReverseProxy{
		Director:  p.director,
		Transport: roundtripper.RoundTripperFunc(p.roundTripper),
		//ErrorLog:  log.New(k.log.Writer(), "", 0),
	}

	authenticator, err := auth.NewAuthenticator(s.config, coreClient, farosClient, defaultOIDCCallbackPathPrefix)
	if err != nil {
		return err
	}

	s.proxy = proxy
	s.kcpClient = kcpClient
	s.farosClient = farosClient
	s.coreClients = coreClient
	s.authenticator = authenticator
	s.seeded = true
	return nil
}

func (s *Service) GetHandlers() []func(h http.Handler) http.HandlerFunc {
	return []func(h http.Handler) http.HandlerFunc{
		s.oidcCallback(),
		s.oidcLogin(),
		s.customTunnels(),
		s.workspacesHandler(),
	}
}
