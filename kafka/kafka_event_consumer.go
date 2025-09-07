package kafka

import (
	"BankingLedgerSystem/internal/config"
	"BankingLedgerSystem/repository"
	"context"

	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
)

// TransactionEvent represents what producer published

type TransactionEvent struct {
	AccountNumber string  `json:"accountNumber"`
	Type          string  `json:"type"` // deposit | withdraw
	Amount        float64 `json:"amount"`
	Timestamp     int64   `json:"timestamp,omitempty"`
	Status        string  `json:"status,omitempty"` // success | failed
	Error         string  `json:"error,omitempty"`  // failure reason
}

// Consumer handles reading messages and processing
type Consumer struct {
	reader    *kafka.Reader
	db        *repository.ConnectionProvider
	topic     string
	retries   int
	workers   int
	mongo     *repository.MongoProvider
	mongoColl string
}

// NewConsumer initializes consumer with config + DB
func NewConsumer(cfg *config.Config, db *repository.ConnectionProvider) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        cfg.Kafka.Brokers,
			Topic:          cfg.Kafka.Topic,
			GroupID:        "banking-ledger-consumer", // consumer group
			CommitInterval: 0,                         // disable auto-commit
		}),
		db:      db,
		topic:   cfg.Kafka.Topic,
		retries: 3,                    // retry count
		workers: cfg.ProcessorWorkers, // parallelism
	}
}

// processTransaction applies business logic + DB write
func (c *Consumer) processTransaction(ctx context.Context, ev *TransactionEvent) error {
	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Example: lock row to prevent race conditions
	var balance float64
	err = tx.QueryRow(ctx,
		`SELECT balance FROM accounts WHERE account_number = $1 FOR UPDATE`,
		ev.AccountNumber).Scan(&balance)
	if err != nil {
		return err
	}

	switch ev.Type {
	case "deposit":
		balance += ev.Amount
	case "withdraw":
		if balance < ev.Amount {
			return errors.New("insufficient funds")
		}
		balance -= ev.Amount
	default:
		return errors.New("unknown transaction type")
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = $1 WHERE account_number = $2`,
		balance, ev.AccountNumber)
	if err != nil {
		return err
	}

	// commit only after success
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// Start runs the consumer loop
func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("failed to fetch message: %v", err)
			continue
		}

		var ev TransactionEvent
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			log.Printf("invalid message: %v", err)
			ev.Status = "failed"
			ev.Error = "invalid json"
			c.logTransaction(ctx, ev)
			_ = c.reader.CommitMessages(ctx, m)
			continue
		}

		// retry loop
		var procErr error
		for attempt := 1; attempt <= c.retries; attempt++ {
			procErr = c.processTransaction(ctx, &ev)
			if procErr == nil {
				break
			}
			log.Printf("processing failed (attempt %d/%d): %v", attempt, c.retries, procErr)
			time.Sleep(time.Second * time.Duration(attempt))
		}

		if procErr != nil {
			ev.Status = "failed"
			ev.Error = procErr.Error()
			log.Printf("dropping message after retries: %v", procErr)
		} else {
			ev.Status = "success"
		}

		// always write to Mongo (audit trail)
		c.logTransaction(ctx, ev)

		// commit offset only after we’re done
		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("commit failed: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

// logTransaction writes the event to MongoDB (success or fail)
func (c *Consumer) logTransaction(ctx context.Context, ev TransactionEvent) {
	coll := c.mongo.GetCollection(c.mongoColl)
	_, err := coll.InsertOne(ctx, bson.M{
		"accountNumber": ev.AccountNumber,
		"type":          ev.Type,
		"amount":        ev.Amount,
		"timestamp":     time.Now().Unix(),
		"status":        ev.Status,
		"error":         ev.Error,
	})
	if err != nil {
		log.Printf("failed to insert log in Mongo: %v", err)
	}
}
