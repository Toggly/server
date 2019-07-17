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
)

// PackageEndpoints API struct
type PackageEndpoints struct {
	Dbs      *dbStore.DbStorage
	Ctx      context.Context
	Config   *models.Config
	Logger   *logging.Logger
	Services map[string]interface{}
}

// Routes returns api endpoints
func (a *PackageEndpoints) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Use(WithProjectCtx(a.Services["project"].(*service.Project)))
		group.Get("/", a.list)
		group.Post("/", a.create)
		group.Get("/{PackageCode}", a.get)
		group.Put("/{PackageCode}", a.update)
		group.Delete("/{PackageCode}", a.delete)
	})
	return router
}

func (a *PackageEndpoints) withProjectService(r *http.Request) *service.Project {
	log := GetLogger(r)
	srv := a.Services["project"].(*service.Project)
	srv.Logger = log
	return srv
}

func (a *PackageEndpoints) withPackageService(r *http.Request) *service.Package {
	srv := a.Services["pkgs"].(*service.Package)
	srv.Logger = GetLogger(r)
	srv.Project = r.Context().Value(ContextProjectKey).(*models.Project)
	return srv
}

func (a *PackageEndpoints) list(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	recs := a.withPackageService(r).List()
	log.Debugf("Package.list: %d items found", len(recs))
	models.JSONResponse(w, r, recs)
}

func (a *PackageEndpoints) create(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Package
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// fill up default values
	data.ProjectID = a.withPackageService(r).Project.ID

	// create param
	resp, err := a.withPackageService(r).Create(data)
	if err != nil {
		log.Errorf("Package.Service.Create: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Pkg: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *PackageEndpoints) update(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "PackageCode")

	// verify param existence
	if ok := a.withPackageService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Package with code [%s] is not found", code))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Package
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// fill up default values
	data.Code = code
	data.ProjectID = a.withPackageService(r).Project.ID

	// update params desc and values
	resp, err := a.withPackageService(r).Update(data)
	if err != nil {
		log.Errorf("Package.Service.Update: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Pkg: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *PackageEndpoints) get(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "PackageCode")

	// verify project existance
	if ok := a.withPackageService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Package with code [%s] is not found", code))
		return
	}

	resp, err := a.withPackageService(r).Get(code)
	if err != nil {
		a.Logger.Error("Not found")
		models.NotFoundResponse(w, r, "Resource not found")
		return
	}

	log.Debugf("Pkg: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *PackageEndpoints) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
