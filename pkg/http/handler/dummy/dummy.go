package dummy

import (
	"net/http"
)

type Dummy struct{}

func (h *Dummy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {}
