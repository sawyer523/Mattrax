package windows

import "net/http"

func (mdm *MDM) ManageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test!"))
	}
}
