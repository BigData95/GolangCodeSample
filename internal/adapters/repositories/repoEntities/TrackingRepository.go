package repoEntities

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main/internal/core/domain"
)

type TrackingRepositoryInterface interface {
	BaseRepositoryInterface
	GetTrackingUpdates(trackingId string) ([]domain.TrackingDbModel, error)
	Persist(tracking domain.TrackingDbModel) error
}

type TrackingRepositoryImpl struct {
	config RepositoryConfig
}

func (t TrackingRepositoryImpl) Initialize(ctx context.Context, client *mongo.Client) {
	collection := client.Database("shippers").Collection("tracking-updates")
	t.config = RepositoryConfig{Db: collection, Ctx: ctx}
}

func (t TrackingRepositoryImpl) GetTrackingUpdates(trackingId string) ([]domain.TrackingDbModel, error) {
	var result []domain.TrackingDbModel
	var cursor *mongo.Cursor
	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{{"created_at", 1}})

	cursor, err := t.config.Db.Find(t.config.Ctx, bson.D{{"tracking_id", trackingId}}, queryOptions)
	if errors.Is(mongo.ErrNoDocuments, err) {
		fmt.Printf("No document was found with the trackingId %s\n", trackingId)
		return result, err
	}
	if err = cursor.All(t.config.Ctx, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (t TrackingRepositoryImpl) Persist(tracking domain.TrackingDbModel) error {
	_, err := t.config.Db.InsertOne(t.config.Ctx, tracking)
	return err
}
