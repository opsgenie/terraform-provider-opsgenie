package schedule

import (
	"net/http"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type GetOnCallsRequest struct {
	client.BaseRequest
	Flat                   bool
	Date                   *time.Time
	ScheduleIdentifierType Identifier
	ScheduleIdentifier     string
}

func (r *GetOnCallsRequest) Validate() error {
	err := validateIdentifiers(r.ScheduleIdentifier, "Schedule identifier cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *GetOnCallsRequest) Method() string {
	return http.MethodGet
}

func (r *GetOnCallsRequest) ResourcePath() string {
	return "/v2/schedules/" + r.ScheduleIdentifier + "/on-calls"
}

func (r *GetOnCallsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}
	if r.Flat {
		params["flat"] = "true"
	}

	if r.Date != nil {
		params["date"] = r.Date.Format("2006-01-02T15:04:05.000Z")
	}

	return params
}

type GetNextOnCallsRequest struct {
	client.BaseRequest
	Flat                   bool
	Date                   *time.Time
	ScheduleIdentifierType Identifier
	ScheduleIdentifier     string
}

func (r *GetNextOnCallsRequest) Validate() error {
	err := validateIdentifiers(r.ScheduleIdentifier, "Schedule identifier cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *GetNextOnCallsRequest) Method() string {
	return http.MethodGet
}

func (r *GetNextOnCallsRequest) ResourcePath() string {
	return "/v2/schedules/" + r.ScheduleIdentifier + "/next-on-calls"
}

func (r *GetNextOnCallsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}
	if r.Flat {
		params["flat"] = "true"
	}

	if r.Date != nil {
		params["date"] = r.Date.Format("2006-01-02T15:04:05.000Z")
	}

	return params
}

type ExportOnCallUserRequest struct {
	client.BaseRequest
	UserIdentifier   string
	ExportedFilePath string
}

func (r *ExportOnCallUserRequest) Validate() error {
	err := validateIdentifiers(r.UserIdentifier, "User identifier cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *ExportOnCallUserRequest) Method() string {
	return http.MethodGet
}

func (r *ExportOnCallUserRequest) getFileName() string {
	return r.UserIdentifier + ".ics"
}

func (r *ExportOnCallUserRequest) ResourcePath() string {
	return "/v2/schedules/on-calls/" + r.getFileName()
}
