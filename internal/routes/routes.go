package routes

import (
	"myapi/internal/handlers"
	"myapi/internal/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Global Middleware
	r.Use(middleware.JsonContentType)

	// Item Routes
	ItemRoutes(r)

	// Categoria Routes
	CategoriaRoutes(r)

	// Swagger and Docs (Not using JsonContentType middleware explicitly here, 
	// but r.Use applies to all sub-routes unless bypassed)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/docs", handlers.ScalarHandler).Methods("GET")

	return r
}