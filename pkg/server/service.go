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

	//kcptenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	//workloadv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/workload/v1alpha1"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
	"go.uber.org/zap"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
)

func init() {
	utilruntime.Must(tenancyv1alpha1.AddToScheme(scheme))
}

var _ Interface = &Service{}

type Interface interface {
	Run(ctx context.Context) error
}

const (
	pathAPIVersion   = "/faros.sh/api/v1alpha1"
	pathWorkspaces   = "/workspaces"
	pathOIDC         = "/oidc"
	pathOIDCLogin    = "/oidc/login"
	pathOIDCCallback = "/oidc/callback"
)

type Service struct {
	config        *config.APIConfig
	authenticator auth.Authenticator
	server        *http.Server
	router        *mux.Router
	health        *health.Health
	cluster       logicalcluster.Name

	// tunneling tooling
	kcpClient   kcpclient.ClusterInterface
	farosClient farosclient.ClusterInterface
	coreClients kubernetes.ClusterInterface

	//proxy       *httputil.ReverseProxy
}

func New(config *config.APIConfig) (*Service, error) {
	//p := newKubeConfigProxy(config.RestConfig)

	kcpClient, err := kcpclient.NewClusterForConfig(config.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	farosClient, err := farosclient.NewClusterForConfig(config.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	coreClient, err := kubernetes.NewClusterForConfig(config.KCPClusterRestConfig)
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
		cluster:       logicalcluster.New(config.ControllersTenantWorkspace),
		health:        health.New(),
		kcpClient:     kcpClient,
		farosClient:   farosClient,
		coreClients:   coreClient,
		authenticator: authenticator,
	}

	s.router = setupRouter()
	apiRouter := s.router.PathPrefix(pathAPIVersion).Subrouter()
	apiRouter.HandleFunc("/healthz", healthhandlers.NewJSONHandlerFunc(s.health, nil))
	apiRouter.HandleFunc(pathOIDCLogin, s.oidcLogin)
	apiRouter.HandleFunc(pathOIDCCallback, s.oidcCallback)

	apiRouter.HandleFunc(pathWorkspaces, s.workspacesHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc(path.Join(pathWorkspaces, "{workspace}"), s.workspacesHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc(path.Join(pathWorkspaces, "{workspace}"), s.workspacesHandler).Methods(http.MethodDelete)
	apiRouter.HandleFunc(pathWorkspaces, s.workspacesHandler).Methods(http.MethodPost)

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
