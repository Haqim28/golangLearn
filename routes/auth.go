package routes

import (
	"golang/handlers"
	"golang/pkg/mysql"
	"golang/repositories"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerAuth(userRepository)

	r.HandleFunc("/register", h.Register).Methods("POST")
	// Create "/login" route using handler Login and method POST here ...
}
