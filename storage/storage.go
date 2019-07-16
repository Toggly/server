package storage

import (
	"bitbucket.org/toggly/toggly-server/models"
)

// Project interface
type Project interface {
	List() []*models.Project
	Get(code string) (*models.Project, error)
	Create(data *models.Project) (*models.Project, error)
	Update(data *models.Project) (*models.Project, error)
	Delete(code string)

	IsExist(code string) bool
}

// Parameter interface
type Parameter interface {
	List() []*models.Parameter
	Get(code string) (*models.Parameter, error)
	Create(data *models.Parameter) (*models.Parameter, error)
	Update(data *models.Parameter) (*models.Parameter, error)
	Delete(code string)

	IsExist(code string) bool
}

// Environment interface
type Environment interface {
}
