package app

import (
	"context"
	"encoding/json"
	"errors"
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
		r.Get("/{ParamCode}", a.get)
		r.Put("/{ParamCode}", a.update)
		r.Delete("/{ParamCode}", a.delete)

		r.Post("/{ParamCode}/assign", a.override)

		r.Get("/validateCode", a.validateCode)
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

func (a *ParamEndpoints) withPackageService(r *http.Request) *service.Package {
	log := GetLogger(r)
	srv := a.Services["pkgs"].(*service.Package)
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
	q := r.URL.Query().Get("q")
	recs := a.withParamService(r).List(q)
	log.Debugf("Param.list: %d items found", len(recs))
	models.JSONResponse(w, r, recs)
}

func (a *ParamEndpoints) get(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "ParamCode")

	res, err := a.withParamService(r).Get(code)
	if err != nil {
		log.Error(err.Error())
		models.ErrorResponseWithStatus(w, r, err, http.StatusNotFound)
		return
	}

	models.JSONResponse(w, r, res)
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

	// verify param existence
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

	if err := a.withParamService(r).Override(code, data); err != nil {
		models.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *ParamEndpoints) update(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "ParamCode")

	// verify param existence
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

func (a *ParamEndpoints) validateCode(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Errorf("Code [%s] is not defined", code)
		models.ErrorResponseWithStatus(w, r, errors.New("Code is not defined"), http.StatusBadRequest)
		return
	}

	// update params desc and values
	if ok := a.withParamService(r).IsExist(code); ok {
		log.Warningf("Code [%s] is exist", code)
		models.ErrorResponseWithStatus(w, r, errors.New("Code is exist"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *ParamEndpoints) delete(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "ParamCode")
	if code == "" {
		log.Errorf("Code [%s] is not defined", code)
		models.ErrorResponseWithStatus(w, r, errors.New("Code is not defined"), http.StatusBadRequest)
		return
	}

	// check parameter existance
	if ok := a.withParamService(r).IsExist(code); !ok {
		log.Warningf("Code [%s] is not exist", code)
		models.ErrorResponseWithStatus(w, r, errors.New("Code is not exist"), http.StatusBadRequest)
		return
	}

	// todo add check for parameter inheritance

	// delete parameter
	if err := a.withParamService(r).Delete(code); err != nil {
		log.Error(err.Error())
		models.ErrorResponseWithStatus(w, r, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *ParamEndpoints) getParamValue(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)

	pkgCode := chi.URLParam(r, "PackageCode")
	paramCode := chi.URLParam(r, "ParamCode")

	// verify package existence
	if ok := a.withPackageService(r).IsExist(pkgCode); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Package with code [%s] is not found", pkgCode))
		return
	}

	// verify param existence
	if ok := a.withParamService(r).IsExist(paramCode); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Param with code [%s] is not found", paramCode))
		return
	}

	log.Debugf("Getting value for param [%s] from package [%s]", paramCode, pkgCode)

	// trying to get value from package
	resp, err := a.withParamService(r).GetParamValue(pkgCode, paramCode)
	if err != nil {
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Param: [%+v]", resp)

	models.JSONResponse(w, r, resp)
}
