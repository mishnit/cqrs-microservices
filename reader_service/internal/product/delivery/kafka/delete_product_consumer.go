package kafka

import (
	"context"
	"github.com/mishnit/cqrs-microservices/pkg/tracing"
	kafkaMessages "github.com/mishnit/cqrs-microservices/proto/kafka"
	"github.com/mishnit/cqrs-microservices/reader_service/internal/product/commands"
	"github.com/avast/retry-go"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (s *readerMessageProcessor) processProductDeleted(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.DeleteProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "readerMessageProcessor.processProductDeleted")
	defer span.Finish()

	msg := &kafkaMessages.ProductDeleted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	productUUID, err := uuid.FromString(msg.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	command := commands.NewDeleteProductCommand(productUUID)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.ps.Commands.DeleteProduct.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteProduct.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}

	s.commitMessage(ctx, r, m)
}
