package windows

import (
	"net/http"
)

func defaultHeaders(userAgent string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "-1")

		next(w, r)
	}
}
