package queries

import (
	"context"
	"github.com/mishnit/cqrs-microservices/pkg/logger"
	"github.com/mishnit/cqrs-microservices/reader_service/config"
	"github.com/mishnit/cqrs-microservices/reader_service/internal/models"
	"github.com/mishnit/cqrs-microservices/reader_service/internal/product/repository"
)

type SearchProductHandler interface {
	Handle(ctx context.Context, query *SearchProductQuery) (*models.ProductsList, error)
}

type searchProductHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.Repository
	redisRepo repository.CacheRepository
}

func NewSearchProductHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.Repository, redisRepo repository.CacheRepository) *searchProductHandler {
	return &searchProductHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo}
}

func (s *searchProductHandler) Handle(ctx context.Context, query *SearchProductQuery) (*models.ProductsList, error) {
	return s.mongoRepo.Search(ctx, query.Text, query.Pagination)
}
