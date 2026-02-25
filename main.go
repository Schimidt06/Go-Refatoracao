package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"myapi/internal/config"
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
	http.HandleFunc("/itens", listItensHandler)                // GET para listar todos os itens
	http.HandleFunc("/itens/get", getItenHandler)              // GET para buscar um item (espera id via query: ?id=1)
	http.HandleFunc("/itens/get-code", getItenByCodigoHandler) // get-code?codigo=TEC001
	http.HandleFunc("/itens/create", createItenHandler)        // POST para criar um item
	http.HandleFunc("/itens/update", updateItenHandler)        // PUT para atualizar um item (JSON com id)
	http.HandleFunc("/itens/delete", deleteItenHandler)        // DELETE para deletar um item (espera id via query: ?id=1)

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

// ==================== HANDLERS PARA ITENS ====================

// Listar todos os itens
// @Summary Listar todos os itens
// @Description Retorna todos os itens cadastrados no banco de dados
// @Tags itens
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Iten
// @Router /itens [get]
func listItensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var itens []model.Iten
	if err := config.DB.Find(&itens).Error; err != nil {
		http.Error(w, "Erro ao buscar itens", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(itens)
}

// Buscar um único item pelo id (via query string: ?id=1)
// @Summary Buscar item por ID
// @Description Retorna um único item baseado no ID fornecido via query string
// @Tags itens
// @Accept  json
// @Produce  json
// @Param id query int true "ID do Item"
// @Success 200 {object} model.Iten
// @Failure 400 {string} string "ID não fornecido ou inválido"
// @Failure 404 {string} string "Item não encontrado"
// @Router /itens/get [get]
func getItenHandler(w http.ResponseWriter, r *http.Request) {
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
	var item model.Iten
	if err := config.DB.First(&item, id).Error; err != nil {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

// Buscar um item pelo campo "codigo"
// @Summary Buscar item por código
// @Description Retorna um único item baseado no código fornecido via query string
// @Tags itens
// @Accept  json
// @Produce  json
// @Param codigo query string true "Código do Item"
// @Success 200 {object} model.Iten
// @Failure 400 {string} string "Código não fornecido"
// @Failure 404 {string} string "Item não encontrado"
// @Router /itens/get-code [get]
func getItenByCodigoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	cod := r.URL.Query().Get("codigo")
	if cod == "" {
		http.Error(w, "Código não fornecido", http.StatusBadRequest)
		return
	}
	var item model.Iten
	// Busca o item onde o campo "codigo" é igual ao valor fornecido
	if err := config.DB.Where("codigo = ?", cod).First(&item).Error; err != nil {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

// Criar um novo item (envie JSON via POST)
// @Summary Criar um novo item
// @Description Cria um item baseado no JSON fornecido
// @Tags itens
// @Accept  json
// @Produce  json
// @Param item body model.Iten true "Dados do Item"
// @Success 200 {object} model.Iten
// @Failure 400 {string} string "Erro ao decodificar o item"
// @Failure 500 {string} string "Erro ao criar o item"
// @Router /itens/create [post]
func createItenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var item model.Iten
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Erro ao decodificar o item", http.StatusBadRequest)
		return
	}
	if err := config.DB.Create(&item).Error; err != nil {
		http.Error(w, "Erro ao criar o item", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

// Atualizar um item (envie JSON via PUT, com o campo id preenchido)
// @Summary Atualizar um item
// @Description Atualiza os dados de um item existente baseado no ID no JSON
// @Tags itens
// @Accept  json
// @Produce  json
// @Param item body model.Iten true "Dados do Item (deve conter ID)"
// @Success 200 {object} model.Iten
// @Failure 400 {string} string "Erro ao decodificar o item"
// @Failure 500 {string} string "Erro ao atualizar o item"
// @Router /itens/update [put]
func updateItenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var item model.Iten
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Erro ao decodificar o item", http.StatusBadRequest)
		return
	}
	if err := config.DB.Save(&item).Error; err != nil {
		http.Error(w, "Erro ao atualizar o item", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

// Deletar um item (via query string: ?id=1)
// @Summary Deletar um item
// @Description Deleta um item baseado no ID fornecido via query string
// @Tags itens
// @Accept  json
// @Produce  plain
// @Param id query int true "ID do Item"
// @Success 200 {string} string "Item deletado com sucesso"
// @Failure 400 {string} string "ID não fornecido ou inválido"
// @Failure 500 {string} string "Erro ao deletar o item"
// @Router /itens/delete [delete]
func deleteItenHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := config.DB.Delete(&model.Iten{}, id).Error; err != nil {
		http.Error(w, "Erro ao deletar o item", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Item deletado com sucesso"))
}

// ==================== HANDLERS PARA CATEGORIAS ====================

// Listar todas as categorias
// @Summary Listar todas as categorias
// @Description Retorna todas as categorias cadastradas no banco de dados
// @Tags categorias
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Cat
// @Router /categorias [get]
func listCategoriasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var cats []model.Cat
	if err := config.DB.Find(&cats).Error; err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cats)
}

// Buscar uma única categoria pelo id (via query string: ?id=1)
// @Summary Buscar categoria por ID
// @Description Retorna uma única categoria baseada no ID fornecido via query string
// @Tags categorias
// @Accept  json
// @Produce  json
// @Param id query int true "ID da Categoria"
// @Success 200 {object} model.Cat
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
	var cat model.Cat
	if err := config.DB.First(&cat, id).Error; err != nil {
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(cat)
}

// Criar uma nova categoria (envie JSON via POST)
// @Summary Criar uma nova categoria
// @Description Cria uma categoria baseada no JSON fornecido
// @Tags categorias
// @Accept  json
// @Produce  json
// @Param cat body model.Cat true "Dados da Categoria"
// @Success 200 {object} model.Cat
// @Failure 400 {string} string "Erro ao decodificar a categoria"
// @Failure 500 {string} string "Erro ao criar a categoria"
// @Router /categorias/create [post]
func createCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var cat model.Cat
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}
	if err := config.DB.Create(&cat).Error; err != nil {
		http.Error(w, "Erro ao criar a categoria", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cat)
}

// Atualizar uma categoria (envie JSON via PUT, com o campo id preenchido)
// @Summary Atualizar uma categoria
// @Description Atualiza os dados de uma categoria existente baseada no ID no JSON
// @Tags categorias
// @Accept  json
// @Produce  json
// @Param cat body model.Cat true "Dados da Categoria (deve conter ID)"
// @Success 200 {object} model.Cat
// @Failure 400 {string} string "Erro ao decodificar a categoria"
// @Failure 500 {string} string "Erro ao atualizar a categoria"
// @Router /categorias/update [put]
func updateCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var cat model.Cat
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}
	if err := config.DB.Save(&cat).Error; err != nil {
		http.Error(w, "Erro ao atualizar a categoria", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cat)
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
	if err := config.DB.Delete(&model.Cat{}, id).Error; err != nil {
		http.Error(w, "Erro ao deletar a categoria", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Categoria deletada com sucesso"))
}
