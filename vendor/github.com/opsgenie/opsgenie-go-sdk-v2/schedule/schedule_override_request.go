package schedule

import (
	"errors"
	"net/http"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type RotationIdentifier struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type CreateScheduleOverrideRequest struct {
	client.BaseRequest
	Alias                  string               `json:"alias,omitempty"`
	User                   Responder            `json:"user,omitempty"`
	StartDate              time.Time            `json:"startDate,omitempty"`
	EndDate                time.Time            `json:"endDate,omitempty"`
	Rotations              []RotationIdentifier `json:"rotations,omitempty"`
	ScheduleIdentifierType Identifier
	ScheduleIdentifier     string
}

func (r *CreateScheduleOverrideRequest) Validate() error {
	err := validateIdentifiers(r.ScheduleIdentifier, "Schedule identifier cannot be empty.")
	if err != nil {
		return err
	}
	err = validateUser(&r.User)
	if err != nil {
		return err
	}
	err = validateDates(&r.StartDate, "Start date cannot be empty.")
	if err != nil {
		return err
	}
	err = validateDates(&r.EndDate, "End date cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateScheduleOverrideRequest) ResourcePath() string {

	return "/v2/schedules/" + r.ScheduleIdentifier + "/overrides"
}

func (r *CreateScheduleOverrideRequest) Method() string {
	return http.MethodPost
}

func (r *CreateScheduleOverrideRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type GetScheduleOverrideRequest struct {
	client.BaseRequest
	ScheduleIdentifierType Identifier
	ScheduleIdentifier     string
	Alias                  string
}

func (r *GetScheduleOverrideRequest) Validate() error {
	err := validateIdentifiers(r.ScheduleIdentifier, "Schedule identifier cannot be empty.")
	if err != nil {
		return err
	}
	err = validateIdentifiers(r.Alias, "Alias cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *GetScheduleOverrideRequest) ResourcePath() string {

	return "/v2/schedules/" + r.ScheduleIdentifier + "/overrides/" + r.Alias
}

func (r *GetScheduleOverrideRequest) Method() string {
	return http.MethodGet
}

func (r *GetScheduleOverrideRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type ListScheduleOverrideRequest struct {
	client.BaseRequest
	ScheduleIdentifierType Identifier
	ScheduleIdentifier     string
}

func (r *ListScheduleOverrideRequest) Validate() error {
	err := validateIdentifiers(r.ScheduleIdentifier, "Schedule identifier cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *ListScheduleOverrideRequest) ResourcePath() string {
	return "/v2/schedules/" + r.ScheduleIdentifier + "/overrides"
}

func (r *ListScheduleOverrideRequest) Method() string {
	return http.MethodGet
}

func (r *ListScheduleOverrideRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type DeleteScheduleOverrideRequest struct {
	client.BaseRequest
	ScheduleIdentifierType Identifier
	ScheduleIdentifier     string
	Alias                  string
}

func (r *DeleteScheduleOverrideRequest) Validate() error {
	err := validateIdentifiers(r.ScheduleIdentifier, "Schedule identifier cannot be empty.")
	if err != nil {
		return err
	}
	err = validateIdentifiers(r.Alias, "Alias cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *DeleteScheduleOverrideRequest) ResourcePath() string {

	return "/v2/schedules/" + r.ScheduleIdentifier + "/overrides/" + r.Alias
}

func (r *DeleteScheduleOverrideRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteScheduleOverrideRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

type UpdateScheduleOverrideRequest struct {
	client.BaseRequest
	Alias                  string
	User                   Responder            `json:"user,omitempty"`
	StartDate              time.Time            `json:"startDate,omitempty"`
	EndDate                time.Time            `json:"endDate,omitempty"`
	Rotations              []RotationIdentifier `json:"rotations,omitempty"`
	ScheduleIdentifierType Identifier
	ScheduleIdentifier     string
}

func (r *UpdateScheduleOverrideRequest) Validate() error {
	err := validateIdentifiers(r.ScheduleIdentifier, "Schedule identifier cannot be empty.")
	if err != nil {
		return err
	}
	err = validateIdentifiers(r.Alias, "Alias cannot be empty.")
	if err != nil {
		return err
	}
	err = validateUser(&r.User)
	if err != nil {
		return err
	}
	err = validateDates(&r.StartDate, "Start date cannot be empty.")
	if err != nil {
		return err
	}
	err = validateDates(&r.EndDate, "End date cannot be empty.")
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateScheduleOverrideRequest) ResourcePath() string {

	return "/v2/schedules/" + r.ScheduleIdentifier + "/overrides/" + r.Alias
}

func (r *UpdateScheduleOverrideRequest) Method() string {
	return http.MethodPut
}

func (r *UpdateScheduleOverrideRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.ScheduleIdentifierType == Name {
		params["scheduleIdentifierType"] = "name"
	} else {
		params["scheduleIdentifierType"] = "id"
	}

	return params
}

func validateIdentifiers(identifier string, message string) error {
	if identifier == "" {
		return errors.New(message)
	}
	return nil
}

func validateUser(user *Responder) error {
	if *user == (Responder{}) {
		return errors.New("User cannot be empty.")
	}
	return nil
}

func validateDates(date *time.Time, message string) error {
	if *date == (time.Time{}) {
		return errors.New(message)
	}
	return nil
}
