package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hascorp/hasmail/internal/hasmailtemplates/noop"
	"github.com/hascorp/hasmail/internal/hasmailtemplates/sample"
	"github.com/hascorp/hasmail/internal/healthcheck"
	"github.com/hascorp/hasmail/internal/middleware"
)

var (
	authFlag = flag.Bool("bypass", false, "denotes whether to skip API level authentication checks")
	portFlag = flag.Int("port", 8000, "port that the HTTP server listens on")
)

func main() {
	// TODO: add env var overrides for input flags to easily override
	//       when running in a container
	flag.Parse()
	start(*authFlag, *portFlag)
}

func start(bypassAuth bool, port int) {
	log.Printf("Starting with bypass=%v, port=%d", bypassAuth, port)

	if port < 1 {
		log.Fatalf("Invalid port: %d", port)
	}
	addr := fmt.Sprintf(":%d", port)

	r := mux.NewRouter()
	r.HandleFunc("/", healthcheck.PingHandler).Methods("GET")

	// sub routes for different email handlers, separated by
	// the template names
	s := r.PathPrefix("/mail").Subrouter()
	s.HandleFunc("/noop", noop.NoOpHandler).Methods("POST")
	s.HandleFunc("/sample", sample.SampleHandler).Methods("POST")

	http.Handle("/", r)

	r.Use(middleware.LoggingMiddleware)
	if !bypassAuth {
		s.Use(middleware.AuthVerifyMiddleware)
	}

	srv := &http.Server{
		Handler:           r,
		Addr:              addr,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Printf("Starting server on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
