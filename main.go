// main.go
package main

import (
	"log"
	"net/http"

	"github.com/jeypc/go-jwt-mux/router"
)

func main() {
	r := router.SetupRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}
