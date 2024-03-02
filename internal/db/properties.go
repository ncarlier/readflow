package db

import "github.com/ncarlier/readflow/internal/model"

// PropertiesRepository is the repository interface to manage Propertiess
type PropertiesRepository interface {
	CreateProperties(properties model.Properties) (*model.Properties, error)
	GetProperties() (*model.Properties, error)
}
