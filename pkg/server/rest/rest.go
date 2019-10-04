package rest

import "net/http"

func ServeFile(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}
