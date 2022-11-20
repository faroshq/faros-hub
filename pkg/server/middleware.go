package server

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/handlers"
	"go.uber.org/zap"
	"k8s.io/klog/v2"
)

func Panic() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					klog.Error("panic")
					klog.Error(e)

					klog.Error(string(debug.Stack()))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}

type logResponseWriter struct {
	http.ResponseWriter

	statusCode int
	path       string
	bytes      int
}

func (w *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *logResponseWriter) Flush() {
	flucher := w.ResponseWriter.(http.Flusher)
	flucher.Flush()
}

func (w *logResponseWriter) CloseNotify() <-chan bool {
	notify := w.ResponseWriter.(http.CloseNotifier)
	return notify.CloseNotify()
}

type logReadCloser struct {
	io.ReadCloser

	bytes int
}

func (rc *logReadCloser) Read(b []byte) (int, error) {
	n, err := rc.ReadCloser.Read(b)
	rc.bytes += n
	return n, err
}

func Log() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			ctx := r.Context()

			r.Body = &logReadCloser{ReadCloser: r.Body}
			w = &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK, path: r.URL.Path}

			logger := klog.FromContext(ctx)

			r = r.WithContext(ctx)

			logger = logger.WithValues(
				zap.String("request_method", r.Method),
				zap.String("request_path", r.URL.Path),
				zap.String("request_proto", r.Proto),
				zap.String("request_remote_addr", r.RemoteAddr),
				zap.String("request_user_agent", r.UserAgent()),
			)

			// enrich context with logger and add back to request
			ctx = klog.NewContext(ctx, logger)
			r = r.WithContext(ctx)

			defer func() {
				if shouldLog(w) {
					logger.WithValues(
						zap.Int("body_read_bytes", r.Body.(*logReadCloser).bytes),
						zap.Int("body_written_bytes", w.(*logResponseWriter).bytes),
						zap.Float64("duration", time.Since(t).Seconds()),
						zap.Int("response_status_code", w.(*logResponseWriter).statusCode),
					).V(4).Info("sent response")
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}

func Gzip() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.CompressHandler(h)
	}
}

// shouldLog defines if we should log request.
// Current rules:
// 1. If we returning an error (>=399) - check. Else not log
// 2. If 401 (unauthorized) - don't log
// 3. If request originated not from our up tp date agent/cli/ui - don't log
func shouldLog(w http.ResponseWriter) bool {
	statusCode := w.(*logResponseWriter).statusCode

	// TODO: Once agent checks are rollout and we don't see many default agent ("Go-http-client/2.0)
	// metrics - drop all these
	if statusCode >= 399 {
		// we don't log unauth as they are noisy
		if statusCode == 401 {
			return false
		}

		// currently we dont support non-synpse based agent logging.
		// this should filter all noise
		agent := w.Header().Get("Synpse-User-Agent")
		return agent != ""

	}
	return false
}
