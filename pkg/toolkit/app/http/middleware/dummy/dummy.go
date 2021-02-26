package dummy

import "net/http"

func Dummy(h http.Handler) http.Handler { return h }
