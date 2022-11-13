package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/aojea/h2rev2"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

var _ Interface = &Service{}

type Interface interface {
	Run(context.Context) error
	Shutdown(context.Context) error
}

type Service struct {
	addr     string
	server   *http.Server
	listener net.Listener
}

func New(addr, certFile, keyFile, clientCertFile string) (*Service, error) {
	b, err := ioutil.ReadFile(clientCertFile)
	if err != nil {
		return nil, err
	}

	clientCert, err := x509.ParseCertificate(b)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AddCert(clientCert)

	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	b, err = ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}

	l, err := tls.Listen("tcp", addr, &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{
					cert,
				},
				PrivateKey: key,
			},
		},
		//ClientCAs:  pool,
		//ClientAuth: tls.RequireAndVerifyClientCert,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		NextProtos:               []string{"h2"},
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
	})
	if err != nil {
		return nil, err
	}

	s := &Service{
		addr:     addr,
		listener: l,
	}

	revPool := h2rev2.NewReversePool()
	mux := http.NewServeMux()
	mux.Handle("/", revPool)

	//http.HandleFunc("/", s.Handler())
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	s.server = &server

	return s, nil
}

func (s *Service) Run(ctx context.Context) error {
	klog.V(2).Infof("Starting remote server on %s", s.addr)

	return s.server.Serve(s.listener)
}

func (s *Service) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "shutdown remote server")
	}
	return nil
}
