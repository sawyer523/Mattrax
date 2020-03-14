package windows

import (
	"fmt"
	"net/http"
)

func (mdm *MDM) PolicyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POLICY")
		w.Write([]byte("Test!"))
	}
}
