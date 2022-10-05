package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/faroshq/faros-hub/pkg/util/clientcache"
	"github.com/faroshq/faros-hub/pkg/util/responsewriter"
	workloadv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/workload/v1alpha1"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

type kubeConfigProxy struct {
	rest *rest.Config
	// clientCache is a cache for direct client connections
	clientCache clientcache.ClientCache
}

func newKubeConfigProxy(rest *rest.Config) *kubeConfigProxy {
	return &kubeConfigProxy{
		rest:        rest,
		clientCache: clientcache.NewClientCache(time.Hour),
	}
}

// director is called by the ReverseProxy. It converts an incoming request into
// the one that'll go out to the API server. It also resolves an HTTP client
// that will be able to make the ongoing request.
//
// Unfortunately the signature of httputil.ReverseProxy.Director does not allow
// us to return values.  We get around this limitation slightly naughtily by
// storing return information in the request context.
func (k *kubeConfigProxy) director(r *http.Request) {
	ctx := r.Context()

	syncer, ok := ctx.Value(contextKeySyncer).(*workloadv1alpha1.SyncTarget)
	if !ok {
		k.error(r, http.StatusInternalServerError, fmt.Errorf("no syncer in context"))
		return
	}
	if syncer == nil {
		k.error(r, http.StatusForbidden, nil)
		return
	}

	clusterName, ok := ctx.Value(contextKeyClusterName).(string)
	if !ok {
		k.error(r, http.StatusInternalServerError, fmt.Errorf("no cluster name in context"))
		return
	}

	key := struct {
		namespace string
		name      string
		cluster   string
	}{
		syncer.Namespace,
		syncer.Name,
		clusterName,
	}

	cli := k.clientCache.Get(key)
	if cli == nil {
		var err error
		cli, err = k.cli(ctx)
		if err != nil {
			k.error(r, http.StatusInternalServerError, err)
			return
		}

		k.clientCache.Put(key, cli)
	}

	r.RequestURI = ""
	r.URL.Scheme = "https"
	r.URL.Host = strings.ReplaceAll(k.rest.Host, "https://", "")

	// We basically piggy-back on syncer tunnel from KCP:
	// https://localhost:6443/services/faros-tunnels/clusters/<ws>/apis/access.faros.sh/v1alpha1/access/<name>/proxy ->
	// https://localhost:6443/services/syncer-tunnels/clusters/<ws>/apis/workload.kcp.dev/v1alpha1/synctargets/<name>/proxy

	r.URL.Path = strings.Replace(r.URL.Path, "faros-tunnels", "syncer-tunnels", 1)
	r.URL.Path = strings.Replace(r.URL.Path, "access.faros.sh/v1alpha1/clusters", "workload.kcp.dev/v1alpha1/synctargets", 1)

	r.Header.Del("Authorization")
	r.Host = r.URL.Host

	// http.Request.WithContext returns a copy of the original Request with the
	// new context, but we have no way to return it, so we overwrite our
	// existing request.
	*r = *r.WithContext(context.WithValue(ctx, contextKeyClient, cli))

}

// cli returns an appropriately configured HTTP client for forwarding the
// incoming request to a cluster
func (k *kubeConfigProxy) cli(ctx context.Context) (*http.Client, error) {
	return rest.HTTPClientFor(k.rest)
}

// roundTripper is called by ReverseProxy to make the onward request happen.  We
// check if we had an error earlier and return that if we did. Otherwise we dig
// out the client and call it.
func (k *kubeConfigProxy) roundTripper(r *http.Request) (*http.Response, error) {
	if resp, ok := r.Context().Value(contextKeyResponse).(*http.Response); ok {
		return resp, nil
	}

	cli := r.Context().Value(contextKeyClient).(*http.Client)
	resp, err := cli.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusSwitchingProtocols {
		resp.Body = newCancelBody(resp.Body.(io.ReadWriteCloser), kubeconfigTimeout)
	}

	return resp, err
}

func (k *kubeConfigProxy) error(r *http.Request, statusCode int, err error) {
	if err != nil {
		klog.V(5).Info(err)
	}

	w := responsewriter.New(r)
	http.Error(w, http.StatusText(statusCode), statusCode)

	*r = *r.WithContext(context.WithValue(r.Context(), contextKeyResponse, w.Response()))
}

// cancelBody is a workaround for the fact that http timeouts are incompatible
// with hijacked connections. See: https://github.com/golang/go/issues/31391:
type cancelBody struct {
	io.ReadWriteCloser
	t *time.Timer
	c chan struct{}
}

func (b *cancelBody) wait() {
	select {
	case <-b.t.C:
		b.ReadWriteCloser.Close()
	case <-b.c:
		b.t.Stop()
	}
}

func (b *cancelBody) Close() error {
	select {
	case b.c <- struct{}{}:
	default:
	}

	return b.ReadWriteCloser.Close()
}

func newCancelBody(rwc io.ReadWriteCloser, d time.Duration) io.ReadWriteCloser {
	b := &cancelBody{
		ReadWriteCloser: rwc,
		t:               time.NewTimer(d),
		c:               make(chan struct{}),
	}

	go b.wait()

	return b
}
