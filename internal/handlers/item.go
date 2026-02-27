package handlers

import (
	"encoding/json"
	"myapi/internal/config"
	"myapi/internal/models"
	"myapi/internal/services"
	"net/http"
	"strconv"
)

// ==================== HANDLERS PARA ITENS ====================

// Listar todos os itens
// @Summary Listar todos os itens
// @Description Retorna todos os itens cadastrados no banco de dados
// @Tags itens
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Iten
// @Router /itens [get]
func ListItensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var itens []models.Iten
	if err := config.DB.Find(&itens).Error; err != nil {
		http.Error(w, "Erro ao buscar itens", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(itens); err != nil {
		http.Error(w, "Erro ao codificar os itens", http.StatusInternalServerError)
	}
}

// Buscar um único item pelo id (via query string: ?id=1)
// @Summary Buscar item por ID
// @Description Retorna um único item baseado no ID fornecido via query string
// @Tags itens
// @Accept  json
// @Produce  json
// @Param id query int true "ID do Item"
// @Success 200 {object} models.Iten
// @Failure 400 {string} string "ID não fornecido ou inválido"
// @Failure 404 {string} string "Item não encontrado"
// @Router /itens/get [get]
func GetItenHandler(w http.ResponseWriter, r *http.Request) {
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
	var item models.Iten
	if err := config.DB.First(&item, id).Error; err != nil {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Erro ao codificar o item", http.StatusInternalServerError)
	}
}

// Buscar um item pelo campo "codigo"
// @Summary Buscar item por código
// @Description Retorna um único item baseado no código fornecido via query string
// @Tags itens
// @Accept  json
// @Produce  json
// @Param codigo query string true "Código do Item"
// @Success 200 {object} models.Iten
// @Failure 400 {string} string "Código não fornecido"
// @Failure 404 {string} string "Item não encontrado"
// @Router /itens/get-code [get]
func GetItenByCodigoHandler(w http.ResponseWriter, r *http.Request) {
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
	var item models.Iten
	// Busca o item onde o campo "codigo" é igual ao valor fornecido
	if err := config.DB.Where("codigo = ?", cod).First(&item).Error; err != nil {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Erro ao codificar o item", http.StatusInternalServerError)
	}
}

// Criar um novo item (envie JSON via POST)
// @Summary Criar um novo item
// @Description Cria um item baseado no JSON fornecido
// @Tags itens
// @Accept  json
// @Produce  json
// @Param item body models.Iten true "Dados do Item"
// @Success 200 {object} models.Iten
// @Failure 400 {string} string "Erro ao decodificar o item"
// @Failure 500 {string} string "Erro ao criar o item"
// @Router /itens/create [post]
func CreateItenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var item models.Iten
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Erro ao decodificar o item", http.StatusBadRequest)
		return
	}

	createdItem, err := services.CreateItem(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(createdItem); err != nil {
		http.Error(w, "Erro ao codificar o item criado", http.StatusInternalServerError)
	}
}

// Atualizar um item (envie JSON via PUT, com o campo id preenchido)
// @Summary Atualizar um item
// @Description Atualiza os dados de um item existente baseado no ID no JSON
// @Tags itens
// @Accept  json
// @Produce  json
// @Param item body models.Iten true "Dados do Item (deve conter ID)"
// @Success 200 {object} models.Iten
// @Failure 400 {string} string "Erro ao decodificar o item"
// @Failure 500 {string} string "Erro ao atualizar o item"
// @Router /itens/update [put]
func UpdateItenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var item models.Iten
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Erro ao decodificar o item", http.StatusBadRequest)
		return
	}
	if err := config.DB.Save(&item).Error; err != nil {
		http.Error(w, "Erro ao atualizar o item", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Erro ao codificar o item", http.StatusInternalServerError)
	}
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
func DeleteItenHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := config.DB.Delete(&models.Iten{}, id).Error; err != nil {
		http.Error(w, "Erro ao deletar o item", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write([]byte("Item deletado com sucesso")); err != nil {
		http.Error(w, "Erro ao escrever resposta", http.StatusInternalServerError)
	}
}
