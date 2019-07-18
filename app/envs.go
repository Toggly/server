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

// EnvEndpoints API struct
type EnvEndpoints struct {
	Dbs      *dbStore.DbStorage
	Ctx      context.Context
	Config   *models.Config
	Logger   *logging.Logger
	Services map[string]interface{}
}

// Routes returns api endpoints
func (a *EnvEndpoints) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Use(WithProjectCtx(a.Services["project"].(*service.Project)))
		group.Get("/", a.list)
		group.Post("/", a.create)
		group.Get("/{EnvCode}", a.get)
		group.Put("/{EnvCode}", a.update)
		group.Delete("/{EnvCode}", a.delete)

		// key stuff
		group.Post("/{EnvCode}/key", a.createKey)
	})
	return router
}

func (a *EnvEndpoints) withProjectService(r *http.Request) *service.Project {
	log := GetLogger(r)
	srv := a.Services["project"].(*service.Project)
	srv.Logger = log
	return srv
}

func (a *EnvEndpoints) withEnvService(r *http.Request) *service.Environment {
	srv := a.Services["envs"].(*service.Environment)
	srv.Logger = GetLogger(r)
	srv.Project = r.Context().Value(ContextProjectKey).(*models.Project)
	return srv
}

func (a *EnvEndpoints) list(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	recs := a.withEnvService(r).List()
	log.Debugf("Env.list: %d items found", len(recs))
	models.JSONResponse(w, r, recs)
}

func (a *EnvEndpoints) create(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Environment
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// fill up default values
	data.ProjectID = a.withEnvService(r).Project.ID

	// create param
	resp, err := a.withEnvService(r).Create(data)
	if err != nil {
		log.Errorf("Env.Service.Create: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Env: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *EnvEndpoints) update(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "EnvCode")

	// verify param existance
	if ok := a.withEnvService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Env with code [%s] is not found", code))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.Environment
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}

	// fill up default values
	data.Code = code
	data.ProjectID = a.withEnvService(r).Project.ID

	// update params desc and values
	resp, err := a.withEnvService(r).Update(data)
	if err != nil {
		log.Errorf("Env.Service.Update: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Env: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *EnvEndpoints) get(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "EnvCode")

	// verify project existance
	if ok := a.withEnvService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Project with code [%s] is not found", code))
		return
	}

	resp, err := a.withEnvService(r).Get(code)
	if err != nil {
		a.Logger.Error("Not found")
		models.NotFoundResponse(w, r, "Resource not found")
		return
	}

	log.Debugf("Env: %+v", resp)

	models.JSONResponse(w, r, resp)
}

func (a *EnvEndpoints) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func (a *EnvEndpoints) createKey(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	code := chi.URLParam(r, "EnvCode")

	// verify project existance
	if ok := a.withEnvService(r).IsExist(code); !ok {
		models.NotFoundResponse(w, r, fmt.Sprintf("Project with code [%s] is not found", code))
		return
	}

	env, _ := a.withEnvService(r).Get(code)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Can't read request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	var data models.EnvAPIKey
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Can't parse request body")
		models.ErrorResponseWithStatus(w, r, err, http.StatusInternalServerError)
		return
	}
	data.EnvID = env.ID

	resp, err := a.withEnvService(r).KeyProvision(data)
	if err != nil {
		log.Errorf("Env.Key.Provision: %s", err.Error())
		models.ErrorResponse(w, r, err)
		return
	}

	log.Debugf("Key: %+v", resp)

	models.JSONResponse(w, r, resp)
}
