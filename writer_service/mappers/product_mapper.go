package mappers

import (
	kafkaMessages "github.com/mishnit/cqrs-microservices/proto/kafka"
	"github.com/mishnit/cqrs-microservices/writer_service/internal/models"
	writerService "github.com/mishnit/cqrs-microservices/writer_service/proto/product_writer"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductToGrpcMessage(product *models.Product) *kafkaMessages.Product {
	return &kafkaMessages.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func ProductFromGrpcMessage(product *kafkaMessages.Product) (*models.Product, error) {

	proUUID, err := uuid.FromString(product.GetProductID())
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ProductID:   proUUID,
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       product.GetPrice(),
		CreatedAt:   product.GetCreatedAt().AsTime(),
		UpdatedAt:   product.GetUpdatedAt().AsTime(),
	}, nil
}

func WriterProductToGrpc(product *models.Product) *writerService.Product {
	return &writerService.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}
