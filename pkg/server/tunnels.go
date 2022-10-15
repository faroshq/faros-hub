package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/util/roundtripper"
)

var (
	defaultTunnelPathPrefix = "/services/faros-tunnels"
)

type Tunnel interface {
	WithCustomTunnels(apiHandler http.Handler) http.HandlerFunc
	SeedClients(rest *rest.Config) error
}

type tunneler struct {
	kcpClient   kcpclient.ClusterInterface
	farosClient farosclient.ClusterInterface
	coreClients kubernetes.ClusterInterface
	proxy       *httputil.ReverseProxy

	seeded bool
}

const (
	kubeconfigTimeout = time.Hour * 24
)

type contextKey int

const (
	contextKeyForwardPath contextKey = iota
	contextKeyResponse
	contextKeyClient
)

// NewTunneler creates a new tunnel handler without initializing the clients
func NewTunneler() Tunnel {
	return &tunneler{}
}

// SeedClients will inject all api server clients with post-start-hook
func (t *tunneler) SeedClients(rest *rest.Config) error {
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

	t.proxy = proxy
	t.kcpClient = kcpClient
	t.farosClient = farosClient
	t.coreClients = coreClient
	t.seeded = true
	return nil
}

// HTTP Handler that handles reverse connections and reverse proxy requests using 2 different paths:
//
// https://host/services/faros-tunnels ...
func (t *tunneler) WithCustomTunnels(apiHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fall through, faros tunnels URL start by /services/faros-tunnels
		if !strings.HasPrefix(r.URL.Path, defaultTunnelPathPrefix) {
			apiHandler.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		// route the request
		p := strings.TrimPrefix(r.URL.Path, defaultTunnelPathPrefix)
		path := strings.Split(strings.Trim(p, "/"), "/")
		if len(path) < 11 {
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}

		gv := accessv1alpha1.SchemeGroupVersion
		if path[2] != "apis" ||
			path[3] != gv.Group ||
			path[4] != gv.Version ||
			path[7] != "access" {
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}

		clusterName := path[1]
		requestNamespace := path[6]
		requestName := path[8]
		command := path[10:]

		klog.V(5).Infof("tunnel connection received", "clusterName", clusterName, "namespace", requestNamespace, "request", requestName)

		access, err := t.farosClient.Cluster(logicalcluster.New(clusterName)).AccessV1alpha1().Requests(requestNamespace).Get(ctx, requestName, metav1.GetOptions{})
		if err != nil {
			klog.V(2).Infof("failed to get access", err)
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}

		secret, err := t.coreClients.Cluster(logicalcluster.New(clusterName)).CoreV1().Secrets(access.Namespace).Get(ctx, access.Name, metav1.GetOptions{})
		if err != nil {
			klog.V(2).Infof("failed to get access", err)
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}

		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			klog.V(2).Infof("header does not contain Bearer token")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")

		switch {
		case token != string(secret.Data["token"]):
			klog.V(2).Infof("token does not match")
			w.WriteHeader(http.StatusForbidden)
			return
		case token == "":
			klog.V(2).Infof("token is empty")
			w.WriteHeader(http.StatusForbidden)
			return
		case token == string(secret.Data["token"]):
			// ok
		default:
			klog.V(2).Infof("token error")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		forwardPath := fmt.Sprintf("/services/syncer-tunnels/clusters/%s/apis/workload.kcp.dev/v1alpha1/synctargets/%s/proxy/%s", clusterName, access.Spec.ClusterName, strings.Join(command, "/"))

		ctx = context.WithValue(ctx, contextKeyForwardPath, forwardPath)
		r = r.WithContext(ctx)
		t.proxy.ServeHTTP(w, r)
	}
}
