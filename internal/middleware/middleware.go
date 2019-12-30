package middleware

import "net/http"

import "github.com/rs/zerolog/log"

// Global is executed on every request and it:
// - Sets the default HTTP headers
// - Logs if in development mode
func Global(handler http.Handler, devMode bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "MattraxMDM")
		// TODO: Security HTTPS Headers

		if devMode {
			log.Debug().Str("method", r.Method).Str("path", r.URL.String()).Str("domain", r.Host).Str("remote", r.RemoteAddr).Str("user-agent", r.UserAgent()).Msg("request")
		}

		handler.ServeHTTP(w, r)
	})
}
