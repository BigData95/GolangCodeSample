package dependencies

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main/internal/adapters/repositories/repoEntities"
	"os"
	"sync"
)

type MongoManager struct {
	client   *mongo.Client
	initOnce sync.Once
}

func NewMongoManager() *MongoManager {
	return &MongoManager{}
}

func GetMongoClient(ctx context.Context) (*mongo.Client, error) {
	var err error
	var mongoClientOnce sync.Once
	var mongoClient *mongo.Client

	mongoClientOnce.Do(func() {
		MongoURI := os.Getenv("MONGO_DB_PROTOCOL") + "://" + os.Getenv("MONGO_DB_USER") + ":" + os.Getenv("MONGO_DB_PASSWORD") + "@" + os.Getenv("MONGO_DB_DNS") + "/shippers?readPreference=secondaryPreferred"
		mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(MongoURI))
		if err != nil {
			log.Fatalf("Mongo db failed: %v", err)
		}
	})
	return mongoClient, err
}

func (m *MongoManager) Initialize(ctx context.Context, repository repoEntities.BaseRepositoryInterface) error {
	var err error
	m.initOnce.Do(func() {
		m.client, err = GetMongoClient(ctx)
		if err != nil {
			log.Fatalf("Mongo db failed: %v", err)
		}
		repository.Initialize(ctx, m.client)
	})
	return err
}

func (m *MongoManager) Cleanup(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
