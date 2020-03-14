package datastoreinit

import (
	"errors"
	"strings"

	"github.com/mattrax/Mattrax/internal/datastore"
	"github.com/mattrax/Mattrax/internal/datastore/boltdb"
)

func Init(connUrl string) (datastore.Store, error) {
	connURL := strings.SplitN(connUrl, ":", 2)
	if len(connURL) != 2 {
		return nil, errors.New("invalid database connection url")
	}

	var store datastore.Store
	switch connURL[0] {
	case "boltdb":
		store = &boltdb.Store{}
	default:
		return nil, errors.New("unsupported datastore provider")
	}

	store.Init(connURL[1])

	return store, nil
}
