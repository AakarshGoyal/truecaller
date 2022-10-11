package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"truecaller/router"
)

func main() {
	r := router.Router()
	os.Setenv("PORT", "7358")
	port := os.Getenv("PORT")
	fmt.Println("TrueCaller Login with Rupifi!")
	fmt.Printf("Starting server on the port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
