package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// --- FIX: CHECK INITIALIZATION ERRORS ---
    // Now initDB returns an error, so we must check it.
	if err := initDB(); err != nil { 
		// If DB connection fails, print error and exit process with failure code 1
		fmt.Fprintf(os.Stderr, "🚨 Critical startup failure: %v\n", err)
		os.Exit(1)
	}

	// loadJWTConfig does not return an error but should be called after initDB
	loadJWTConfig()

	// ----------------------------------------
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // local dev fallback
    }

    r := mux.NewRouter()

    // Home route
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        jsonResponse(w, map[string]string{"message": "Leave Management System API with MongoDB & JWT"})
    }).Methods("GET")

    // Employee routes
    r.HandleFunc("/employee/{id}", getEmployeeHandler).Methods("GET")
    r.HandleFunc("/employee/{id}/apply-leave", applyLeaveHandler).Methods("POST")
    r.HandleFunc("/employee/{id}/leaves", getEmployeeLeavesHandler).Methods("GET")

    // HR routes
    hr := r.PathPrefix("/hr").Subrouter()
    hr.HandleFunc("/login", hrLoginHandler).Methods("POST")

    // Protected HR routes
    hr.Handle("/pending-leaves", requireHRAuth(http.HandlerFunc(hrPendingLeavesHandler))).Methods("GET")
    hr.Handle("/leave/{index}/grant", requireHRAuth(http.HandlerFunc(hrGrantLeaveHandler))).Methods("POST")
    hr.Handle("/leave/{index}/decline", requireHRAuth(http.HandlerFunc(hrDeclineLeaveHandler))).Methods("POST")

    allowedOrigins := []string{"*"}
    headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
    origins := handlers.AllowedOrigins(allowedOrigins)

    // --- FIX: Dynamic Port Logging and Error Handling for ListenAndServe ---
    fmt.Printf("🚀 Server starting on 0.0.0.0:%s\n", port)
    
    // Bind to 0.0.0.0 and check for binding errors
    err := http.ListenAndServe("0.0.0.0:"+port, handlers.CORS(headers, methods, origins)(r))
    
    if err != nil {
        // This will now catch and log the specific error if the port is in use or the binding failed.
        fmt.Fprintf(os.Stderr, "🚨 HTTP server failed to start/run: %v\n", err)
        os.Exit(1)
    }
}