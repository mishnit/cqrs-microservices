package main

import (
	"flag"
	"github.com/mishnit/cqrs-microservices/pkg/logger"
	"github.com/mishnit/cqrs-microservices/reader_service/config"
	"github.com/mishnit/cqrs-microservices/reader_service/internal/server"
	"log"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("ReaderService")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
