package repository

import (
	"BankingLedgerSystem/internal/config"
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoProviderInterface defines contract for MongoDB operations
type MongoProviderInterface interface {
	GetDB() *mongo.Database
	GetCollection(name string) *mongo.Collection
	Close(ctx context.Context) error
}

// mongoClient wraps the actual MongoDB client + db
type mongoClient struct {
	client *mongo.Client
	db     *mongo.Database
}

// MongoProvider implements MongoProviderInterface
type MongoProvider struct {
	cfg    *config.Config
	client *mongo.Client
	db     *mongo.Database
	once   sync.Once
}

// NewMongoProvider creates a new instance
func NewMongoProvider(cfg *config.Config) *MongoProvider {
	return &MongoProvider{cfg: cfg}
}

// InitMongo connects to MongoDB and verifies connection
func (m *MongoProvider) InitMongo(ctx context.Context) error {
	var err error

	m.once.Do(func() {
		clientOpts := options.Client().
			ApplyURI(m.cfg.MongoURL).
			SetConnectTimeout(10 * time.Second).
			SetServerSelectionTimeout(10 * time.Second)

		m.client, err = mongo.Connect(ctx, clientOpts)
		if err != nil {
			err = fmt.Errorf("failed to connect to Mongo: %w", err)
			return
		}

		// Ping to ensure connection
		if pingErr := m.client.Ping(ctx, nil); pingErr != nil {
			err = fmt.Errorf("failed to ping Mongo: %w", pingErr)
			return
		}

		m.db = m.client.Database(m.cfg.MongoDB)
	})

	return err
}

// GetDB returns the connected Mongo database
func (m *MongoProvider) GetDB() *mongo.Database {
	return m.db
}

// GetCollection returns a collection reference
func (m *MongoProvider) GetCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

// Close disconnects the client
func (m *MongoProvider) Close(ctx context.Context) error {
	if m.client != nil {
		return m.client.Disconnect(ctx)
	}
	return nil
}
