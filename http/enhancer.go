package http

import (
	"net/http"
)

// ServeMuxEnhancer are functions that enhance net/http Multiplexers
type ServeMuxEnhancer func(*http.ServeMux) *http.ServeMux
