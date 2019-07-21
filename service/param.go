package service

import (
	"context"
	"fmt"
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
func (a *Param) Override(code string, data bson.M) error {
	if data["type"] == "environment" {

	} else if data["type"] == "package" {
		pkgCode := data["packageCode"].(string)
		// validate that this param is not in that package
		pkg, err := a.Storage.PackageCRUD(a.Project.ID).Get(pkgCode)
		if err != nil {
			a.Logger.Errorf("Param.Override: %s", err.Error())
			return models.ErrNotFound(err.Error())
		}
		param, err := a.Storage.ParamCRUD(a.Project.ID).Get(code)
		if err != nil {
			a.Logger.Errorf("Param.Override: %s", err.Error())
			return models.ErrNotFound(err.Error())
		}
		if ok := a.Storage.PackageCRUD(a.Project.ID).IsParamExist(pkg.ID, param.ID); ok {
			return models.ErrBadRequest(fmt.Sprintf("Param [%s] is already overriden", param.Code))
		}

		return a.Storage.PackageCRUD(a.Project.ID).Override(pkg.ID, param.ID, data["value"])

	}
	return models.ErrBadRequest("Type is not valid")
}

// GetParamValue gets param value from env/package
func (a *Param) GetParamValue(pkgCode string, paramCode string) (*models.ParameterValue, error) {
	// gets pkg
	pkg, err := a.Storage.PackageCRUD(a.Project.ID).Get(pkgCode)
	if err != nil {
		return nil, models.ErrInternalServer(err.Error())
	}
	// gets param
	param, err := a.Storage.ParamCRUD(a.Project.ID).Get(paramCode)
	if err != nil {
		return nil, models.ErrInternalServer(err.Error())
	}
	// reads param override
	link, err := a.Storage.PackageCRUD(a.Project.ID).ReadParam(pkg.ID, param.ID)
	if link == nil || err != nil {
		return nil, models.ErrNotFound(fmt.Sprintf("Param [%s] is not found", paramCode))
	}
	// finally returns data
	return &models.ParameterValue{
		Code:  param.Code,
		Type:  param.Type,
		Value: link.Value,
	}, nil
}
