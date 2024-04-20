package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/logger"
)

// Logger is a middleware to log HTTP request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}
		start := time.Now()
		defer func() {
			requestID, ok := r.Context().Value(ContextRequestID).(string)
			if !ok {
				requestID = ""
			}
			event := logger.Debug()
			if o.status >= http.StatusBadRequest && o.status != http.StatusPaymentRequired {
				event = logger.Info()
			}
			event.Str(
				"req-id", requestID,
			).Int(
				"status", o.status,
			).Str(
				"remote-addr", getRequestIP(r),
			).Str(
				"user-agent", r.UserAgent(),
			).Int64(
				"size", o.written,
			).Dur(
				"duration", time.Since(start),
			).Msgf("%s %s", r.Method, r.URL.Path)
		}()
		next.ServeHTTP(o, r)
	})
}

func getRequestIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if comma := strings.Index(ip, ","); comma != -1 {
		ip = ip[0:comma]
	}
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func (o *responseObserver) Flush() {
	flusher, ok := o.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
