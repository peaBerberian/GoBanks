package api

import "strconv"
import "log"
import "net/http"

const apiVersion = 1

// Start starts the API on localhost with the given port.
func Start(port int, certPath string, keyPath string) {
	var portStr = strconv.Itoa(port)

	var displayedVersion = "/v" + strconv.Itoa(apiVersion) + "/"
	mux := http.NewServeMux()
	mux.Handle(displayedVersion, http.HandlerFunc(handlerV1))

	// Start listing on a given port with these routes on this server.
	log.Print("Listening on port " + portStr + " ... ")
	err := http.ListenAndServeTLS(":"+portStr, certPath, keyPath, mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
