// package gopostgres

package main

import (
	// "database/sql"
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"

	// _ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello, World!")

	r := router.Router()
	
	fmt.Println("Starting Server on POPRT 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}


