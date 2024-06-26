package middlewares

import (
	"github.com/gin-gonic/gin"
	"main/internal/adapters/repositories/repoCloudTask"
	"main/internal/adapters/repositories/repoEntities"
	"main/internal/adapters/repositories/repoPubSub"
	"main/internal/adapters/repositories/repoStorage"
)

func TrackingRepositoryMiddleware(trackingRepository repoEntities.TrackingRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("trackingRepository", trackingRepository)
		c.Next()
	}
}

func PublisherMiddleware(publisher repoPubSub.PublisherInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("publisher", publisher)
		c.Next()
	}
}

func StorageMiddleware(storage repoStorage.StorageInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage", storage)
		c.Next()
	}
}

func CloudTasksMiddleware(cloudTasks repoCloudTask.CloudTasksInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cloudTasks", cloudTasks)
		c.Next()
	}
}
