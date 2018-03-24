package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	api "github.com/phob0s-pl/perfchat/apiv1"
	"github.com/phob0s-pl/perfchat/chat"
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
	log.Debugf("Read configuration from %q: %+v", *configPath, config)

	serverAPI.Engine.SetRoomLimit(config.RoomLimit)
	serverAPI.Engine.SetUsersLimit(config.UsersLimit)

	// Register all API calls
	AddAPI(router, serverAPI.AddUserRoute())
	AddAPI(router, serverAPI.PingRoute())
	AddAPI(router, serverAPI.GetUsersRoute())
	AddAPI(router, serverAPI.GetRoomsRoute())
	AddAPI(router, serverAPI.CreateRoomRoute())
	AddAPI(router, serverAPI.JoinRoomRoute())
	AddAPI(router, serverAPI.ExitRoomRoute())
	AddAPI(router, serverAPI.SendMessageRoute())
	AddAPI(router, serverAPI.ReceiveMessageRoute())
	AddAPI(router, serverAPI.StatsRoute())

	admin := &chat.User{
		AuthID: config.AuthID,
		Name:   "admin",
		Role:   chat.AdminRole,
		Token:  config.Token,
	}
	if err := serverAPI.Engine.AddUser(admin); err != nil {
		log.Fatalf("Failed to add admin, err=%s", err)
	}

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
	if route.WithPrefix {
		router.
			Methods(route.Method).
			PathPrefix(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
		return
	}

	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(route.HandlerFunc)
}
