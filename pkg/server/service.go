package server

import (
	"context"
	"net/http"
	"path"
	"time"

	health "github.com/InVisionApp/go-health/v2"
	healthhandlers "github.com/InVisionApp/go-health/v2/handlers"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/server/auth"
	"github.com/faroshq/faros-hub/pkg/util/recover"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
)

var (
	scheme       = runtime.NewScheme()
	codecs       = serializer.NewCodecFactory(scheme)
	limit  int64 = 1024 * 1024 * 10
)

func init() {
	utilruntime.Must(tenancyv1alpha1.AddToScheme(scheme))
	utilruntime.Must(pluginsv1alpha1.AddToScheme(scheme))
}

var _ Interface = &Service{}

type Interface interface {
	Run(ctx context.Context) error
}

const (
	pathAPIVersion   = "/faros.sh/api/v1alpha1"
	pathWorkspaces   = "/workspaces"
	pathPlugins      = "/plugins"
	pathOIDC         = "/oidc"
	pathOIDCLogin    = "/oidc/login"
	pathOIDCCallback = "/oidc/callback"
)

type Service struct {
	config         *config.APIConfig
	authenticator  auth.Authenticator
	server         *http.Server
	router         *mux.Router
	health         *health.Health
	tenantsCluster logicalcluster.Name
	pluginsCluster logicalcluster.Name

	// tunneling tooling
	kcpClient   kcpclient.Interface
	farosClient farosclient.Interface
	coreClients kubernetes.Interface

	//proxy       *httputil.ReverseProxy
}

func New(config *config.APIConfig) (*Service, error) {
	//p := newKubeConfigProxy(config.RestConfig)

	kcpClient, err := kcpclient.NewForConfig(config.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	farosClient, err := farosclient.NewForConfig(config.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	coreClient, err := kubernetes.NewForConfig(config.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	//proxy := &httputil.ReverseProxy{
	//	Director:  p.director,
	//	Transport: roundtripper.RoundTripperFunc(p.roundTripper),
	//	//ErrorLog:  log.New(k.log.Writer(), "", 0),
	//}

	authenticator, err := auth.NewAuthenticator(config, coreClient, farosClient, path.Join(pathAPIVersion, pathOIDCCallback))
	if err != nil {
		return nil, err
	}

	s := &Service{
		config: config,
		//proxy:         proxy,
		tenantsCluster: logicalcluster.New(config.ControllersTenantWorkspace),
		pluginsCluster: logicalcluster.New(config.ControllersPluginsWorkspace),
		health:         health.New(),
		kcpClient:      kcpClient,
		farosClient:    farosClient,
		coreClients:    coreClient,
		authenticator:  authenticator,
	}

	s.router = setupRouter()
	apiRouter := s.router.PathPrefix(pathAPIVersion).Subrouter()
	apiRouter.HandleFunc("/healthz", healthhandlers.NewJSONHandlerFunc(s.health, nil)) // /healthz
	apiRouter.HandleFunc(pathOIDCLogin, s.oidcLogin)                                   // /faros.sh/api/v1alpha1/oidc/login
	apiRouter.HandleFunc(pathOIDCCallback, s.oidcCallback)                             // /faros.sh/api/v1alpha1/oidc/callback

	apiRouter.HandleFunc(pathWorkspaces, s.workspacesHandler).Methods(http.MethodGet)                              // /faros.sh/api/v1alpha1/workspaces
	apiRouter.HandleFunc(path.Join(pathWorkspaces, "{workspace}"), s.workspacesHandler).Methods(http.MethodGet)    // /faros.sh/api/v1alpha1/workspaces/{workspace}
	apiRouter.HandleFunc(path.Join(pathWorkspaces, "{workspace}"), s.workspacesHandler).Methods(http.MethodDelete) // /faros.sh/api/v1alpha1/workspaces/{workspace}
	apiRouter.HandleFunc(pathWorkspaces, s.workspacesHandler).Methods(http.MethodPost)                             // /faros.sh/api/v1alpha1/workspaces

	apiRouter.HandleFunc(pathPlugins, s.pluginsHandler).Methods(http.MethodGet) // /faros.sh/api/v1alpha1/plugins

	s.server = &http.Server{
		Addr: config.Addr,
		Handler: handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
		)(s),
	}

	return s, nil
}

func (s *Service) Run(ctx context.Context) error {
	klog.Info("Starting API Service")
	go func() {
		defer recover.Panic()
		<-ctx.Done()

		err := s.health.Stop()
		if err != nil {
			klog.Error(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		err = s.server.Shutdown(ctx)
		if err != nil {
			klog.Error("api shutdown error", zap.Error(err))
		}
		klog.Info("Stopped API Service")
	}()

	klog.Info("Server will now listen", "url", s.config.Addr)
	return s.server.ListenAndServe()
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(Panic())
	r.Use(Gzip())
	r.Use(Log())

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
