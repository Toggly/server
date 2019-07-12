package storage

import "bitbucket.org/toggly/toggly-server/models"

// Project interface
type Project interface {
	List() []*models.Project
	Get(code string) *models.Project
	Create(data *models.Project) (*models.Project, error)
	Update(data *models.Project) (*models.Project, error)
	Delete(code string)
	IsExist(code string) bool
}

// Environment interface
type Environment interface {
}
