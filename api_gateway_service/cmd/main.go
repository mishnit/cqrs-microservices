package main

import (
	"flag"
	"github.com/mishnit/cqrs-microservices/api_gateway_service/config"
	"github.com/mishnit/cqrs-microservices/api_gateway_service/internal/server"
	"github.com/mishnit/cqrs-microservices/pkg/logger"
	"log"
)

// @contact.name Nitin Mishra
// @contact.url https://github.com/mishnit
// @contact.email geekymishnit@gmail.com
func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("ApiGateway")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
