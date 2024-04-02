package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"net/http"
)

const (
	AnalyticsMetricsMinXscVersion = "1.7.1"
	xscEventApi                   = "api/v1/event"
)

type AnalyticsEventService struct {
	client     *jfroghttpclient.JfrogHttpClient
	XscDetails auth.ServiceDetails
}

func NewAnalyticsEventService(client *jfroghttpclient.JfrogHttpClient) *AnalyticsEventService {
	return &AnalyticsEventService{client: client}
}

// GetXscDetails returns the Xsc details
func (vs *AnalyticsEventService) GetXscDetails() auth.ServiceDetails {
	return vs.XscDetails
}

// AddGeneralEvent add general event in Xsc and returns msi generated by Xsc.
func (vs *AnalyticsEventService) AddGeneralEvent(event XscAnalyticsGeneralEvent) (string, error) {
	httpDetails := vs.XscDetails.CreateHttpClientDetails()
	requestContent, err := json.Marshal(event)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	resp, body, err := vs.client.SendPost(vs.XscDetails.GetUrl()+xscEventApi, requestContent, &httpDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusCreated); err != nil {
		return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	var response XscAnalyticsGeneralEventResponse
	err = json.Unmarshal(body, &response)
	return response.MultiScanId, errorutils.CheckError(err)
}

// UpdateGeneralEvent update finalized analytics metrics info of an existing event.
func (vs *AnalyticsEventService) UpdateGeneralEvent(event XscAnalyticsGeneralEventFinalize) error {
	httpDetails := vs.XscDetails.CreateHttpClientDetails()
	requestContent, err := json.Marshal(event)
	if err != nil {
		return errorutils.CheckError(err)
	}
	resp, body, err := vs.client.SendPut(vs.XscDetails.GetUrl()+xscEventApi, requestContent, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	return nil
}

// GetGeneralEvent returns event's data matching the provided multi scan id.
func (vs *AnalyticsEventService) GetGeneralEvent(msi string) (*XscAnalyticsGeneralEvent, error) {
	httpDetails := vs.XscDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(fmt.Sprintf("%s%s/%s", vs.XscDetails.GetUrl(), xscEventApi, msi), true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	var response XscAnalyticsGeneralEvent
	err = json.Unmarshal(body, &response)
	return &response, errorutils.CheckError(err)
}

// XscAnalyticsGeneralEvent extend the basic struct with Frogbot related info.
type XscAnalyticsGeneralEvent struct {
	XscAnalyticsBasicGeneralEvent
	GitInfo       *services.XscGitInfoContext
	IsGitInfoFlow bool `json:"is_gitinfo_flow,omitempty"`
}

type XscAnalyticsGeneralEventFinalize struct {
	XscAnalyticsBasicGeneralEvent
	MultiScanId string `json:"multi_scan_id,omitempty"`
}

type XscAnalyticsBasicGeneralEvent struct {
	EventType              EventType   `json:"event_type,omitempty"`
	EventStatus            EventStatus `json:"event_status,omitempty"`
	Product                ProductName `json:"product,omitempty"`
	ProductVersion         string      `json:"product_version,omitempty"`
	TotalFindings          int         `json:"total_findings,omitempty"`
	TotalIgnoredFindings   int         `json:"total_ignored_findings,omitempty"`
	IsDefaultConfig        bool        `json:"is_default_config,omitempty"`
	JfrogUser              string      `json:"jfrog_user,omitempty"`
	OsPlatform             string      `json:"os_platform,omitempty"`
	OsArchitecture         string      `json:"os_architecture,omitempty"`
	MachineId              string      `json:"machine_id,omitempty"`
	AnalyzerManagerVersion string      `json:"analyzer_manager_version,omitempty"`
	JpdVersion             string      `json:"jpd_version,omitempty"`
	TotalScanDuration      string      `json:"total_scan_duration,omitempty"`
	FrogbotSourceMsi       string      `json:"frogbot_source_msi,omitempty"`
	FrogbotScanType        string      `json:"frogbot_scan_type,omitempty"`
	FrogbotCiProvider      string      `json:"frogbot_ci_provider,omitempty"`
}

type XscAnalyticsGeneralEventResponse struct {
	MultiScanId string `json:"multi_scan_id,omitempty"`
}

type EventStatus string

const (
	Started   EventStatus = "started"
	Completed EventStatus = "completed"
	Cancelled EventStatus = "cancelled"
	Failed    EventStatus = "failed"
)

type ProductName string

const (
	CliProduct     ProductName = "cli"
	FrogbotProduct ProductName = "frogbot"
)

type EventType int

const (
	CliEventType EventType = 1
	FrogbotType  EventType = 8
)
