package http

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	"github.com/mattrax/Mattrax/internal/mattrax"
	"github.com/rs/zerolog/log"
)

func Serve(srv *mattrax.Server) {
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(srv.Config.Port),
		Handler:      HeaderMiddleware(srv.Router, srv.Config.DevelopmentMode),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
		},
		// FUTURE: ErrorLog: {MAKE COMPATIBLE},
	}

	// Everything Good
	if srv.Config.DevelopmentMode {
		log.Warn().Msg("Development Mode Enabled. DO NOT use the server in production!")
	}
	log.Info().Str("domain", srv.Config.Domain).Int("port", srv.Config.Port).Msg("Initialised Mattrax MDM Server!")

	// TODO: Gracefull HTTP Shutdown
	if err := server.ListenAndServeTLS(srv.Config.CertFile, srv.Config.KeyFile); err != nil {
		log.Fatal().Int("port", srv.Config.Port).Str("certfile", srv.Config.CertFile).Str("keyfile", srv.Config.KeyFile).Err(err).Msg("Webserver error!")
	}
}

func HeaderMiddleware(handler http.Handler, devMode bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "MattraxMDM")
		if devMode {
			log.Debug().Str("method", r.Method).Str("path", r.URL.String()).Str("domain", r.Host).Str("remote", r.RemoteAddr).Str("user-agent", r.UserAgent()).Msg("request")
		}
		handler.ServeHTTP(w, r)
	})
}
