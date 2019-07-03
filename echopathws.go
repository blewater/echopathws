// Launches an echo path web server
// at the requested port
// or at the default 8770 port.
package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const shutdownSecondsAllowance = 15
const defaultHTTPPort = 8770

// Accepts a web listening port to launch the web server.
func main() {

	logger := log.New(os.Stdout, "", log.LstdFlags)

	var port int

	if len(os.Args) == 1 {

		// Default
		port = defaultHTTPPort

	} else {

		argport, err := strconv.Atoi(os.Args[1])
		if err != nil {
			logger.Fatalf("Not acceptable http listening port %v", os.Args[1])
		}
		port = argport
	}

	workflow(logger, port)

	logger.Printf("Http Server at container port %d completed its shutdown.\n", port)
}

// Performs steps for launching the web server.
func workflow(logger *log.Logger, port int) {

	httpServer := getHTTPServer(logger, port)

	setupTerminateSignal(logger, httpServer, port)

	launchHTTPListener(logger, httpServer, port)

}

// getHTTPServer constructs an HTTP listening server with 2 request handlers
func getHTTPServer(logger *log.Logger, port int) *http.Server {

	router := http.NewServeMux()

	// mute favicon requests
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {

	})

	// handler echoes the Path component of the request URL r.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		urlPath := r.URL.Path
		var respString string

		if len(urlPath) > 1 {
			respString = fmt.Sprintf("%q\n", html.EscapeString(urlPath[1:len(urlPath)]))
		} else {
			respString = fmt.Sprintf("Root Domain Request: /\n")
		}

		// -> client
		fmt.Fprint(w, respString)

		// -> stdout
		logger.Println(respString)
	}) // each request calls handler

	return &http.Server{
		Addr:     fmt.Sprintf(":%d", port),
		Handler:  router,
		ErrorLog: logger,
	}
}

// setupTerminateSignal connects the os.Interrupt signal to a quit channel to
// start teardown.
func setupTerminateSignal(logger *log.Logger, httpServer *http.Server, port int) {

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go httpServerShutdown(logger, httpServer, port, quit)

}

// Final step in launching an http server: Start accepting requests.
func launchHTTPListener(logger *log.Logger, httpServer *http.Server, port int) {

	logger.Printf("Http server at container port %d listening...\n", port)

	err := httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Can't launch http listener at container %d...\n", port)
	}

}

// httpServerShutdown handles the termination signal by shutting down the http server
// by closing connections and forcing shutdown if needed: "shutdownSecondsAllowance" max allowance.
func httpServerShutdown(logger *log.Logger, httpServer *http.Server, port int, quit <-chan os.Signal) {

	<-quit
	logger.Printf("Http server at container port %d is shutting down...\n", port)

	// Allow
	ctx, cancel := context.WithTimeout(context.Background(), shutdownSecondsAllowance*time.Second)
	defer cancel()

	httpServer.SetKeepAlivesEnabled(false)

	err := httpServer.Shutdown(ctx)
	if err != nil {
		logger.Fatalf("Could not shutdown the server @ %d. Error: %v\n", port, err)
	}
}
