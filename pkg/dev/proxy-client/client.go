package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/aojea/h2rev2"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/http2"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"k8s.io/utils/clock"
)

type Client struct {
	upstreamURL    string
	downstreamURL  string
	upstreamClient *http.Client
	tlsConfig      *tls.Config
}

func New(upstreamURL, downstreamURL, clientCertFile, clientCertKeyFile, servingCertFile string) (*Client, error) {
	certFile, err := ioutil.ReadFile(clientCertFile)
	if err != nil {
		return nil, err
	}

	clientCert, err := x509.ParseCertificate(certFile)
	if err != nil {
		return nil, err
	}

	servingFile, err := ioutil.ReadFile(servingCertFile)
	if err != nil {
		return nil, err
	}

	servingCert, err := x509.ParseCertificate(servingFile)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AddCert(clientCert)
	pool.AddCert(servingCert)

	keyFile, err := ioutil.ReadFile(clientCertKeyFile)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs: pool,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{
					certFile,
				},
				PrivateKey: key,
			},
		},
		ServerName: "proxy",
	}

	upstreamClient := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &Client{
		upstreamURL:    upstreamURL,
		downstreamURL:  downstreamURL,
		upstreamClient: upstreamClient,
		tlsConfig:      tlsConfig,
	}, nil
}

// startTunnel blocks until the context is cancelled trying to establish a tunnel against the specified target
func (c *Client) Run(ctx context.Context) {
	// connect to create the reverse tunnels
	var (
		initBackoff   = 5 * time.Second
		maxBackoff    = 5 * time.Minute
		resetDuration = 1 * time.Minute
		backoffFactor = 2.0
		jitter        = 1.0
		clock         = &clock.RealClock{}
		sliding       = true
	)

	backoffMgr := wait.NewExponentialBackoffManager(initBackoff, maxBackoff, resetDuration, backoffFactor, jitter, clock)
	logger := klog.FromContext(ctx)

	wait.BackoffUntil(func() {
		logger.V(2).Info("starting tunnel")
		err := c.startTunneler(ctx)
		if err != nil {
			logger.Error(err, "failed to create tunnel")
		}
	}, backoffMgr, sliding, ctx.Done())
}

func (c *Client) startTunneler(ctx context.Context) error {
	logger := klog.FromContext(ctx)

	logger = logger.WithValues("to", c.downstreamURL).WithValues("from", c.upstreamURL)
	logger.Info("connecting to destination URL")

	l, err := h2rev2.NewListener(c.upstreamClient, c.upstreamURL, "bob")
	if err != nil {
		panic(err)
	}

	// client --> local dev instance
	downstreamURL, err := url.Parse(c.downstreamURL)
	if err != nil {
		return err
	}

	// dev-proxy-server --> local dev instance
	proxy := httputil.NewSingleHostReverseProxy(downstreamURL)
	if err != nil {
		return err
	}

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		//req.URL.Path = strings.TrimPrefix(req.URL.Path, cmdTunnelProxy)
		director(req)
	}
	clientDownstream := http.DefaultClient
	proxy.Transport = clientDownstream.Transport

	spew.Dump("listener done")
	// reverse proxy the request coming from the reverse connection to the apiserver
	server := &http.Server{Handler: proxy}
	defer server.Close()

	logger.Info("serving on reverse connection")
	errCh := make(chan error)
	go func() {
		errCh <- server.Serve(l)
	}()

	select {
	case err = <-errCh:
	case <-ctx.Done():
		err = server.Close()
	}
	logger.V(2).Info("stop serving on reverse connection")
	return err
}
