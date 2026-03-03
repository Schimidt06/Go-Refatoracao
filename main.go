package main

import (
	"log"
	"net/http"

	"myapi/internal/config"
	"myapi/internal/routes"

	_ "myapi/docs"
)

func main() {
	config.ConnectDatabase()

	r := routes.SetupRoutes()

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
