package commands

import (
	"context"
	kafkaClient "github.com/mishnit/cqrs-microservices/pkg/kafka"
	"github.com/mishnit/cqrs-microservices/pkg/logger"
	"github.com/mishnit/cqrs-microservices/pkg/tracing"
	kafkaMessages "github.com/mishnit/cqrs-microservices/proto/kafka"
	"github.com/mishnit/cqrs-microservices/writer_service/config"
	"github.com/mishnit/cqrs-microservices/writer_service/internal/models"
	"github.com/mishnit/cqrs-microservices/writer_service/internal/product/repository"
	"github.com/mishnit/cqrs-microservices/writer_service/mappers"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type CreateProductCmdHandler interface {
	Handle(ctx context.Context, command *CreateProductCommand) error
}

type createProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        repository.Repository
	kafkaProducer kafkaClient.Producer
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, pgRepo repository.Repository, kafkaProducer kafkaClient.Producer) *createProductHandler {
	return &createProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *createProductHandler) Handle(ctx context.Context, command *CreateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createProductHandler.Handle")
	defer span.Finish()

	productDto := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price}

	product, err := c.pgRepo.CreateProduct(ctx, productDto)
	if err != nil {
		return err
	}

	msg := &kafkaMessages.ProductCreated{Product: mappers.ProductToGrpcMessage(product)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return c.kafkaProducer.PublishMessage(ctx, message)
}
