package internal

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Details struct {
	Code string
	Desc string
}

type ErrorCode struct {
	MISSING_PARAMETER Details
	INVALID_PARAMETER Details
	NOT_FOUND         Details
	DEFAULT_ERROR     Details
	ALREADY_EXISTS    Details
}

type ErrorDetails struct {
	Parameter string `json:"parameter"`
	Code      string `json:"code"`
	Details   string `json:"details"`
	In        string `json:"in"`
	Message   string `json:"message"`
	Value     string `json:"value"`
	Ide       string `json:"ide"`
}

type ErrorsModel struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details"`
	Location  string `json:"location"`
	Parameter string `json:"parameter"`
}

type OutputModel struct {
	Version   string                 `default:"1.0.0" json:"version"`
	Status    string                 `json:"status"`
	ProcessID string                 `json:"processID"`
	Timestamp int64                  `json:"timestamp"`
	Errors    []ErrorsModel          `json:"errors,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func NewResponseController() *OutputModel {
	responseController := new(OutputModel)
	responseController.ProcessID = uuid.New().String()
	return responseController
}

func (o *OutputModel) SetResponse(status string) {
	o.Version = "1.0.0"
	o.Status = status
	o.Timestamp = time.Now().UnixMilli()
}

func (o *OutputModel) AddError(error ErrorsModel) {
	o.Errors = append(o.Errors, error)
}

func GetPubSubMessage(ctx *gin.Context, logger *logrus.Entry) []byte {
	var pubSubBody map[string]interface{}
	err := ctx.ShouldBindJSON(&pubSubBody)
	if err != nil {
		logger.Errorf("Failed to read request body from pubsub: %v", err)
	}
	b64String := pubSubBody["message"].(map[string]interface{})["data"].(string)
	b64String += strings.Repeat("=", (4-len(b64String)%4)%4)
	rawDecodedText, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		logger.Errorf("Failed to decode string 64: %v", err)
	}
	return rawDecodedText

}
