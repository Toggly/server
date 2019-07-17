package service

import (
	"context"
	"strings"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/storage"
	"gopkg.in/toggly/go-utils.v2"
)

// Package Service
type Package struct {
	Storage *storage.MongoStorage
	Ctx     context.Context
	Config  *models.Config
	Logger  *utils.StructuredLogger
	Project *models.Project
}

// IsExist checks that param exists by code
func (a *Package) IsExist(code string) bool {
	return a.Storage.PackageCRUD(a.Project.ID).IsExist(code)
}

// Get param by code
func (a *Package) Get(code string) (*models.Package, error) {
	return a.Storage.PackageCRUD(a.Project.ID).Get(code)
}

// List envs by code
func (a *Package) List() []*models.Package {
	a.Logger.Debugf("Package.List for project [%s]", a.Project.Code)
	return a.Storage.PackageCRUD(a.Project.ID).List()
}

// Create param
func (a *Package) Create(data models.Package) (*models.Package, error) {
	if data.Code == "" {
		return nil, models.ErrBadRequest("Code is invalid")
	}

	a.Logger.Debugf("Package.Create: %+v", data)

	resp, err := a.Storage.PackageCRUD(a.Project.ID).Create(&data)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return nil, models.ErrConflict("Code is already exist")
		}
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}

// Update param
func (a *Package) Update(data models.Package) (*models.Package, error) {

	a.Logger.Debugf("Package.Update: %+v", data)

	item, _ := a.Storage.PackageCRUD(a.Project.ID).Get(data.Code)

	// revalue existing data

	item.Description = data.Description

	resp, err := a.Storage.PackageCRUD(a.Project.ID).Update(item)
	if err != nil {
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}
