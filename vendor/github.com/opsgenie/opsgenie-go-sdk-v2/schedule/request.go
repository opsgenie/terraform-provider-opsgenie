package schedule

import (
	"net/http"
	"strconv"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
)

type Identifier uint32

type CreateRequest struct {
	client.BaseRequest
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Timezone    string        `json:"timezone,omitempty"`
	Enabled     *bool         `json:"enabled,omitempty"`
	OwnerTeam   *og.OwnerTeam `json:"ownerTeam,omitempty"`
	Rotations   []og.Rotation `json:"rotations,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("Name cannot be empty.")
	}
	err := og.ValidateRotations(r.Rotations)
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateRequest) ResourcePath() string {
	return "/v2/schedules"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type GetRequest struct {
	client.BaseRequest
	IdentifierType  Identifier
	IdentifierValue string
}

func (r *GetRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {

	return "/v2/schedules/" + r.IdentifierValue
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

func (r *GetRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type UpdateRequest struct {
	client.BaseRequest
	IdentifierType  Identifier
	IdentifierValue string
	Name            string        `json:"name, omitempty"`
	Description     string        `json:"description,omitempty"`
	Timezone        string        `json:"timezone,omitempty"`
	Enabled         *bool         `json:"enabled,omitempty"`
	OwnerTeam       *og.OwnerTeam `json:"ownerTeam,omitempty"`
	Rotations       []og.Rotation `json:"rotations,omitempty"`
}

func (r *UpdateRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}
	err = og.ValidateRotations(r.Rotations)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateRequest) ResourcePath() string {

	return "/v2/schedules/" + r.IdentifierValue
}

func (r *UpdateRequest) Method() string {
	return http.MethodPatch
}

func (r *UpdateRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type DeleteRequest struct {
	client.BaseRequest
	IdentifierType  Identifier
	IdentifierValue string
}

func (r *DeleteRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeleteRequest) ResourcePath() string {

	return "/v2/schedules/" + r.IdentifierValue
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type ListRequest struct {
	client.BaseRequest
	Expand *bool
}

func (r *ListRequest) Validate() error {
	return nil
}

func (r *ListRequest) ResourcePath() string {

	return "/v2/schedules"
}

func (r *ListRequest) Method() string {
	return http.MethodGet
}

func (r *ListRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if *r.Expand {
		params["expand"] = "rotation"

	}

	return params
}

type GetTimelineRequest struct {
	client.BaseRequest
	IdentifierType  Identifier
	IdentifierValue string
	Expands         []ExpandType
	Interval        int
	IntervalUnit    Unit
	Date            *time.Time
}

func (r *GetTimelineRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}

	if r.IntervalUnit != Days && r.IntervalUnit != Months && r.IntervalUnit != Weeks {
		return errors.New("Provided InternalUnit is not valid.")
	}
	return nil
}

func (r *GetTimelineRequest) ResourcePath() string {

	return "/v2/schedules/" + r.IdentifierValue + "/timeline"

}

func (r *GetTimelineRequest) Method() string {
	return http.MethodGet
}

func (r *GetTimelineRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	if len(r.Expands) != 0 {
		expands := ""
		for i, expand := range r.Expands {
			if i != len(r.Expands)-1 {
				expands = expands + string(expand) + ","
			} else {
				expands = expands + string(expand)
			}
		}
		params["expand"] = expands
	}

	if r.Interval > 1 {
		params["interval"] = strconv.Itoa(r.Interval)
	}
	params["intervalUnit"] = string(r.IntervalUnit)

	if r.Date != nil {
		params["date"] = r.Date.Format("2006-01-02T15:04:05.000Z")
	}
	return params
}

func (r *GetTimelineRequest) WithExpands(expands ...ExpandType) GetTimelineRequest {
	r.Expands = expands
	return *r
}

type ExportScheduleRequest struct {
	client.BaseRequest
	IdentifierType   Identifier
	IdentifierValue  string
	ExportedFilePath string
}

func (r *ExportScheduleRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExportScheduleRequest) Method() string {
	return http.MethodGet
}

func (r *ExportScheduleRequest) getFileName() string {
	return r.IdentifierValue + ".ics"
}

func (r *ExportScheduleRequest) ResourcePath() string {
	return "/v2/schedules/" + r.getFileName()
}

func (r *ExportScheduleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type Unit string

const (
	Months Unit = "months"
	Weeks  Unit = "weeks"
	Days   Unit = "days"
)

type ExpandType string

const (
	Base       ExpandType = "base"
	Forwarding ExpandType = "forwarding"
	Override   ExpandType = "override"
)

const (
	Name Identifier = iota
	Id
)

func (r *CreateRequest) WithRotation(rotation *og.Rotation) *CreateRequest {
	r.Rotations = append(r.Rotations, *rotation)
	return r
}

func (r *UpdateRequest) WithRotation(rotation *og.Rotation) *UpdateRequest {
	r.Rotations = append(r.Rotations, *rotation)
	return r
}

func validateIdentifier(identifier string) error {
	if identifier == "" {
		return errors.New("Schedule identifier cannot be empty.")
	}
	return nil
}

//schedule rotation
type CreateRotationRequest struct {
	*og.Rotation
	ScheduleIdentifierType  Identifier
	ScheduleIdentifierValue string
}

func (r *CreateRotationRequest) Validate() error {
	err := validateIdentifier(r.ScheduleIdentifierValue)
	if err != nil {
		return err
	}

	err = r.Rotation.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateRotationRequest) ResourcePath() string {
	return "/v2/schedules/" + r.ScheduleIdentifierValue + "/rotations"

}

func (r *CreateRotationRequest) Method() string {
	return http.MethodPost
}

func (r *CreateRotationRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type GetRotationRequest struct {
	client.BaseRequest
	ScheduleIdentifierType  Identifier
	ScheduleIdentifierValue string
	RotationId              string
}

func (r *GetRotationRequest) Validate() error {

	err := validateIdentifier(r.ScheduleIdentifierValue)
	if err != nil {
		return err
	}

	if r.RotationId == "" {
		return errors.New("Rotation Id cannot be empty.")
	}

	return nil
}

func (r *GetRotationRequest) ResourcePath() string {
	return "/v2/schedules/" + r.ScheduleIdentifierValue + "/rotations/" + r.RotationId

}

func (r *GetRotationRequest) Method() string {
	return http.MethodGet
}

func (r *GetRotationRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type UpdateRotationRequest struct {
	ScheduleIdentifierType  Identifier
	ScheduleIdentifierValue string
	RotationId              string
	*og.Rotation
}

func (r *UpdateRotationRequest) Validate() error {

	err := validateIdentifier(r.ScheduleIdentifierValue)
	if err != nil {
		return err
	}

	if r.RotationId == "" {
		return errors.New("Rotation Id cannot be empty.")
	}

	return nil
}

func (r *UpdateRotationRequest) ResourcePath() string {

	return "/v2/schedules/" + r.ScheduleIdentifierValue + "/rotations/" + r.RotationId

}

func (r *UpdateRotationRequest) Method() string {
	return http.MethodPatch
}

func (r *UpdateRotationRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type DeleteRotationRequest struct {
	client.BaseRequest
	ScheduleIdentifierType  Identifier
	ScheduleIdentifierValue string
	RotationId              string
}

func (r *DeleteRotationRequest) Validate() error {

	err := validateIdentifier(r.ScheduleIdentifierValue)
	if err != nil {
		return err
	}

	if r.RotationId == "" {
		return errors.New("Rotation Id cannot be empty.")
	}

	return nil
}

func (r *DeleteRotationRequest) ResourcePath() string {

	return "/v2/schedules/" + r.ScheduleIdentifierValue + "/rotations/" + r.RotationId

}

func (r *DeleteRotationRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteRotationRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type ListRotationsRequest struct {
	client.BaseRequest
	ScheduleIdentifierType  Identifier
	ScheduleIdentifierValue string
}

func (r *ListRotationsRequest) Validate() error {

	err := validateIdentifier(r.ScheduleIdentifierValue)
	if err != nil {
		return err
	}

	return nil
}

func (r *ListRotationsRequest) ResourcePath() string {

	return "/v2/schedules/" + r.ScheduleIdentifierValue + "/rotations"

}

func (r *ListRotationsRequest) Method() string {
	return http.MethodGet
}

func (r *ListRotationsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}
