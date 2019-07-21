package storage

import (
	"bitbucket.org/toggly/toggly-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	List() []*models.Environment
	Get(code string) (*models.Environment, error)
	Create(data *models.Environment) (*models.Environment, error)
	Update(data *models.Environment) (*models.Environment, error)
	Delete(code string)

	IsExist(code string) bool
}

// EnvironmentKey interface
type EnvironmentKey interface {
	Provision(data *models.EnvAPIKey) (*models.EnvAPIKey, error)
	Check(key string, secret string) error
}

// Package interface
type Package interface {
	List() []*models.Package
	Get(code string) (*models.Package, error)
	Create(data *models.Package) (*models.Package, error)
	Update(data *models.Package) (*models.Package, error)
	Delete(code string)

	Override(id primitive.ObjectID, paramID primitive.ObjectID, value interface{}) error

	IsExist(code string) bool
	IsParamExist(id primitive.ObjectID, paramID primitive.ObjectID) bool

	ReadParam(id primitive.ObjectID, paramID primitive.ObjectID) (*models.PackageParamLink, error)
}
