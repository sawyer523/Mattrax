package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/boltdb"
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

	server := mattrax.Server{
		Config: mattrax.Config{
			TenantName:    "Acme School Inc",
			PrimaryDomain: "mdm.otbeaumont.me",
		},
		UserService:   userService,
		PolicyService: policyService,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Mattrax Server. Created By Oscar Beaumont.")
	})

	// TODO: Gracefull Shutdown
	log.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServeTLS(":8000", "./certs/server.crt", "./certs/server.key", r))

	_ = server // TEMP
}
