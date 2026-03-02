package services

import (
	"errors"
	"myapi/internal/config"
	"myapi/internal/models"
)

func CreateItem(item *models.Iten) (*models.Iten, error) {
	if item.Nome == "" {
		return nil, errors.New("nome do item não pode ser vazio")
	}

	if err := config.DB.Create(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}
