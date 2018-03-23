package main

import (
	"github.com/gorilla/mux"
)

// NewRouter returns new HTTP handler
func NewRouter() *mux.Router {
	return mux.NewRouter()
}
