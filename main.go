package main

import (
	"fmt"
	"log"
	"net/http"
	"test-crud-api/handler"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("Starting on localhost:8080...")

	r.HandleFunc("/api/student/{id}", handler.GetStudent).Methods("GET")
	r.HandleFunc("/api/student", handler.GetAllStudents).Methods("GET")
	r.HandleFunc("/api/newstudent", handler.CreateStudent).Methods("POST")
	r.HandleFunc("/api/student/{id}", handler.UpdateStudent).Methods("PUT")
	r.HandleFunc("/api/deletestudent/{id}", handler.DeleteStudent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
