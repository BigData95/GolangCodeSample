package handlers

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"main/internal"
)

func AlreadyCreatedError(createdAt int64, uuid, shipperId string, logger *logrus.Entry, response *internal.OutputModel)(int){
		date := time.UnixMilli(createdAt).String()
		logger.Error("Already Created: " + uuid + "Tenant: " + shipperId)
		response.AddError(internal.ErrorsModel{
			Code: "ALREADY_EXISTS",
			Message:   "The resource you are trying to create already exists.",
			Details: "The tracking update you are trying to create, already exists. " +
				" Time: " + date,
			Location:  "request.body",
			Parameter: "shipperTrackingId",
		})
		response.SetResponse("CONFLICT")
		return http.StatusOK
}