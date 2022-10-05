package server

import (
	"context"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
	"github.com/faroshq/faros-hub/pkg/util/roundtripper"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

var (
	defaultTunnelPathPrefix = "/services/faros-tunnels"
)

type Tunnel interface {
	WithCustomTunnels(apiHandler http.Handler) http.HandlerFunc
	SeedClients(rest *rest.Config) error
}

type tunneler struct {
	// rest config for kcp syncer tunnels
	rest      *rest.Config
	kcpClient kcpclient.ClusterInterface
	proxy     *httputil.ReverseProxy

	seeded bool
}

const (
	kubeconfigTimeout = time.Hour * 24
)

type contextKey int

const (
	contextKeySyncer contextKey = iota
	contextKeyClusterName
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

	cf, err := utilkubernetes.NewClientFactory(rest)
	if err != nil {
		return err
	}

	client, err := cf.GetRootKCPClient()
	if err != nil {
		return err
	}

	proxy := &httputil.ReverseProxy{
		Director:  p.director,
		Transport: roundtripper.RoundTripperFunc(p.roundTripper),
		//ErrorLog:  log.New(k.log.Writer(), "", 0),
	}

	t.proxy = proxy
	t.rest = p.rest
	t.kcpClient = client
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
		if len(path) < 9 {
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}

		gv := accessv1alpha1.SchemeGroupVersion
		if path[2] != "apis" ||
			path[3] != gv.Group ||
			path[4] != gv.Version ||
			path[5] != "access" {
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}

		clusterName := path[1]
		requestName := path[6]

		klog.V(5).Infof("tunnel connection received", "clusterName", clusterName, "request", requestName)

		//access, err :=

		// TODO: This is very much not effective and should be cached
		syncer, err := t.kcpClient.Cluster(logicalcluster.New(clusterName)).WorkloadV1alpha1().SyncTargets().Get(ctx, requestName, metav1.GetOptions{})
		if err != nil {
			klog.V(2).Infof("failed to get cluster", err)
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, contextKeySyncer, syncer)
		ctx = context.WithValue(ctx, contextKeyClusterName, clusterName)
		r = r.WithContext(ctx)
		t.proxy.ServeHTTP(w, r)
	}
}
