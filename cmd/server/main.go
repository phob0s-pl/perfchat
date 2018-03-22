package main

import (
	"os"

	"flag"

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

	log.Println(config)

}
