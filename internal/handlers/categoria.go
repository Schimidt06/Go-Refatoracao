package handlers

import (
	"encoding/json"
	"fmt"
	"myapi/internal/config"
	"myapi/internal/models"
	"net/http"
	"strconv"
)

func ScalarHandler(w http.ResponseWriter, r *http.Request) {
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

func ListCategoriasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var categorias []models.Categoria
	if err := config.DB.Find(&categorias).Error; err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(categorias); err != nil {
		http.Error(w, "Erro ao codificar categorias", http.StatusInternalServerError)
	}
}

func GetCategoriaHandler(w http.ResponseWriter, r *http.Request) {
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
	var categorias models.Categoria
	if err := config.DB.First(&categorias, id).Error; err != nil {
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(categorias); err != nil {
		http.Error(w, "Erro ao codificar categoria", http.StatusInternalServerError)
	}
}

func CreateCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var categorias models.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categorias); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}
	if err := config.DB.Create(&categorias).Error; err != nil {
		http.Error(w, "Erro ao criar a categoria", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(categorias); err != nil {
		http.Error(w, "Erro ao codificar categoria criada", http.StatusInternalServerError)
	}
}

func UpdateCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var categorias models.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categorias); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}
	if err := config.DB.Save(&categorias).Error; err != nil {
		http.Error(w, "Erro ao atualizar a categoria", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(categorias); err != nil {
		http.Error(w, "Erro ao codificar categoria atualizada", http.StatusInternalServerError)
	}
}

func DeleteCategoriaHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := config.DB.Delete(&models.Categoria{}, id).Error; err != nil {
		http.Error(w, "Erro ao deletar a categoria", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write([]byte("Categoria deletada com sucesso")); err != nil {
		http.Error(w, "Erro ao escrever resposta", http.StatusInternalServerError)
	}
}
