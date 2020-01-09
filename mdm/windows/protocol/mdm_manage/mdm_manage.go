package mdmmanage

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Make work

		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}
}
