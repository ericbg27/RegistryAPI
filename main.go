package main

import (
	"log"

	"github.com/ericbg27/RegistryAPI/api"
	"github.com/ericbg27/RegistryAPI/db"
	"github.com/ericbg27/RegistryAPI/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v\n", err)
	}

	dbManager, err := db.CreateDBManager(config.DBSource)
	if err != nil {
		log.Fatalf("Cannot create database manager: %v\n", err)
	}

	server, err := api.NewServer(dbManager, config)
	if err != nil {
		log.Fatalf("Cannot create server: %v\n", err)
	}

	server.Start()
}
