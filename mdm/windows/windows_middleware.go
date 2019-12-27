package windows

import (
	"log"
	"net/http"
)

func verifyUserAgent(userAgent string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != userAgent {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func verifySoapRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify Content-Type
		if r.Header.Get("Content-type") != "application/soap+xml; charset=utf-8" {
			log.Println("error: Invalid Content-Type '" + r.Header.Get("Content-type") + "'")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check for request body
		if r.ContentLength == 0 {
			log.Println("error: No Body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check request body isn't too large
		if r.ContentLength > 8e+6 /* 8MB */ { // TODO: Verify the 8e+6 works
			log.Println("error: Body too big")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// // Check request body type
		// var rBody []byte
		// if r.Body != nil {
		// 	var err error
		// 	if rBody, err = ioutil.ReadAll(r.Body); err != nil {
		// 		log.Println("error: Error restoring body content")
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		return
		// 	}
		// }
		// // Restore the io.ReadCloser to its original state
		// r.Body = ioutil.NopCloser(bytes.NewBuffer(rBody))
		// log.Println(string(rBody), http.DetectContentType(rBody))
		// if http.DetectContentType(rBody) != "text/xml; charset=utf-8" { // TODO: Not working with non-pretty printed text
		// 	log.Println("error: invalid request body content")
		// 	w.WriteHeader(http.StatusUnsupportedMediaType)
		// 	return
		// }

		// Parse to next handler
		next(w, r)
	}
}
