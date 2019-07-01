// Package ws launches a web server
// for an echo path web server
// at the requested port
// or at the default 8770 port.
package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Accepts a web listening port to launch the web server.
func main() {

	var port int

	if len(os.Args) == 1 {
		port = 8770

	} else {

		argport, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Not acceptable web listening port %v", os.Args[1])
		}
		port = argport
	}
	workflow(port)
}

// Performs 2 steps for the web server.
func workflow(port int) {

	setHTTPHandler(port)

	launchHTTPListener(port)

}

func setHTTPHandler(port int) {

	fmt.Printf("Starting localhost:%d...\n", port)

	http.HandleFunc("/", handler) // each request calls handler

	fmt.Printf("Set HTTP handler @ localhost:%d...\n", port)

}

func launchHTTPListener(port int) {

	err := http.ListenAndServe(fmt.Sprint(":", port), nil)
	if err != nil {
		log.Fatalf("Can't launch http listener at %d...\n", port)
	}

}

// handler echoes the Path component of the request URL r.
func handler(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path
	var respString string
	if len(urlPath) > 1 {
		respString = fmt.Sprintf("%q\n", html.EscapeString(urlPath[1:len(urlPath)]))
	} else {
		respString = fmt.Sprintf("Root Domain Request: /\n")
	}
	// -> client
	fmt.Fprint(w, respString)
}
