package handlers

import (
	"encoding/json"
	"fmt"
	"myapi/internal/models"
	"myapi/internal/repositories"
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
	repository := repositories.NewCategoriaRepository()
	categorias, err := repository.ListAll()
	if err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categorias)
}

func GetCategoriaHandler(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewCategoriaRepository()
	categoria, err := repository.GetByID(id)
	if err != nil {
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(categoria)
}

func CreateCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	var categoria models.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categoria); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}

	repository := repositories.NewCategoriaRepository()
	createdCategoria, err := repository.Create(&categoria)
	if err != nil {
		http.Error(w, "Erro ao criar a categoria", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(createdCategoria)
}

func UpdateCategoriaHandler(w http.ResponseWriter, r *http.Request) {
	var categoria models.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categoria); err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}

	repository := repositories.NewCategoriaRepository()
	if err := repository.Update(&categoria); err != nil {
		http.Error(w, "Erro ao atualizar the categoria", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categoria)
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

	repository := repositories.NewCategoriaRepository()
	if err := repository.Delete(id); err != nil {
		http.Error(w, "Erro ao deletar a categoria", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Categoria deletada com sucesso"))
}
