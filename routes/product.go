package routes

import (
	"golang/handlers"
	"golang/pkg/middleware"
	"golang/pkg/mysql"
	"golang/repositories"

	"github.com/gorilla/mux"
)

func ProductRoutes(r *mux.Router) {
	productRepository := repositories.RepositoryProduct(mysql.DB)
	h := handlers.HandlerProduct(productRepository)

	r.HandleFunc("/products", middleware.Auth(h.FindProducts)).Methods("GET")
	r.HandleFunc("/product/{id}", h.GetProduct).Methods("GET")
	// Create "/product" route using midleware Auth, handler CreateProduct, and method POST here ...
}
