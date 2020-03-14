package authentication

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattrax/Mattrax/internal/mattrax/settings"
)

var Providers = map[string]Provider{}

type Provider interface {
	Init(config map[string]string) error
	Handler(w http.ResponseWriter, r *http.Request)
}

func MountProvider(name string, p Provider) {
	Providers[name] = p
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.String())
		w.Write([]byte("Auth Entrypoint"))
	}
}

func ProviderHandler(ss *settings.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName, ok := mux.Vars(r)["provider"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Auth provider not specified!"))
			return
		}
		provider, ok := Providers[providerName]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Auth provider not found!"))
			return
		}

		settings, ok := ss.Get().AuthProviderSettings[providerName]
		if ok {
			ctx := context.WithValue(context.Background(), "settings", settings)
			r = r.WithContext(ctx)
		}

		provider.Handler(w, r)
	}
}
