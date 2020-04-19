package httputil

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type WriterRecorder interface {
	http.ResponseWriter
	GetCode() int
}

func NewResponseWriterRecorder(rw http.ResponseWriter) WriterRecorder {
	return &ResponseWriterRecorder{
		ResponseWriter: rw,
		statusCode:     http.StatusOK,
	}
}

// ResponseWriterRecorder captures information from the response and preserves it for
// later analysis.
type ResponseWriterRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *ResponseWriterRecorder) GetCode() int {
	return rec.statusCode
}

// WriteHeader captures the status code for later retrieval.
func (rec *ResponseWriterRecorder) WriteHeader(status int) {
	rec.ResponseWriter.WriteHeader(status)
	rec.statusCode = status
}

// Hijack hijacks the connection
func (rec *ResponseWriterRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rec.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("not a hijacker: %T", rec.ResponseWriter)
}

// Flush sends any buffered data to the client.
func (rec *ResponseWriterRecorder) Flush() {
	if f, ok := rec.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
