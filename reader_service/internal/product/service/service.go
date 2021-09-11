package service

import (
	"github.com/mishnit/cqrs-microservices/pkg/logger"
	"github.com/mishnit/cqrs-microservices/reader_service/config"
	"github.com/mishnit/cqrs-microservices/reader_service/internal/product/commands"
	"github.com/mishnit/cqrs-microservices/reader_service/internal/product/queries"
	"github.com/mishnit/cqrs-microservices/reader_service/internal/product/repository"
)

type ProductService struct {
	Commands *commands.ProductCommands
	Queries  *queries.ProductQueries
}

func NewProductService(
	log logger.Logger,
	cfg *config.Config,
	mongoRepo repository.Repository,
	redisRepo repository.CacheRepository,
) *ProductService {

	createProductHandler := commands.NewCreateProductHandler(log, cfg, mongoRepo, redisRepo)
	deleteProductCmdHandler := commands.NewDeleteProductCmdHandler(log, cfg, mongoRepo, redisRepo)
	updateProductCmdHandler := commands.NewUpdateProductCmdHandler(log, cfg, mongoRepo, redisRepo)

	getProductByIdHandler := queries.NewGetProductByIdHandler(log, cfg, mongoRepo, redisRepo)
	searchProductHandler := queries.NewSearchProductHandler(log, cfg, mongoRepo, redisRepo)

	productCommands := commands.NewProductCommands(createProductHandler, updateProductCmdHandler, deleteProductCmdHandler)
	productQueries := queries.NewProductQueries(getProductByIdHandler, searchProductHandler)

	return &ProductService{Commands: productCommands, Queries: productQueries}
}
