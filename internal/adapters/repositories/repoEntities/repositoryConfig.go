package repoEntities

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryConfig struct {
	Db  *mongo.Collection
	Ctx context.Context
}

type BaseRepositoryInterface interface {
	Initialize(ctx context.Context, client *mongo.Client)
}
