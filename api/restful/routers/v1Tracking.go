package routers

import (
	"github.com/gin-gonic/gin"
	"main/api/restful/controllers"
	"main/api/restful/middlewares"
	"main/internal/adapters/repositories/repoCloudTask"
	"main/internal/adapters/repositories/repoEntities"
	"main/internal/adapters/repositories/repoPubSub"
	"main/internal/adapters/repositories/repoStorage"
)

func NewApp(
	router *gin.Engine,
	trackingRepository repoEntities.TrackingRepositoryInterface,
	publisher repoPubSub.PublisherInterface,
	storage repoStorage.StorageInterface,
	cloudTasks repoCloudTask.CloudTasksInterface,
) *controllers.App {
	api := &controllers.App{
		Router: router,
	}
	v1 := api.Router.Group("/v1")
	{
		v1.POST("/pubsub/shipment-tracking",
			middlewares.TrackingRepositoryMiddleware(trackingRepository),
			middlewares.PublisherMiddleware(publisher),
			middlewares.StorageMiddleware(storage),
			middlewares.CloudTasksMiddleware(cloudTasks),
			api.CreateTrackingUpdate)
	}
	return api
}
