package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"main/internal/adapters/repositories/repoStorage"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"main/internal"
	"main/internal/core/domain"
)

type ContextProcessor struct {
	HasEvidences           bool
	HasImageEvidences      bool
	HasAnomaly             bool
	HasDescription         bool
	TrackingId             string
	TenantId               string
	DeduplicationTimestamp string
	IsMissingEvidences     bool
	Storage                repoStorage.StorageInterface
	Evidences              []domain.Evidences
	Token                  string
	Logger                 *logrus.Entry
}

func NewContextProcessor(
	eventContext *domain.Context,
	tenantId, trackingId, deduplicationTimestamp string,
	storage repoStorage.StorageInterface,
	logger *logrus.Entry,
) *ContextProcessor {
	if eventContext == nil {
		logger.Error("eventContext is nil")
		return nil
	}
	return &ContextProcessor{
		Logger:                 logger,
		HasEvidences:           len(eventContext.Evidences) > 1,
		HasAnomaly:             len(eventContext.AnomalyType) > 0,
		DeduplicationTimestamp: deduplicationTimestamp,
		TrackingId:             trackingId,
		TenantId:               tenantId,
		IsMissingEvidences:     false,
		Storage:                storage,
	}

}
func (c *ContextProcessor) setEvidenceName(evidenceType, index string) string {
	return c.TrackingId + ":" + c.DeduplicationTimestamp + ":" + evidenceType + ":" + index
}

type PartnerResponse struct {
	Data struct {
		Evidence    string `json:"evidence"`
		ContentType string `json:"content_type"`
	} `json:"data"`
}
type RequestBody struct {
	ShipperTenantId  string `json:"shipper_tenant_id"`
	TrackingId       string `json:"tracking_id"`
	StorageReference string `json:"storage_reference"`
}

func (c *ContextProcessor) getBlobFromPartners(storageReference string, signedJwt string) (blob, contentType string, err error) {
	requestBody := []byte(
		fmt.Sprintf(
			`{"shipper_tenant_id": "%s","tracking_id": "%s","storage_reference": "%s"}`,
			c.TenantId,
			c.TrackingId,
			storageReference,
		))

	req, err := http.NewRequest("POST", internal.URLS["EVIDENCES"], bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+signedJwt)

	client := &http.Client{}
	startTime := time.Now().UTC().UnixMilli()
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	endTime := time.Now().UTC().UnixMilli()
	c.Logger.Infof("[communication][request] service: shippers-tracking-process | endpoint:POST:%v | %v | %v",
		internal.URLS["EVIDENCES"], startTime, endTime)

	dataResult, ok := result["data"].(map[string]interface{})
	if !ok {
		c.Logger.Errorf("Failed response from partners, response: %v", result)
		return "", "", errors.New(fmt.Sprintf("failed to get response: %w", result))
	}

	blob, ok = dataResult["evidence"].(string)
	if !ok {
		return "", "", errors.New("missing evidence in response")
	}

	contentType, ok = dataResult["content_type"].(string)
	if !ok {
		return "", "", errors.New("missing content_type in response")
	}

	return blob, contentType, nil

}

func (c *ContextProcessor) processImageEvidence(
	evidenceType string,
	evidence domain.Evidences,
	index int,
	signedJwt string,
) (fileReference, completeFileName string, missingEvidence bool) {
	fileName := c.setEvidenceName(evidenceType, strconv.FormatInt(int64(index), 10))
	storageReference := evidence.Value

	blob, contentType, err := c.getBlobFromPartners(storageReference, signedJwt)
	if err != nil {
		c.Logger.Errorf("Error while getting the blob: %w", err)
		return "", "", true
	}

	fileReference, fileExtension, err := c.Storage.UploadToStorage(fileName, blob, contentType)
	if err != nil {
		c.Logger.Error("failed to upload file to shippers repoStorage")
		return "", "", true
	}

	completeFileName = fmt.Sprintf("%s.%s", fileName, fileExtension)
	return fileReference, completeFileName, false

}

type ProcessImagesResponse struct {
	FileReference   string
	FileName        string
	MissingEvidence bool
	EvidenceIndex   int
	SignedJwt       string
}

func (c *ContextProcessor) ProcessEvidences(eventContext *domain.Context) {
	wg := new(sync.WaitGroup)
	once := &sync.Once{}
	evidenceResponse := make(chan ProcessImagesResponse)
	var signedJwt string
	var err error
	for index, evidence := range eventContext.Evidences {
		index := index
		evidence := evidence
		evidenceType := evidence.Type
		if evidenceType == internal.EvidenceTypes["IMAGE"] || evidenceType == internal.EvidenceTypes["SIGNATURE"] {
			wg.Add(1)
			go func(index int, evidenceType string) {
				defer wg.Done()
				once.Do(func() {
					signedJwt, err = GenerateJWT()
					if err != nil {
						c.Logger.Infof("error Generating JWT: %v", err)
					}
				})
				fileReference, fileName, missingEvidence := c.processImageEvidence(evidenceType, evidence, index, signedJwt)
				evidenceResponse <- ProcessImagesResponse{
					FileReference:   fileReference,
					FileName:        fileName,
					MissingEvidence: missingEvidence,
					EvidenceIndex:   index,
					SignedJwt:       signedJwt,
				}
			}(index, evidenceType)
		}
	}
	go func() {
		wg.Wait()
		close(evidenceResponse)
	}()
	c.Evidences = eventContext.Evidences
	for response := range evidenceResponse {
		if response.MissingEvidence {
			c.IsMissingEvidences = true
		}
		c.Evidences[response.EvidenceIndex].EvidenceSource = response.FileReference
		c.Evidences[response.EvidenceIndex].FileName = response.FileName
		c.Token = response.SignedJwt
	}
}
