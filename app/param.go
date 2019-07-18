package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/service"
	"github.com/go-chi/chi"
	dbStore "github.com/nodely/go-mongo-store"
	"github.com/op/go-logging"
	"gopkg.in/mgo.v2/bson"
)

type contextKey int

const (
	//ContextProjectKey key
	ContextProjectKey contextKey = iota
)

// ParamEndpoints API struct
type ParamEndpoints struct {
	Dbs      *dbStore.DbStorage
	Ctx      context.Context
	Config   *models.Config
	Logger   *logging.Logger
	Services map[string]interface{}
}

// Routes returns api endpoints
func (a *ParamEndpoints) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(WithProjectCtx(a.Services["project"].(*service.Project)))
		r.Get("/", a.list)
		r.Post("/", a.create)
		r.Put("/{ParamCode}", a.update)
		r.Delete("/{ParamCode}", a.delete)

		r.Post("/{ParamCode}/assign", a.override)
	})
	return router
}

// GetterRoutes returns routes for getting api
func (a *ParamEndpoints) GetterRoutes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(WithProjectCtx(a.Services["project"].(*service.Project)))
		r.Use(EnvironmentCtx(a.Services["envs"].(*service.Environment)))

		r.Get("/{ParamCode}", a.getParamValue)
	})
	return router
}

func (a *ParamEndpoints) withParamService(r *http.Request) *service.Param {
	log := GetLogger(r)
	srv := a.Services["params"].(*service.Param)
	srv.Project = r.Context().Value(ContextProjectKey).(*models.Project)
	srv.Logger = log
	return srv
}

func (a *ParamEndpoints) withProjectService(r *http.Request) *service.Project {
	log := GetLogger(r)
	srv := a.Services["project"].(*service.Project)
	srv.Logger = log
	return srv
}

func (a *ParamEndpoints) list(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	recs := a.withParamService(r).List()
	log.Debugf("Param.list: %d items found", len(recs))
	models.JSONResponse(w, r, recs)
}

func (a *ParamEndpoints) create(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Parameter
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// fill up default values
	data.ProjectID = a.withParamService(r).Project.ID

	// create param
	resp, err := a.withParamService(r).Create(data)
	if err != nil {
		log.Errorf("Param.Service.Create: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Param: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *ParamEndpoints) override(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "ParamCode")

	// verify param existance
	if ok := a.withParamService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Param with code [%s] is not found", code))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data bson.M
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	if err := a.withParamService(r).Override(data); err != nil {
		models.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *ParamEndpoints) update(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "ParamCode")

	// verify param existance
	if ok := a.withParamService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Param with code [%s] is not found", code))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Parameter
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// fill up default values
	data.Code = code
	data.ProjectID = a.withParamService(r).Project.ID

	// update params desc and values
	resp, err := a.withParamService(r).Update(data)
	if err != nil {
		log.Errorf("Param.Service.Update: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Param: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *ParamEndpoints) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func (a *ParamEndpoints) getParamValue(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
