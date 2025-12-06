package main

import (
	"log"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/server"
	"github.com/laksanagusta/identity/pkg/database"
	"github.com/laksanagusta/identity/pkg/logger"
)

func main() {
	log.Println("Starting App")

	// Get Environtment Config
	// env := os.Getenv("config")
	// if env == "" {
	// 	env = "local"
	// }

	// Load YAML Config
	cfg, err := config.LoadConfigV2()
	if err != nil {
		log.Panicf("Error loading config: %s", err)
	}

	// Init Logger
	appLogger, err := logger.InitLogger(cfg)
	if err != nil {
		log.Panicf("Error on initializing logger: %s", err)
	}
	defer appLogger.Sync()

	// Init DB Connection
	db, err := database.GetPostgreConnection(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		appLogger.Fatal("Error on getting database postgre connection: %s", err)
	}

	// Create & Run Server
	server := server.NewServer(cfg, appLogger, db)
	if err = server.Run(); err != nil {
		log.Panic(err)
	}

	appLogger.Info("App stopped")
}
