package boltdb

import (
	"fmt"
	"net/http"

	"github.com/mattrax/Mattrax/internal/authentication"
)

type Provider struct {
}

func (p *Provider) Init(config map[string]string) error {
	return nil
}

func (p *Provider) Handler(w http.ResponseWriter, r *http.Request) {
	settings := r.Context().Value("settings").(map[string]string)
	fmt.Println(settings)
	w.Write([]byte("AzureAD Provider!"))
}

func init() {
	authentication.MountProvider("azuread", &Provider{})
}
