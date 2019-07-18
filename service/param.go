package service

import (
	"context"
	"strings"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/storage"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/toggly/go-utils.v2"
)

// Param Service
type Param struct {
	Storage *storage.MongoStorage
	Ctx     context.Context
	Config  *models.Config
	Logger  *utils.StructuredLogger
	Project *models.Project
}

// IsExist checks that param exists by code
func (a *Param) IsExist(code string) bool {
	return a.Storage.ParamCRUD(a.Project.ID).IsExist(code)
}

// Get param by code
func (a *Param) Get(code string) (*models.Parameter, error) {
	return a.Storage.ParamCRUD(a.Project.ID).Get(code)
}

// List project by code
func (a *Param) List() []*models.Parameter {
	a.Logger.Debugf("Param.List for project [%s]", a.Project.Code)
	return a.Storage.ParamCRUD(a.Project.ID).List()
}

// Create param
func (a *Param) Create(data models.Parameter) (*models.Parameter, error) {
	if data.Code == "" {
		return nil, models.ErrBadRequest("Code is empty")
	}
	if data.Type == "" {
		return nil, models.ErrBadRequest("Type is empty")
	}
	if data.Value == nil {
		return nil, models.ErrBadRequest("Value is empty")
	}

	if models.IsCodeValid(data.Code) {
		return nil, models.ErrBadRequest("Code is invalid")
	}

	if data.Type != models.ParameterTypeBool &&
		data.Type != models.ParameterTypeFloat &&
		data.Type != models.ParameterTypeInt &&
		data.Type != models.ParameterTypeString {
		return nil, models.ErrBadRequest("Type is invalid")
	}

	a.Logger.Debugf("Param.Create: %+v", data)

	resp, err := a.Storage.ParamCRUD(a.Project.ID).Create(&data)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return nil, models.ErrConflict("Code is already exist")
		}
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}

// Update param
func (a *Param) Update(data models.Parameter) (*models.Parameter, error) {
	a.Logger.Debugf("Param.Update: %+v", data)

	item, _ := a.Storage.ParamCRUD(a.Project.ID).Get(data.Code)

	// revalue existing data
	item.Value = data.Value
	item.AllowedValues = data.AllowedValues
	item.Description = data.Description

	resp, err := a.Storage.ParamCRUD(a.Project.ID).Update(item)
	if err != nil {
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}

// Override parameter by package/env
func (a *Param) Override(data bson.M) error {
	if data["type"] == "environment" {

	} else if data["type"] == "package" {

	}
	return models.ErrBadRequest("Type is not valid")
}
