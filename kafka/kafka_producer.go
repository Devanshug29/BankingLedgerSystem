package kafka

import (
	"BankingLedgerSystem/internal/config"
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

//go:generate sh -c "sh $(git rev-parse --show-toplevel)/scripts/mock_generator.sh $GOFILE"

type ProducerInterface interface {
	Publish(ctx context.Context, key, payload []byte) error
}

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(cfg *config.Config) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.Brokers...),
			Topic:    cfg.Kafka.Topic,
			Balancer: &kafka.Hash{}, // ensures same key (account) goes to same partition
		},
		topic: cfg.Kafka.Topic,
	}
}

func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,   // e.g. accountNumber
		Value: value, // JSON of deposit/withdraw event
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("failed to write message: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
