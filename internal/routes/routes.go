package routes

import (
	"myapi/internal/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/itens", handlers.ListItens).Methods("GET")
	r.HandleFunc("/api/itens/{id}", handlers.GetItem).Methods("GET")
	r.HandleFunc("/api/itens/codigo/{codigo}", handlers.GetItemByCode).Methods("GET")
	r.HandleFunc("/api/itens", handlers.CreateItem).Methods("POST")
	r.HandleFunc("/api/itens", handlers.UpdateItem).Methods("PUT")
	r.HandleFunc("/api/itens/{id}", handlers.DeleteItem).Methods("DELETE")

	// Endpoints para Categorias
	r.HandleFunc("/categorias", handlers.ListCategoriasHandler).Methods("GET")
	r.HandleFunc("/categorias/get", handlers.GetCategoriaHandler).Methods("GET")
	r.HandleFunc("/categorias/create", handlers.CreateCategoriaHandler).Methods("POST")
	r.HandleFunc("/categorias/update", handlers.UpdateCategoriaHandler).Methods("PUT")
	r.HandleFunc("/categorias/delete", handlers.DeleteCategoriaHandler).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/docs", handlers.ScalarHandler).Methods("GET")

	return r
}