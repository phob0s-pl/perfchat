package tests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	api "github.com/phob0s-pl/perfchat/apiv1"
)

// Server is a test server
type Server struct {
	Srv *http.Server
	API *api.API
}

// NewServer returns new test server with all api calls
func NewServer(address string) *Server {
	router := mux.NewRouter()

	httpSrv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
		Addr:         address,
	}

	serverAPI := api.NewAPI()
	AddAPI(router, serverAPI.AddUserRoute())
	AddAPI(router, serverAPI.PingRoute())
	AddAPI(router, serverAPI.GetUsersRoute())
	AddAPI(router, serverAPI.GetRoomsRoute())

	return &Server{
		Srv: httpSrv,
		API: serverAPI,
	}
}

func AddAPI(router *mux.Router, route *api.Route) {
	if route.WithPrefix {
		router.
			Methods(route.Method).
			PathPrefix(route.Pattern).
			Name(route.Name).
			Handler(Logger(route.HandlerFunc, route.Name))
		return
	}

	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(Logger(route.HandlerFunc, route.Name))
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		fmt.Printf(
			"%s %s %s %s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
