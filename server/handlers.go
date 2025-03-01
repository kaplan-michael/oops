package server

import (
	"github.com/kaplan-michael/oops/errorpage"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// errorHandler extracts the error code from headers,
// looks up the corresponding error info, writes the proper status code,
// and renders the template.
func errorHandler(ep *errorpage.ErrorPage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeStr := r.Header.Get("X-Code")
		code := parseCode(codeStr)

		zap.S().Infof("[%s] %s %s (X-Code: %s) from %s",
			time.Now().Format(time.RFC3339), r.Method, r.URL.Path, codeStr, r.RemoteAddr)

		rendered, err := ep.Render(code)
		if err != nil {
			zap.S().Warnf("Error rendering error page: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Write response headers and the rendered HTML.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Server", "server")
		w.WriteHeader(rendered.StatusCode)
		if _, err := w.Write(rendered.Body); err != nil {
			zap.S().Warnf("Error writing rendered error page: %v", err)
		}
	}
}

// healthzHandler responds with a simple 200 OK "ok" message.
// It recovers from any panic to ensure that the health endpoint always responds.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Infof("[%s] %s %s from %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, r.RemoteAddr)
	defer func() {
		if rec := recover(); rec != nil {
			zap.S().Errorf("Recovered in healthzHandler: %v", rec)
			http.Error(w, "unhealthy", http.StatusInternalServerError)
		}
	}()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		zap.S().Warnf("Error writing healthz response: %v", err)
	}
}

func parseCode(codeStr string) int {
	var code int
	if codeStr == "" {
		code = 404 // default if header is missing
	} else {
		var err error
		code, err = strconv.Atoi(codeStr)
		if err != nil {
			code = 404
		}
	}
	return code
}
