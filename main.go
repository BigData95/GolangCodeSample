package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"main/api/restful/dependencies"
	"main/api/restful/routers"
	"main/internal/adapters/repositories/repoCloudTask"
	"main/internal/adapters/repositories/repoEntities"
	"main/internal/adapters/repositories/repoPubSub"
	"main/internal/adapters/repositories/repoStorage"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	ctx := context.Background()

	mongo := dependencies.NewMongoManager()
	trackingRepo := &repoEntities.TrackingRepositoryImpl{}
	if err := mongo.Initialize(ctx, trackingRepo); err != nil {
		panic(fmt.Sprintf("failed to initialize MongoDB manager: %v", err))
	}
	defer mongo.Cleanup(ctx)

	pubSubManager := dependencies.NewPubSubManager()
	trackingPublisher := &repoPubSub.TrackingPublisherImpl{}
	if err := pubSubManager.Initialize(ctx, trackingPublisher); err != nil {
		panic(fmt.Sprintf("failed to initialize PubSub manager: %v", err))
	}
	defer pubSubManager.Cleanup(ctx)

	cloudTasksManager := dependencies.NewCloudTasksManager()
	cloudTasks := &repoCloudTask.CloudTasksImpl{}
	if err := cloudTasksManager.Initialize(ctx, cloudTasks); err != nil {
		panic(fmt.Sprintf("failed to initialize Cloud Tasks manager: %v", err))
	}
	defer cloudTasksManager.Cleanup(ctx)

	storageManager := dependencies.NewStorageManager()
	trackingStorage := &repoStorage.TrackingStorageImpl{}
	if err := storageManager.Initialize(ctx, trackingStorage); err != nil {
		panic(fmt.Sprintf("failed to initialize Storage manager: %v", err))
	}
	defer storageManager.Cleanup(ctx)

	app := routers.NewApp(router, trackingRepo, trackingPublisher, trackingStorage, cloudTasks)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "ready")
	})
	err := app.Router.Run()
	if err != nil {
		return
	}
}
