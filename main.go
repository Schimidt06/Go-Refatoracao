package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"myapi/internal/config"
	"myapi/internal/handlers"
	model "myapi/internal/models"

	_ "myapi/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Refatoração API
// @version 1.0
// @description Esta é uma API de exemplo para o curso de Refatoração em Go.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	config.ConectaComBancoDeDados()

	// Endpoint raiz
	http.HandleFunc("/api", indexHandler)

	// Endpoints para Itens
	http.HandleFunc("/itens", handlers.ListItensHandler)                // GET para listar todos os itens
	http.HandleFunc("/itens/get", handlers.GetItenHandler)              // GET para buscar um item (espera id via query: ?id=1)
	http.HandleFunc("/itens/get-code", handlers.GetItenByCodigoHandler) // get-code?codigo=TEC001
	http.HandleFunc("/itens/create", handlers.CreateItenHandler)        // POST para criar um item
	http.HandleFunc("/itens/update", handlers.UpdateItenHandler)        // PUT para atualizar um item (JSON com id)
	http.HandleFunc("/itens/delete", handlers.DeleteItenHandler)        // DELETE para deletar um item (espera id via query: ?id=1)

	// Endpoints para Categorias
	http.HandleFunc("/categorias", listCategoriasHandler)         // GET para listar todas as categorias
	http.HandleFunc("/categorias/get", getCategoriaHandler)       // GET para buscar uma categoria (espera id via query)
	http.HandleFunc("/categorias/create", createCategoriaHandler) // POST para criar uma categoria
	http.HandleFunc("/categorias/update", updateCategoriaHandler) // PUT para atualizar uma categoria (JSON com id)
	http.HandleFunc("/categorias/delete", deleteCategoriaHandler) // DELETE para deletar uma categoria (espera id via query)

	// Swagger documentation
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Scalar UI
	http.HandleFunc("/docs", scalarHandler)

	log.Println("Servidor rodando na porta 8080")
	log.Println("Documentação disponível em http://localhost:8080/docs")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// scalarHandler serve a interface do Scalar
func scalarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/swagger/doc.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`)
}

// Handler raiz
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API Go!")
}

// ==================== HANDLERS PARA CATEGORIAS ====================

// Listar todas as categorias
// @Summary Listar todas as categorias
// @Description Retorna todas as categorias cadastradas no banco de dados
// @Tags categorias
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Categoria
// @Router /categorias [get]
func listCategoriasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var categorias []model.Categoria
	if err := config.DB.Find(&categorias).Error; err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categorias)
}

// Buscar uma única categoria pelo id (via query string: ?id=1)
// @Summary Buscar categoria por ID
// @Description Retorna uma única categoria baseada no ID fornecido via query string
// @Tags categorias
// @Accept  json
// @Produce  json
// @Param id query int true "ID da Categoria"
// @Success 200 {object} model.Categoria
// @Failure 400 {string} string "ID não fornecido ou inválido"
// @Failure 404 {string} string "Categoria não encontrada"
// @Router /categorias/get [get]
func getCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID não fornecido", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var categorias model.Categoria
	if err := config.DB.First(&categorias, id).Error; err != nil {
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(categorias)
}

// Criar uma nova categoria (envie JSON via POST)
// @Summary Criar uma nova categoria
// @Description Cria uma categoria baseada no JSON fornecido
// @Tags categorias
// @Accept  json
// @Produce  json
// @Param cat body model.Categoria true "Dados da Categoria"
// @Success 200 {object} model.Categoria
// @Failure 400 {string} string "Erro ao decodificar a categoria"
// @Failure 500 {string} string "Erro ao criar a categoria"
// @Router /categorias/create [post]
func createCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var categorias model.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categorias); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}
	if err := config.DB.Create(&categorias).Error; err != nil {
		http.Error(w, "Erro ao criar a categoria", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categorias)
}

// Atualizar uma categoria (envie JSON via PUT, com o campo id preenchido)
// @Summary Atualizar uma categoria
// @Description Atualiza os dados de uma categoria existente baseada no ID no JSON
// @Tags categorias
// @Accept  json
// @Produce  json
// @Param cat body model.Categoria true "Dados da Categoria (deve conter ID)"
// @Success 200 {object} model.Categoria
// @Failure 400 {string} string "Erro ao decodificar a categoria"
// @Failure 500 {string} string "Erro ao atualizar a categoria"
// @Router /categorias/update [put]
func updateCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var categorias model.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categorias); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}
	if err := config.DB.Save(&categorias).Error; err != nil {
		http.Error(w, "Erro ao atualizar a categoria", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categorias)
}

// Deletar uma categoria (via query string: ?id=1)
// @Summary Deletar uma categoria
// @Description Deleta uma categoria baseada no ID fornecido via query string
// @Tags categorias
// @Accept  json
// @Produce  plain
// @Param id query int true "ID da Categoria"
// @Success 200 {string} string "Categoria deletada com sucesso"
// @Failure 400 {string} string "ID não fornecido ou inválido"
// @Failure 500 {string} string "Erro ao deletar a categoria"
// @Router /categorias/delete [delete]
func deleteCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID não fornecido", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	if err := config.DB.Delete(&model.Categoria{}, id).Error; err != nil {
		http.Error(w, "Erro ao deletar a categoria", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Categoria deletada com sucesso"))
}
