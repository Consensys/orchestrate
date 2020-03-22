package httputil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
)

type WriterInterceptor interface {
	io.Writer
	Header() http.Header
}

type ResponseWriterInterceptor struct {
	rw           http.ResponseWriter
	interceptors func(int, http.Header) WriterInterceptor

	header      http.Header
	interceptor WriterInterceptor

	closeNotifyCh chan bool
}

func NewResponseWriterInterceptor(rw http.ResponseWriter, interceptors func(int, http.Header) WriterInterceptor) *ResponseWriterInterceptor {
	return &ResponseWriterInterceptor{
		rw:            rw,
		interceptors:  interceptors,
		closeNotifyCh: make(chan bool, 1),
	}
}

func (i *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	i.interceptor = i.interceptors(statusCode, i.header)
	if i.interceptor == nil {
		// Copy headers in original ResponseWriter
		for key := range i.header {
			i.rw.Header().Set(key, i.header.Get(key))
		}

		// Use original writer to write headers
		i.rw.WriteHeader(statusCode)
	}
}

func (i *ResponseWriterInterceptor) Write(b []byte) (int, error) {
	if i.interceptor != nil {
		// If intercepted we redirect bytes to the response interceptor
		return i.interceptor.Write(b)
	}
	return i.rw.Write(b)
}

func (i *ResponseWriterInterceptor) Header() http.Header {
	if i.header == nil {
		i.header = make(http.Header)
	}
	return i.header
}

func (i *ResponseWriterInterceptor) Interceptor() WriterInterceptor {
	return i.interceptor
}

func (i *ResponseWriterInterceptor) CloseNotify() <-chan bool {
	// This will panic if rw is not an http.CloseNotifier
	return i.rw.(http.CloseNotifier).CloseNotify() //nolint
}

func (i *ResponseWriterInterceptor) Flush() {
	if f, ok := i.rw.(http.Flusher); ok {
		f.Flush()
	}
}

func (i *ResponseWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := i.rw.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("not a hijacker: %T", i.rw)
}

type BytesBufferInterceptor struct {
	bytes.Buffer
	header http.Header
}

func NewBytesBufferInterceptor(header http.Header) *BytesBufferInterceptor {
	return &BytesBufferInterceptor{
		Buffer: bytes.Buffer{},
		header: header,
	}
}

func (i *BytesBufferInterceptor) Header() http.Header {
	return i.header
}
