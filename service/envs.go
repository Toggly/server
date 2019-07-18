package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/storage"
	"github.com/gofrs/uuid"
	"gopkg.in/toggly/go-utils.v2"
)

// Environment Service
type Environment struct {
	Storage *storage.MongoStorage
	Ctx     context.Context
	Config  *models.Config
	Logger  *utils.StructuredLogger
	Project *models.Project
}

// IsExist checks that param exists by code
func (a *Environment) IsExist(code string) bool {
	return a.Storage.EnvCRUD(a.Project.ID).IsExist(code)
}

// Get param by code
func (a *Environment) Get(code string) (*models.Environment, error) {
	return a.Storage.EnvCRUD(a.Project.ID).Get(code)
}

// List envs by code
func (a *Environment) List() []*models.Environment {
	a.Logger.Debugf("Env.List for project [%s]", a.Project.Code)
	return a.Storage.EnvCRUD(a.Project.ID).List()
}

// Create param
func (a *Environment) Create(data models.Environment) (*models.Environment, error) {
	if data.Code == "" {
		return nil, models.ErrBadRequest("Code is invalid")
	}

	a.Logger.Debugf("Env.Create: %+v", data)

	resp, err := a.Storage.EnvCRUD(a.Project.ID).Create(&data)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return nil, models.ErrConflict("Code is already exist")
		}
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}

// Update param
func (a *Environment) Update(data models.Environment) (*models.Environment, error) {

	a.Logger.Debugf("Env.Update: %+v", data)

	item, _ := a.Storage.EnvCRUD(a.Project.ID).Get(data.Code)

	// revalue existing data
	item.Protected = data.Protected
	item.Description = data.Description

	resp, err := a.Storage.EnvCRUD(a.Project.ID).Update(item)
	if err != nil {
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}

// KeyProvision param
func (a *Environment) KeyProvision(data models.EnvAPIKey) (*models.EnvAPIKey, error) {

	if data.Name == "" {
		return nil, models.ErrBadRequest("Name is empty")
	}

	u, _ := uuid.NewV4()
	data.Key = u.String()
	secret := make([]byte, 64)
	rand.Read(secret)
	data.Secret = fmt.Sprintf("%x", secret)
	data.CreatedDate = time.Now()

	resp, err := a.Storage.EnvKeyCRUD(a.Project.ID).Provision(&data)
	if err != nil {
		return nil, models.ErrInternalServer(err.Error())
	}

	return resp, nil
}
