package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/toggly/toggly-server/models"
	"bitbucket.org/toggly/toggly-server/service"
	"bitbucket.org/toggly/toggly-server/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	dbStore "github.com/nodely/go-mongo-store"
	"github.com/op/go-logging"
	"gopkg.in/toggly/go-utils.v2"
)

// Toggly struct
type Toggly struct {
	Dbs    *dbStore.DbStorage
	Ctx    context.Context
	Config *models.Config
	Logger *logging.Logger
}

// Run Toggly App
func (t *Toggly) Run() {
	log := t.Logger
	routes := t.Router("/")
	if t.Config.Port == 0 {
		t.Config.Port = 8080
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", t.Config.Port),
		Handler: chi.ServerBaseContext(t.Ctx, routes),
	}
	go func() {
		<-t.Ctx.Done()
		if err := srv.Shutdown(t.Ctx); err != nil {
			log.Error("REST stop error")
		}
		log.Info("REST server stopped")
	}()
	log.Infof("HTTP server listening on %s", srv.Addr)
	err := srv.ListenAndServe()
	log.Infof("HTTP server terminated, %s", err)
}

// Router returns router configuration
func (t *Toggly) Router(basePath string) chi.Router {
	router := chi.NewRouter()
	router.Use(utils.RequestIDCtx)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.RequestLogger(&utils.StructuredLogger{Logger: t.Logger, R: nil}))
	if t.Config.MultiUserMode {
		router.Use(OwnerCtx(""))
	} else {
		t.Logger.Info("Single user mode is enabled")
		router.Use(OwnerCtx("NO_OWNER_ID_MODE"))
	}
	router.Route(basePath, t.versions)
	return router
}

// versions for routing
func (t *Toggly) versions(router chi.Router) {
	router.Use(utils.VersionCtx("v1"))
	router.Route("/v1", t.v1)
}

// routes for API v1
func (t *Toggly) v1(router chi.Router) {
	services := make(map[string]interface{}, 0)
	services["project"] = &service.Project{
		Storage: &storage.MongoStorage{
			Dbs: t.Dbs,
		},
		Ctx:    t.Ctx,
		Config: t.Config,
	}
	services["params"] = &service.Param{
		Storage: &storage.MongoStorage{
			Dbs: t.Dbs,
		},
		Ctx:    t.Ctx,
		Config: t.Config,
	}
	services["envs"] = &service.Environment{
		Storage: &storage.MongoStorage{
			Dbs: t.Dbs,
		},
		Ctx:    t.Ctx,
		Config: t.Config,
	}

	router.Mount("/project", (&ProjectEndpoints{
		Dbs:      t.Dbs,
		Ctx:      t.Ctx,
		Config:   t.Config,
		Logger:   t.Logger,
		Services: services,
	}).Routes())

	router.Mount("/project/{ProjectCode}/param", (&ParamEndpoints{
		Dbs:      t.Dbs,
		Ctx:      t.Ctx,
		Config:   t.Config,
		Logger:   t.Logger,
		Services: services,
	}).Routes())

	router.Mount("/project/{ProjectCode}/env", (&EnvEndpoints{
		Dbs:      t.Dbs,
		Ctx:      t.Ctx,
		Config:   t.Config,
		Logger:   t.Logger,
		Services: services,
	}).Routes())
}
