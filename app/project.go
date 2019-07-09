package app

import (
	"context"
	"net/http"

	"bitbucket.org/toggly/toggly-server/models"
	"github.com/go-chi/chi"
	dbStore "github.com/nodely/go-mongo-store"
	"github.com/op/go-logging"
)

// ProjectEndpoints API struct
type ProjectEndpoints struct {
	Dbs    *dbStore.DbStorage
	Ctx    context.Context
	Config *models.Config
	Logger *logging.Logger
}

// Routes returns api endpoints
func (a *ProjectEndpoints) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", a.list)
	})
	return router
}

func (a *ProjectEndpoints) list(w http.ResponseWriter, r *http.Request) {
}
