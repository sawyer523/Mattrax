package main

import (
	"crypto/x509/pkix"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/api"
	"github.com/mattrax/Mattrax/internal/boltdb"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/mattrax/Mattrax/mdm/windows"
)

func main() {
	db, err := boltdb.Init()
	if err != nil {
		panic(err) // TODO
	}
	defer db.Close()

	userService, err := boltdb.NewUserService(db)
	if err != nil {
		panic(err) // TODO
	}

	policyService, err := boltdb.NewPolicyService(db)
	if err != nil {
		panic(err) // TODO
	}

	settingsService, err := boltdb.NewSettingsService(db)
	if err != nil {
		panic(err) // TODO
	}

	certificateService, err := boltdb.NewCertificateService(db, types.IdentityCertificateConfig{
		KeyLength: 4096,
		Subject: pkix.Name{
			// TODO: Does Configuring it work?
			Country:            []string{"US"},
			Organization:       []string{"groob-io"},
			OrganizationalUnit: []string{"SCEP CA"},
		},
	})
	if err != nil {
		panic(err) // TODO
	}

	server := mattrax.Server{
		Config: mattrax.Config{
			Port:                   443,
			PrimaryDomain:          "mdm.otbeaumont.me",
			WindowsDiscoveryDomain: "enterpriseenrollment.otbeaumont.me",
			CertFile:               "./certs/cert.pem",
			KeyFile:                "./certs/privkey.pem",
			DevelopmentMode:        true,
		},
		UserService:        userService,
		PolicyService:      policyService,
		SettingsService:    settingsService,
		CertificateService: certificateService,
	}

	r := mux.NewRouter()

	windows.Init(server, r)

	api.InitAPI(server, r)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Mattrax Server. Created By Oscar Beaumont.")
	})

	// TODO: Gracefull Shutdown
	port := strconv.Itoa(server.Config.Port)
	log.Println("Listening on port " + port + "...")
	log.Fatal(http.ListenAndServeTLS(":"+port, server.Config.CertFile, server.Config.KeyFile, logRequest(r)))
}

// TEMP
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
