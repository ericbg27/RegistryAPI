package main

import (
	"log"

	"github.com/ericbg27/RegistryAPI/api"
	"github.com/ericbg27/RegistryAPI/db"
	"github.com/ericbg27/RegistryAPI/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v\n", err)
	}

	dbConn, err := gorm.Open(postgres.Open(config.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot open DB connection: %v\n", err)
	}

	dbManager := db.NewDBManager(dbConn)

	server, err := api.NewServer(dbManager, config)
	if err != nil {
		log.Fatalf("Cannot create server: %v\n", err)
	}

	server.Start()
}
