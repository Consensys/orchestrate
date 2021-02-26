package httputil

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type WriterRecorder interface {
	http.ResponseWriter
	CloseNotify() <-chan bool
	GetCode() int
}

func NewResponseWriterRecorder(rw http.ResponseWriter) WriterRecorder {
	return &responseWriterRecorder{
		ResponseWriter: rw,
		statusCode:     http.StatusOK,
		closeNotifyCh:  make(chan bool, 1),
	}
}

// ResponseWriterRecorder captures information from the response and preserves it for
// later analysis.
type responseWriterRecorder struct {
	http.ResponseWriter
	statusCode    int
	closeNotifyCh chan bool
}

func (rec *responseWriterRecorder) GetCode() int {
	return rec.statusCode
}

// WriteHeader captures the status code for later retrieval.
func (rec *responseWriterRecorder) WriteHeader(status int) {
	rec.ResponseWriter.WriteHeader(status)
	rec.statusCode = status
}

// Hijack hijacks the connection
func (rec *responseWriterRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rec.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("not a hijacker: %T", rec.ResponseWriter)
}

// Flush sends any buffered data to the client.
func (rec *responseWriterRecorder) Flush() {
	if f, ok := rec.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (rec *responseWriterRecorder) CloseNotify() <-chan bool {
	// This will panic if rw is not an http.CloseNotifier
	if rw2, ok := rec.ResponseWriter.(http.CloseNotifier); ok { //nolint
		return rw2.CloseNotify()
	}

	return rec.closeNotifyCh
}
