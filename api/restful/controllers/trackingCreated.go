package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"main/internal"
	"main/internal/adapters/handlers"
	"main/internal/adapters/repositories/repoCloudTask"
	"main/internal/adapters/repositories/repoEntities"
	"main/internal/adapters/repositories/repoPubSub"
	"main/internal/adapters/repositories/repoStorage"
	"main/internal/core/domain"
	"main/internal/core/services"
	"net/http"
)

type App struct {
	Router *gin.Engine
}

func (app *App) CreateTrackingUpdate(ctx *gin.Context) {
	ResponseController := internal.NewResponseController()
	logger := handlers.NewLogger(ResponseController.ProcessID)

	rawDecodedText := internal.GetPubSubMessage(ctx, logger)
	var body domain.TrackingModel

	err := json.Unmarshal(rawDecodedText, &body)
	if err != nil {
		logger.Error("Error Unmarshall base64: %v", err)
	}
	validate := validator.New()
	err = validate.Struct(body)
	if err != nil {
		logger.Error("Validation Error")
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		for _, validationError := range validationErrors {
			ResponseController.AddError(internal.ErrorsModel{
				Code:      "TrackingCreation",
				Message:   validationError.Error(),
				Parameter: validationError.Field(),
			},
			)
		}
		ctx.JSON(http.StatusOK, ResponseController)
		return
	}
	printableBody, err := json.Marshal(body)
	if err != nil {
		logger.Infof("Failed to marshal input body, struct body: %v", body)
	} else {
		logger.Infof("Input TrackingCreated: Body: %v ", string(printableBody))
	}

	trackingRepo := ctx.MustGet("trackingRepository").(repoEntities.TrackingRepositoryInterface)
	publisher := ctx.MustGet("publisher").(repoPubSub.PublisherInterface)
	storage := ctx.MustGet("storage").(repoStorage.StorageInterface)
	cloudTask := ctx.MustGet("cloudTasks").(repoCloudTask.CloudTasksInterface)
	service := services.TrackingCreator{
		TrackingDb: trackingRepo,
		Publisher:  publisher,
		Storage:    storage,
		CloudTask:  cloudTask,
		Logger:     logger,
		Response:   ResponseController,
	}
	status := service.MainProcess(body)
	ctx.JSON(status, ResponseController)
}
