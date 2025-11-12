package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	initDB() 
	loadJWTConfig()

	r := mux.NewRouter()


	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, map[string]string{"message": "Leave Management System API with MongoDB & JWT"})
	}).Methods("GET")

	
	r.HandleFunc("/employee/{id}", getEmployeeHandler).Methods("GET")
	r.HandleFunc("/employee/{id}/apply-leave", applyLeaveHandler).Methods("POST")
	r.HandleFunc("/employee/{id}/leaves", getEmployeeLeavesHandler).Methods("GET")

	
	hr := r.PathPrefix("/hr").Subrouter()
	hr.HandleFunc("/login", hrLoginHandler).Methods("POST")


	hr.Handle("/pending-leaves", requireHRAuth(http.HandlerFunc(hrPendingLeavesHandler))).Methods("GET")
	hr.Handle("/leave/{index}/grant", requireHRAuth(http.HandlerFunc(hrGrantLeaveHandler))).Methods("POST")
	hr.Handle("/leave/{index}/decline", requireHRAuth(http.HandlerFunc(hrDeclineLeaveHandler))).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:5174"})

	fmt.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r))

}
