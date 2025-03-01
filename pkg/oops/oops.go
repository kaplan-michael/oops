package oops

import (
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v3"
)

// ErrorInfo holds details for an error page.
type ErrorInfo struct {
	Code    int    `yaml:"code"`
	Title   string `yaml:"title"`
	Message string `yaml:"message"`
}

// ErrorsData is a container for a list of error definitions.
type ErrorsData struct {
	Errors []ErrorInfo `yaml:"errors"`
}

var (
	// tmpl holds the parsed HTML template.
	tmpl *template.Template
	// errorsMap maps error codes (e.g. 404) to ErrorInfo.
	errorsMap map[int]ErrorInfo
)

// loadTemplate reads the template file from disk and parses it.
func loadTemplate(filename string) (*template.Template, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return template.New("errorTemplate").Parse(string(data))
}

// loadErrors reads a YAML file and builds a map of error codes to ErrorInfo.
func loadErrors(filename string) (map[int]ErrorInfo, error) {
	var errorsData ErrorsData
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, &errorsData); err != nil {
		return nil, err
	}
	m := make(map[int]ErrorInfo)
	for _, e := range errorsData.Errors {
		m[e.Code] = e
	}
	return m, nil
}

// healthzHandler responds with a simple 200 OK "ok" message.
// It recovers from any panic to ensure that the health endpoint always responds.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s from %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, r.RemoteAddr)
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered in healthzHandler: %v", rec)
			http.Error(w, "unhealthy", http.StatusInternalServerError)
		}
	}()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Printf("Error writing healthz response: %v", err)
	}
}

// errorHandler extracts the error code from the URL (e.g. "/404"),
// looks up the corresponding error info, writes the proper status code,
// and renders the template.
func errorHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := r.Header.Get("X-Code")
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

	log.Printf("[%s] %s %s (X-Code: %s) from %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, codeStr, r.RemoteAddr)
	// Write the correct HTTP status code.
	w.WriteHeader(code)
	// Ensure we serve HTML.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", "oops")
	// Look up the error info from our map.
	errInfo, ok := errorsMap[code]
	if !ok {
		errInfo = ErrorInfo{
			Code:    code,
			Title:   "Error",
			Message: "An unexpected error occurred. Please try again later.",
		}
	}
	if err := tmpl.Execute(w, errInfo); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func Oops() {
	var err error
	// Load the HTML template from file.
	tmpl, err = loadTemplate("template.tmpl")
	if err != nil {
		log.Fatalf("Error loading template: %v", err)
	}
	// Load error definitions from the YAML file.
	errorsMap, err = loadErrors("errors.yaml")
	if err != nil {
		log.Fatalf("Error loading errors.yaml: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check endpoint.
	r.Get("/healthz", healthzHandler)

	// Handle error pages. Since the ingress forwards errors with headers,
	// we use a catch-all endpoint here.
	r.Get("/*", errorHandler)

	log.Println("Starting error service on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
