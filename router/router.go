package router

import (
	"truecaller/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/callback", middleware.Callback).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/details", middleware.Details).Methods("GET", "OPTIONS")
	// https://api-stage.rupifi.com/auth/v1/tokens
	return router
}
