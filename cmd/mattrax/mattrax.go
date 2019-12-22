package main

import (
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

	server := mattrax.Server{
		Config: mattrax.Config{
			TenantName:    "Acme School Inc",
			PrimaryDomain: "mdm.otbeaumont.me",
		},
		UserService: userService,
	}
}
