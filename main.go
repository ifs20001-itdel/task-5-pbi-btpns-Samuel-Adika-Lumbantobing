// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeypc/go-jwt-mux/controllers/authController"
	"github.com/jeypc/go-jwt-mux/controllers/photoController"
	"github.com/jeypc/go-jwt-mux/models"
)

func main() {
	models.ConnectDatabase()
	r := mux.NewRouter()

	// Authentication Endpoints
	r.HandleFunc("/users/register", authController.Register).Methods("POST")
	r.HandleFunc("/users/login", authController.Login).Methods("POST")
	r.HandleFunc("/users/{userID}", authController.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{userID}", authController.DeleteUser).Methods("DELETE")

	// Photos Endpoints
	r.HandleFunc("/photos", photoController.GetPhotos).Methods("GET")
	r.HandleFunc("/photos", photoController.CreatePhoto).Methods("POST")
	r.HandleFunc("/photos/{photoID}", photoController.UpdatePhoto).Methods("PUT")
	r.HandleFunc("/photos/{photoID}", photoController.DeletePhoto).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
