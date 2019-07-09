package main

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/toggly/toggly-server/app"
	"bitbucket.org/toggly/toggly-server/models"
	dbStore "github.com/nodely/go-mongo-store"
	"github.com/op/go-logging"
	"gopkg.in/nodely/mongo-session.v3"
	"gopkg.in/session.v3"
	"gopkg.in/yaml.v2"
)

const logo = `

::::::::::: ::::::::   ::::::::   ::::::::  :::     :::   :::       ::::::::   ::::::::  :::::::::  :::::::::: 
    :+:    :+:    :+: :+:    :+: :+:    :+: :+:     :+:   :+:      :+:    :+: :+:    :+: :+:    :+: :+:        
    +:+    +:+    +:+ +:+        +:+        +:+      +:+ +:+       +:+        +:+    +:+ +:+    +:+ +:+        
    +#+    +#+    +:+ :#:        :#:        +#+       +#++:        +#+        +#+    +:+ +#++:++#:  +#++:++#   
    +#+    +#+    +#+ +#+   +#+# +#+   +#+# +#+        +#+         +#+        +#+    +#+ +#+    +#+ +#+        
    #+#    #+#    #+# #+#    #+# #+#    #+# #+#        #+#         #+#    #+# #+#    #+# #+#    #+# #+#        
    ###     ########   ########   ########  ########## ###          ########   ########  ###    ### ########## 

`

func init() {
	var format = logging.MustStringFormatter(
		`%{color} ▶ %{level:-8s}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	log := logging.MustGetLogger("logger")
	log.Info(logo)

	ctx = context.WithValue(ctx, models.ContextLoggerKey, log)

	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warning("interrupt signal \x1b[31m✘\x1b[0m")
		cancel()
	}()

	config := loadConfigs(os.Getenv("APP_CONFIG_PATH"))

	// connects to session storage
	mgoStore, err := mongo.NewMongoStore(&mongo.Options{
		Connection: config.Storage.Connection,
		DB:         config.Storage.Name,
		Collection: "sessions",
		Logger:     log,
	})
	if err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}

	session.InitManager(
		session.SetCookieName("TGLY_SID"),
		session.SetSign([]byte(config.Sessions["key"])),
		session.SetStore(mgoStore),
	)

	// connects to db storage
	dbs, err := dbStore.NewMongoStorage(config.Storage.Connection)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	// set db
	dbs.WithName(config.Storage.Name)

	log.Info("API server started")

	app := &app.Toggly{
		Dbs:    dbs,
		Ctx:    ctx,
		Config: config,
		Logger: log,
	}

	app.Run()

	log.Info("Bye! 🖐")
}

func loadConfigs(cfgPath string) *models.Config {
	confContent, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}
	// expand environment variables
	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &models.Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		panic(err)
	}
	return conf
}
