package errorpage

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"sync"
)

// ErrorInfo holds details for an error code.
type ErrorInfo struct {
	Code    int    `yaml:"code"`
	Title   string `yaml:"title"`
	Message string `yaml:"message"`
}

// ErrorsData is used for unmarshaling the YAML file.
type ErrorsData struct {
	Errors []ErrorInfo `yaml:"errors"`
}

// RenderedError holds the HTTP status code and the rendered HTML body.
type RenderedError struct {
	StatusCode int
	Body       []byte
}

// ErrorPage encapsulates the templating logic and error definitions.
type ErrorPage struct {
	tmpl      *template.Template
	errorsMap map[int]ErrorInfo
}

var (
	instance *ErrorPage
	once     sync.Once
)

// NewErrorPage returns the singleton instance of ErrorPage. It loads the template and YAML file only once.
// Subsequent calls return the same pointer.
func NewErrorPage(templateFile, errorsFile string) *ErrorPage {
	once.Do(func() {
		ep := &ErrorPage{}
		var err error
		ep.tmpl, err = loadTemplate(templateFile)
		if err != nil {
			log.Fatalf("Failed to load template: %v", err)
		}

		ep.errorsMap, err = loadErrors(errorsFile)
		if err != nil {
			log.Fatalf("Failed to load errors file: %v", err)
		}
		instance = ep
	})
	return instance
}

// loadTemplate reads and parses the HTML template file.
func loadTemplate(filename string) (*template.Template, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return template.New("errorTemplate").Parse(string(data))
}

// loadErrors reads the YAML file and builds a map of error codes to ErrorInfo.
// It returns an error if any error code is 0 (assuming 0 is invalid).
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
		if e.Code == 0 {
			return nil, fmt.Errorf("invalid error code for error %+v: must be a non-zero integer", e)
		}
		m[e.Code] = e
	}
	return m, nil
}

// Render renders the error page for the given HTTP status code and returns a pointer to a RenderedError.
func (ep *ErrorPage) Render(code int) (*RenderedError, error) {
	// Look up error info; if not found, use a default.
	errInfo, ok := ep.errorsMap[code]
	if !ok {
		errInfo = ErrorInfo{
			Code:    code,
			Title:   "Error",
			Message: "An unexpected error occurred. Please try again later.",
		}
	}

	// Render the template into a buffer.
	var buf bytes.Buffer
	if err := ep.tmpl.Execute(&buf, errInfo); err != nil {
		return nil, err
	}

	return &RenderedError{
		StatusCode: code,
		Body:       buf.Bytes(),
	}, nil
}
