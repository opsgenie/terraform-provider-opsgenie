package maintenance

import (
	"net/http"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

type CreateRequest struct {
	client.BaseRequest
	Description string `json:"description"`
	Time        Time   `json:"time"`
	Rules       []Rule `json:"rules"`
}

func (r *CreateRequest) Validate() error {
	err := validateTime(r.Time)
	if err != nil {
		return err
	}
	if len(r.Rules) == 0 {
		return errors.New("There should be at least one rule.")
	}
	err = validateRules(r.Rules)
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateRequest) ResourcePath() string {
	return "/v1/maintenance"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type GetRequest struct {
	client.BaseRequest
	Id string
}

func (r *GetRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Maintenance ID cannot be blank.")
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {
	return "/v1/maintenance/" + r.Id
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

type UpdateRequest struct {
	client.BaseRequest
	Id          string
	Description string `json:"description"`
	Time        Time   `json:"time"`
	Rules       []Rule `json:"rules"`
}

func (r *UpdateRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Maintenance ID cannot be blank.")
	}
	err := validateTime(r.Time)
	if err != nil {
		return err
	}
	if len(r.Rules) == 0 {
		return errors.New("There should be at least one rule.")
	}
	err = validateRules(r.Rules)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateRequest) ResourcePath() string {
	return "/v1/maintenance/" + r.Id
}

func (r *UpdateRequest) Method() string {
	return http.MethodPut
}

type ChangeEndDateRequest struct {
	client.BaseRequest
	Id      string
	EndDate *time.Time `json:"endDate"`
}

func (r *ChangeEndDateRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Maintenance ID cannot be blank.")
	}
	if r.EndDate == nil {
		return errors.New("Maintenance End Date cannot be blank.")
	}
	return nil
}

func (r *ChangeEndDateRequest) ResourcePath() string {
	return "/v1/maintenance/" + r.Id + "/change-end-date"
}

func (r *ChangeEndDateRequest) Method() string {
	return http.MethodPost
}

type DeleteRequest struct {
	client.BaseRequest
	Id string
}

func (r *DeleteRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Maintenance ID cannot be blank.")
	}
	return nil
}

func (r *DeleteRequest) ResourcePath() string {
	return "/v1/maintenance/" + r.Id
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

type ListRequest struct {
	client.BaseRequest
	Type StatusType
}

func (r *ListRequest) Validate() error {
	err := validateStatusType(r.Type)
	if err != nil {
		return err
	}
	return nil
}

func (r *ListRequest) ResourcePath() string {
	return "/v1/maintenance"
}

func (r *ListRequest) Method() string {
	return http.MethodGet
}

type CancelRequest struct {
	client.BaseRequest
	Id string
}

func (r *CancelRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Maintenance ID cannot be blank.")
	}
	return nil
}

func (r *CancelRequest) ResourcePath() string {
	return "/v1/maintenance/" + r.Id + "/cancel"
}

func (r *CancelRequest) Method() string {
	return http.MethodPost
}

type Rule struct {
	State  RuleState `json:"state"`
	Entity Entity    `json:"entity"`
}

type Entity struct {
	Id   string         `json:"id"`
	Type RuleEntityType `json:"type"`
}

type Time struct {
	Type      TimeType   `json:"type"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}

type TimeType string
type RuleEntityType string
type RuleState string
type StatusType string

const (
	For5Minutes  TimeType = "for-5-minutes"
	For30Minutes TimeType = "for-30-minutes"
	For1Hour     TimeType = "for-1-hour"
	Indefinitely TimeType = "indefinitely"
	Schedule     TimeType = "schedule"

	Integration RuleEntityType = "integration"
	Policy      RuleEntityType = "policy"

	Enabled  RuleState = "enabled"
	Disabled RuleState = "disabled"

	All        StatusType = "all"
	NonExpired StatusType = "non-expired"
	Past       StatusType = "past"
)

func validateTime(t Time) error {
	switch t.Type {
	case For5Minutes, For30Minutes, For1Hour, Indefinitely, Schedule:
		break
	default:
		return errors.New("Time.Type should be one of these: " +
			"'For5Minutes', 'For30Minutes', 'For1Hour', 'Indefinitely' and 'Schedule'")
	}
	if t.Type == Schedule {
		if t.EndDate == nil || t.StartDate == nil {
			return errors.New("For 'Schedule' type both 'StartDate' and 'EndDate'" +
				" fields cannot be empty.")
		}
		sub := t.EndDate.Sub(*t.StartDate).Minutes()
		if sub <= 0 {
			return errors.New("EndDate should be after the StartDate.")
		}
	}
	return nil
}

func validateRules(rules []Rule) error {
	for _, rule := range rules {
		if rule.Entity.Type != Policy && rule.Entity.Type != Integration {
			return errors.New("Rule.Entity.Id should be one of these: " +
				"'Policy', 'Integration'.")
		}
		if rule.State != Enabled && rule.State != Disabled && rule.State != "" {
			return errors.New("Rule.State field should be one of these: " +
				"'Enabled' and 'Disabled' or empty.")
		}
		if rule.Entity.Type != Integration && rule.State == "" {
			return errors.New("Rule.State field cannot be empty " +
				"when the Rule.Entity.Type is not 'Integration'.")
		}
		if rule.Entity.Type == Integration && rule.State == Enabled {
			return errors.New("Rule.State field cannot be 'Enabled'" +
				" when the Rule.Entity.Type is 'Integration'.")
		}
	}
	return nil
}

func validateStatusType(status StatusType) error {
	switch status {
	case All, NonExpired, Past, "":
		return nil
	}
	return errors.New("Priority should be one of these: " +
		"'All', 'NonExpired' and 'Past' or empty.")
}
