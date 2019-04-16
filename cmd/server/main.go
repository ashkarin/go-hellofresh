package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/ashkarin/ashkarin-api-test/internal/config"
	"github.com/ashkarin/ashkarin-api-test/internal/server"
)

func main() {
	configPath := flag.String("config", "", "path to the configuration (JSON)")
	flag.Parse()

	var cfg *config.Config
	var err error
	if *configPath == "" {
		log.Infof("Load configuration from the environment variables")
		cfg, err = config.GetConfigFromEnv()
	} else {
		log.Infof("Load configuration from the file: %s", *configPath)
		cfg, err = config.GetConfig(*configPath)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Create and run the server
	srv := &server.Server{}
	srv.Initialize(cfg)
	srv.ListenAndServe()
}
