package main

import (
	"fmt"
	"net/http"

	"app/config"
	"app/controller"
	"app/pkg/logger"
	"app/storage/postgresql"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	var loggerLevel = new(string)

	*loggerLevel = logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		*loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger("tu_go_admin_api_gateway", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()

	store, err := postgresql.NewConnectPostgresql(&cfg)
	if err != nil {
		log.Panic("Error connect to postgresql: ", logger.Error(err))
		return
	}

	defer store.CloseDB()

	newController := controller.NewController(&cfg, store)

	http.HandleFunc("/book", newController.BookController)

	fmt.Println("Listening Server", cfg.ServerHost+cfg.ServerPort)
	err = http.ListenAndServe(cfg.ServerHost+cfg.ServerPort, nil)
	if err != nil {
		log.Panic("Error listening server:", logger.Error(err))
		return
	}
}
