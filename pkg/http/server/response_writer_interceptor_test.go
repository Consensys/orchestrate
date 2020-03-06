package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockInterceptor struct {
	called bool
	header http.Header
}

func (i *MockInterceptor) Write(b []byte) (int, error) {
	i.called = true
	return 0, nil
}

func (i *MockInterceptor) Header() http.Header {
	return i.header
}

func TestResponseWriterInterceptor(t *testing.T) {
	rec := httptest.NewRecorder()
	interceptor := &MockInterceptor{header: make(http.Header)}
	rw := NewResponseWriterInterceptor(
		rec,
		func(code int, header http.Header) WriterInterceptor {
			if code == 400 {
				interceptor.header = header
				return interceptor
			}
			return nil
		},
	)
	rw.Header().Set("test-key", "test-value")
	rw.WriteHeader(400)
	_, _ = rw.Write([]byte(``))
	assert.True(t, interceptor.called, "Mock should have been called")
	assert.Equal(t, "test-value", interceptor.header.Get("test-key"), "Header should have been set")
}

func TestBytesBufferInterceptor(t *testing.T) {
	header := make(http.Header)
	i := NewBytesBufferInterceptor(header)
	assert.Equal(t, header, i.Header(), "Header should have been initialized correctly")
	_, err := i.Write([]byte("foo"))
	assert.NoError(t, err, "Write should not error")
	b, _ := i.ReadByte()
	assert.Equal(t, b, byte('f'), "Radbyte should return correct value")
}
