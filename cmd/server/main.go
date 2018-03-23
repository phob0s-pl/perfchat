package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	api "github.com/phob0s-pl/perfchat/apiv1"
	log "github.com/sirupsen/logrus"
)

const (
	// ConfigPath is default configuration path
	ConfigPath = "server.conf"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	var (
		serverAPI = api.NewAPI()
		router    = NewRouter()
	)

	configPath := flag.String("conf", ConfigPath, "path to config file")
	flag.Parse()

	config, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config, err=%s", err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("Read configuration from %q: %+v", configPath, config)

	// Register all API calls
	AddAPI(router, serverAPI.AddUserRoute())

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
		Addr:         config.Address,
	}

	log.Infof("Serving at %s", config.Address)

	if err := srv.ListenAndServe(); err != nil {
		log.Errorf("Failed to serve: %s", err)
	}
}

func AddAPI(router *mux.Router, route *api.Route) {
	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(route.HandlerFunc)
	log.Debugf("Registered API for %q", route.Pattern)
}
