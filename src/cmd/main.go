package main

import (
	"clean_api/src/api/routers"
	"clean_api/src/config"
	"clean_api/src/data/db"
	"clean_api/src/pkg/logger"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.GetConfig()
	_, err := db.InitDB(cfg)
	defer db.CloseDB()
	if err != nil {
		log.Fatal(err)
	}
	logger := logger.New(os.Stdout, logger.LevelInfo)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      routers.Router(),
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
		IdleTimeout:  time.Minute,
	}
	logger.PrintInfo("server started", map[string]string{"port": srv.Addr})
	log.Fatal(srv.ListenAndServe())

}
