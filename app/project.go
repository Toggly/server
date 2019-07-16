package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/service"
	"github.com/go-chi/chi"
	dbStore "github.com/nodely/go-mongo-store"
	"github.com/op/go-logging"
)

// ProjectEndpoints API struct
type ProjectEndpoints struct {
	Dbs      *dbStore.DbStorage
	Ctx      context.Context
	Config   *models.Config
	Logger   *logging.Logger
	Services map[string]interface{}
}

// Routes returns api endpoints
func (a *ProjectEndpoints) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", a.list)
		group.Post("/", a.create)
		group.Put("/{ProjectCode}", a.update)
		group.Get("/{ProjectCode}", a.get)
	})
	return router
}

func (a *ProjectEndpoints) withProjectService(r *http.Request) *service.Project {
	log := GetLogger(r)
	srv := a.Services["project"].(*service.Project)
	srv.Logger = log
	return srv
}

func (a *ProjectEndpoints) list(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	recs := a.withProjectService(r).List()
	log.Debugf("Project.list: %d items found", len(recs))
	models.JSONResponse(w, r, recs)
}

func (a *ProjectEndpoints) create(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Project
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// fill up default values
	data.OwnerID = models.OwnerFromContext(r)
	data.RegDate = time.Now()
	data.Status = 1

	// create project
	resp, err := a.withProjectService(r).Create(data)
	if err != nil {
		log.Errorf("Project.Service.Create: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Project: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *ProjectEndpoints) update(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "ProjectCode")

	// verify project existance
	if ok := a.withProjectService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Project with code [%s] is not found", code))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Project
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// update project
	resp, err := a.withProjectService(r).Update(data)
	if err != nil {
		log.Errorf("Project.Service.Update: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Project: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *ProjectEndpoints) get(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "ProjectCode")

	// verify project existance
	if ok := a.withProjectService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Project with code [%s] is not found", code))
		return
	}

	resp, err := a.withProjectService(r).Get(code)
	if err != nil {
		a.Logger.Error("Not found")
		models.NotFoundResponse(w, r, "Resource not found")
		return
	}

	log.Debugf("Project: %+v", resp)

	models.JSONResponse(w, r, resp)
}
