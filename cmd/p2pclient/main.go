package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

var adapterType = flag.String("adapter", "sim", `node adapter to use (one of "sim", "exec")`)

func main() {
	flag.Parse()

	// set the log level to info
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	// register a single chat-service
	services := map[string]adapters.ServiceFunc{
		"chat-service": func(ctx *adapters.ServiceContext) (node.Service, error) {
			return newChatService(ctx.Config.ID), nil
		},
	}
	adapters.RegisterServices(services)

	// create the NodeAdapter
	adapter := adapters.NewSimAdapter(services)

	// start the HTTP API
	log.Info("starting simulation server on 0.0.0.0:8888...")
	network := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "chat-service",
	})
	if err := http.ListenAndServe(":8888", simulations.NewServer(network)); err != nil {
		log.Crit("error starting simulation server", "err", err)
	}
}
