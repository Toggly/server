package service

import (
	"context"
	"strings"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/storage"
	"github.com/op/go-logging"
)

// Project Service
type Project struct {
	Storage *storage.MongoStorage
	Ctx     context.Context
	Config  *models.Config
	Logger  *logging.Logger
}

// IsExist checks that project exists by code
func (a *Project) IsExist(code string) bool {
	return a.Storage.ProjectCRUD().IsExist(code)
}

// Get project by code
func (a *Project) Get(code string) *models.Project {
	return a.Storage.ProjectCRUD().Get(code)
}

// List project by code
func (a *Project) List() []*models.Project {
	return a.Storage.ProjectCRUD().List()
}

// Create project
func (a *Project) Create(data models.Project) (*models.Project, error) {
	if data.Code == "" {
		return nil, models.ErrBadRequest("Code is invalid")
	}
	if data.Name == "" {
		return nil, models.ErrBadRequest("Name is invalid")
	}

	a.Logger.Debugf("Project.Create: %+v", data)

	resp, err := a.Storage.ProjectCRUD().Create(&data)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return nil, models.ErrConflict("Code is already exist")
		}
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}

// Update project
func (a *Project) Update(data models.Project) (*models.Project, error) {
	if data.Name == "" {
		return nil, models.ErrBadRequest("Name is invalid")
	}

	a.Logger.Debugf("Project.Update: %+v", data)

	item := a.Storage.ProjectCRUD().Get(data.Code)

	// revalue existing data
	item.Name = data.Name
	item.Description = data.Description

	resp, err := a.Storage.ProjectCRUD().Update(item)
	if err != nil {
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}
