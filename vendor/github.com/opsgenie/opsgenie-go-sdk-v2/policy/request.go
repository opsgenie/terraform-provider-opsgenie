package policy

import (
	"errors"
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type CreateAlertPolicyRequest struct {
	client.BaseRequest
	MainFields
	Message                  string             `json:"message,omitempty"`
	Continue                 *bool              `json:"continue,omitempty"`
	Alias                    string             `json:"alias,omitempty"`
	AlertDescription         string             `json:"alertDescription,omitempty"`
	Entity                   string             `json:"entity,omitempty"`
	Source                   string             `json:"source,omitempty"`
	IgnoreOriginalDetails    *bool              `json:"ignoreOriginalDetails,omitempty"`
	Actions                  []string           `json:"actions,omitempty"`
	IgnoreOriginalActions    *bool              `json:"ignoreOriginalActions,omitempty"`
	Details                  []string           `json:"details,omitempty"`
	IgnoreOriginalResponders *bool              `json:"ignoreOriginalResponders,omitempty"`
	Responders               *[]alert.Responder `json:"responders,omitempty"`
	IgnoreOriginalTags       *bool              `json:"ignoreOriginalTags,omitempty"`
	Tags                     []string           `json:"tags,omitempty"`
	Priority                 alert.Priority     `json:"priority,omitempty"`
}

type CreateNotificationPolicyRequest struct {
	client.BaseRequest
	MainFields
	AutoRestartAction   *AutoRestartAction   `json:"autoRestartAction,omitempty"`
	AutoCloseAction     *AutoCloseAction     `json:"autoCloseAction,omitempty"`
	DeDuplicationAction *DeDuplicationAction `json:"deduplicationAction,omitempty"`
	DelayAction         *DelayAction         `json:"delayAction,omitempty"`
	Suppress            *bool                `json:"suppress,omitempty"`
}

type MainFields struct {
	PolicyType        string              `json:"type,omitempty"`
	Name              string              `json:"name,omitempty"`
	Enabled           *bool               `json:"enabled"`
	PolicyDescription string              `json:"policyDescription"`
	Filter            *og.Filter          `json:"filter,omitempty"`
	TimeRestriction   *og.TimeRestriction `json:"timeRestrictions,omitempty"`
	TeamId            string
}

func (r *CreateAlertPolicyRequest) Validate() error {
	err := ValidateMainFields(&r.MainFields)
	if err != nil {
		return err
	}
	if r.Message == "" {
		return errors.New("alert message cannot be empty")
	}
	if r.Responders != nil {
		err = ValidateResponders(r.Responders)
		if err != nil {
			return err
		}
	}
	if r.Priority != "" {
		err = alert.ValidatePriority(r.Priority)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CreateAlertPolicyRequest) ResourcePath() string {
	return "/v2/policies"
}

func (r *CreateAlertPolicyRequest) Method() string {
	return http.MethodPost
}

func (r *CreateAlertPolicyRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

func (r *CreateNotificationPolicyRequest) Validate() error {
	err := ValidateMainFields(&r.MainFields)
	if err != nil {
		return err
	}
	if r.TeamId == "" {
		return errors.New("policy team id should be provided")
	}
	if r.AutoRestartAction != nil {
		err = ValidateAutoRestartAction(*r.AutoRestartAction)
		if err != nil {
			return err
		}
	}
	if r.AutoCloseAction != nil {
		err = ValidateAutoCloseAction(*r.AutoCloseAction)
		if err != nil {
			return err
		}
	}
	if r.DeDuplicationAction != nil {
		err = ValidateDeDuplicationAction(*r.DeDuplicationAction)
		if err != nil {
			return err
		}
	}
	if r.DelayAction != nil {
		err = ValidateDelayAction(*r.DelayAction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CreateNotificationPolicyRequest) ResourcePath() string {
	return "/v2/policies"
}

func (r *CreateNotificationPolicyRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

func (r *CreateNotificationPolicyRequest) Method() string {
	return http.MethodPost
}

type GetAlertPolicyRequest struct {
	client.BaseRequest
	Id     string
	TeamId string
}

func (r *GetAlertPolicyRequest) Validate() error {
	if r.Id == "" {
		return errors.New("policy id should be provided")
	}
	return nil
}

func (r *GetAlertPolicyRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id
}

func (r *GetAlertPolicyRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

func (r *GetAlertPolicyRequest) Method() string {
	return http.MethodGet
}

type GetNotificationPolicyRequest struct {
	client.BaseRequest
	Id     string
	TeamId string
}

func (r *GetNotificationPolicyRequest) Validate() error {
	if r.Id == "" {
		return errors.New("policy id should be provided")
	}
	if r.TeamId == "" {
		return errors.New("policy team id should be provided")
	}
	return nil
}

func (r *GetNotificationPolicyRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id
}

func (r *GetNotificationPolicyRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

func (r *GetNotificationPolicyRequest) Method() string {
	return http.MethodGet
}

type UpdateAlertPolicyRequest struct {
	client.BaseRequest
	MainFields
	Message                  string                 `json:"message,omitempty"`
	Continue                 *bool                  `json:"continue,omitempty"`
	Alias                    string                 `json:"alias,omitempty"`
	AlertDescription         string                 `json:"alertDescription,omitempty"`
	Entity                   string                 `json:"entity,omitempty"`
	Source                   string                 `json:"source,omitempty"`
	IgnoreOriginalDetails    *bool                  `json:"ignoreOriginalDetails,omitempty"`
	Actions                  []string               `json:"actions,omitempty"`
	IgnoreOriginalActions    *bool                  `json:"ignoreOriginalActions,omitempty"`
	Details                  map[string]interface{} `json:"details,omitempty"`
	IgnoreOriginalResponders *bool                  `json:"ignoreOriginalResponders,omitempty"`
	Responders               *[]alert.Responder     `json:"responders,omitempty"`
	IgnoreOriginalTags       *bool                  `json:"ignoreOriginalTags,omitempty"`
	Tags                     []string               `json:"tags,omitempty"`
	Priority                 alert.Priority         `json:"priority,omitempty"`
	Id                       string
}

func (r *UpdateAlertPolicyRequest) Validate() error {
	err := ValidatePolicyIdentifier("alert", r.Id, r.TeamId)
	if err != nil {
		return err
	}
	err = ValidateMainFields(&r.MainFields)
	if err != nil {
		return err
	}
	if r.Message == "" {
		return errors.New("alert message cannot be empty")
	}
	if r.Responders != nil {
		err = ValidateResponders(r.Responders)
		if err != nil {
			return err
		}
	}
	if r.Priority != "" {
		err = alert.ValidatePriority(r.Priority)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UpdateAlertPolicyRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id
}

func (r *UpdateAlertPolicyRequest) RequestParams() map[string]string {
        params := make(map[string]string)
        params["teamId"] = r.TeamId
        return params
}

func (r *UpdateAlertPolicyRequest) Method() string {
	return http.MethodPut
}

type UpdateNotificationPolicyRequest struct {
	client.BaseRequest
	MainFields
	AutoRestartAction   *AutoRestartAction   `json:"autoRestartAction,omitempty"`
	AutoCloseAction     *AutoCloseAction     `json:"autoCloseAction,omitempty"`
	DeDuplicationAction *DeDuplicationAction `json:"deduplicationAction,omitempty"`
	DelayAction         *DelayAction         `json:"delayAction,omitempty"`
	Suppress            *bool                `json:"suppress,omitempty"`
	Id                  string
}

func (r *UpdateNotificationPolicyRequest) Validate() error {
	err := ValidatePolicyIdentifier("notification", r.Id, r.TeamId)
	if err != nil {
		return err
	}
	err = ValidateMainFields(&r.MainFields)
	if err != nil {
		return err
	}
	if r.TeamId == "" {
		return errors.New("policy team id should be provided")
	}
	if r.AutoRestartAction != nil {
		err = ValidateAutoRestartAction(*r.AutoRestartAction)
		if err != nil {
			return err
		}
	}
	if r.AutoCloseAction != nil {
		err = ValidateAutoCloseAction(*r.AutoCloseAction)
		if err != nil {
			return err
		}
	}
	if r.DeDuplicationAction != nil {
		err = ValidateDeDuplicationAction(*r.DeDuplicationAction)
		if err != nil {
			return err
		}
	}
	if r.DelayAction != nil {
		err = ValidateDelayAction(*r.DelayAction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UpdateNotificationPolicyRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id
}

func (r *UpdateNotificationPolicyRequest) RequestParams() map[string]string {
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

func (r *UpdateNotificationPolicyRequest) Method() string {
	return http.MethodPut
}

type DeletePolicyRequest struct {
	client.BaseRequest
	Id     string `json:"id,omitempty"`
	TeamId string
	Type   PolicyType
}

func (r *DeletePolicyRequest) Validate() error {
	if r.Type != AlertPolicy && r.Type != NotificationPolicy {
		return errors.New("policy type should be one of alert or notification")
	}
	err := ValidatePolicyIdentifier(string(r.Type), r.Id, r.TeamId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeletePolicyRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id
}

func (r *DeletePolicyRequest) Method() string {
	return http.MethodDelete
}

func (r *DeletePolicyRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

type DisablePolicyRequest struct {
	client.BaseRequest
	Id     string `json:"id,omitempty"`
	TeamId string
	Type   PolicyType
}

func (r *DisablePolicyRequest) Validate() error {
	if r.Type != AlertPolicy && r.Type != NotificationPolicy {
		return errors.New("policy type should be one of alert or notification")
	}
	err := ValidatePolicyIdentifier(string(r.Type), r.Id, r.TeamId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DisablePolicyRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id + "/disable"
}

func (r *DisablePolicyRequest) Method() string {
	return http.MethodPost
}

func (r *DisablePolicyRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

type EnablePolicyRequest struct {
	client.BaseRequest
	Id     string `json:"id,omitempty"`
	TeamId string
	Type   PolicyType
}

func (r *EnablePolicyRequest) Validate() error {
	if r.Type != AlertPolicy && r.Type != NotificationPolicy {
		return errors.New("policy type should be one of alert or notification")
	}
	err := ValidatePolicyIdentifier(string(r.Type), r.Id, r.TeamId)
	if err != nil {
		return err
	}
	return nil
}

func (r *EnablePolicyRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id + "/enable"
}

func (r *EnablePolicyRequest) Method() string {
	return http.MethodPost
}

func (r *EnablePolicyRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

type ChangeOrderRequest struct {
	client.BaseRequest
	Id          string `json:"id,omitempty"`
	TeamId      string
	Type        PolicyType
	TargetIndex int `json:"targetIndex,omitempty"`
}

func (r *ChangeOrderRequest) Validate() error {
	if r.Type != AlertPolicy && r.Type != NotificationPolicy {
		return errors.New("policy type should be one of alert or notification")
	}
	err := ValidatePolicyIdentifier(string(r.Type), r.Id, r.TeamId)
	if err != nil {
		return err
	}
	if r.TargetIndex < 0 {
		return errors.New("target index should be at least 0")
	}
	return nil
}

func (r *ChangeOrderRequest) ResourcePath() string {
	return "/v2/policies/" + r.Id + "/change-order"
}

func (r *ChangeOrderRequest) Method() string {
	return http.MethodPost
}

func (r *ChangeOrderRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

type ListAlertPoliciesRequest struct {
	client.BaseRequest
	TeamId string
}

func (r *ListAlertPoliciesRequest) Validate() error {
	return nil
}

func (r *ListAlertPoliciesRequest) ResourcePath() string {
	return "/v2/policies/alert"
}

func (r *ListAlertPoliciesRequest) Method() string {
	return http.MethodGet
}

func (r *ListAlertPoliciesRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

type ListNotificationPoliciesRequest struct {
	client.BaseRequest
	TeamId string
}

func (r *ListNotificationPoliciesRequest) Validate() error {
	if r.TeamId == "" {
		return errors.New("team id should be provided")
	}
	return nil
}

func (r *ListNotificationPoliciesRequest) ResourcePath() string {
	return "/v2/policies/notification"
}

func (r *ListNotificationPoliciesRequest) Method() string {
	return http.MethodGet
}

func (r *ListNotificationPoliciesRequest) RequestParams() map[string]string {
	if r.TeamId == "" {
		return nil
	}
	params := make(map[string]string)
	params["teamId"] = r.TeamId
	return params
}

type PolicyType string

type Duration struct {
	TimeAmount int         `json:"timeAmount,omitempty"`
	TimeUnit   og.TimeUnit `json:"timeUnit,omitempty"`
}

type AutoRestartAction struct {
	Duration       *Duration `json:"duration,omitempty"`
	MaxRepeatCount int       `json:"maxRepeatCount,omitempty"`
}

type AutoCloseAction struct {
	Duration *Duration `json:"duration,omitempty"`
}

type DeDuplicationAction struct {
	DeDuplicationActionType DeDuplicationActionType `json:"deduplicationActionType,omitempty"`
	Duration                *Duration               `json:"duration,omitempty"`
	Count                   int                     `json:"count,omitempty"`
}

type DelayAction struct {
	DelayOption DelayType `json:"delayOption,omitempty"`
	UntilMinute *int      `json:"untilMinute,omitempty"`
	UntilHour   *int      `json:"untilHour,omitempty"`
	Duration    *Duration `json:"duration,omitempty"`
}

type DeDuplicationActionType string
type DelayType string

const (
	ValueBased     DeDuplicationActionType = "value-based"
	FrequencyBased DeDuplicationActionType = "frequency-based"

	ForDuration        DelayType  = "for-duration"
	NextTime           DelayType  = "next-time"
	NextWeekday        DelayType  = "next-weekday"
	NextMonday         DelayType  = "next-monday"
	NextTuesday        DelayType  = "next-tuesday"
	NextWednesday      DelayType  = "next-wednesday"
	NextThursday       DelayType  = "next-thursday"
	NextFriday         DelayType  = "next-friday"
	NextSaturday       DelayType  = "next-saturday"
	NextSunday         DelayType  = "next-sunday"
	AlertPolicy        PolicyType = "alert"
	NotificationPolicy PolicyType = "notification"
)

func ValidateMainFields(fields *MainFields) error {
	if fields == nil {
		return errors.New("policy main fields should be provided")
	}
	if fields.PolicyType != "alert" && fields.PolicyType != "notification" {
		return errors.New("policy type should be alert or notification")
	}
	if fields.Name == "" {
		return errors.New("policy name cannot be empty")
	}
	if fields.Filter != nil {
		err := og.ValidateFilter(*fields.Filter)
		if err != nil {
			return err
		}
	}
	if fields.TimeRestriction != nil {
		err := og.ValidateRestrictions(fields.TimeRestriction)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateDuration(duration *Duration) error {
	if duration != nil && duration.TimeUnit != "" && duration.TimeUnit != og.Days && duration.TimeUnit != og.Hours && duration.TimeUnit != og.Minutes {
		return errors.New("timeUnit provided for duration is not valid")
	}
	if duration != nil && duration.TimeAmount <= 0 {
		return errors.New("duration timeAmount should be greater than zero")
	}
	if duration != nil && duration.TimeUnit == "" {
		duration.TimeUnit = og.Minutes
	}
	return nil
}

func ValidateDeDuplicationAction(action DeDuplicationAction) error {
	if action.DeDuplicationActionType != ValueBased && action.DeDuplicationActionType != FrequencyBased {
		return errors.New("deDuplication action type should be one of value-based or frequency-based")
	}
	if action.Duration != nil {
		err := ValidateDuration(action.Duration)
		if err != nil {
			return err
		}
	}
	if action.Count < 0 {
		return errors.New("deDuplication count is not valid")
	}
	return nil
}

func ValidateAutoRestartAction(action AutoRestartAction) error {
	if action.Duration == nil {
		return errors.New("autoRestart action duration cannot be empty")
	}
	err := ValidateDuration(action.Duration)
	if err != nil {
		return err
	}
	if action.MaxRepeatCount < 0 {
		return errors.New("autoRestart maxRepeatCount is not valid")
	}
	return nil
}

func ValidateAutoCloseAction(action AutoCloseAction) error {
	if action.Duration == nil {
		return errors.New("autoClose action duration cannot be empty")
	}
	err := ValidateDuration(action.Duration)
	if err != nil {
		return err
	}
	return nil
}

func ValidateDelayAction(action DelayAction) error {
	if action.DelayOption != ForDuration && action.DelayOption != NextTime && action.DelayOption != NextWeekday &&
		action.DelayOption != NextMonday && action.DelayOption != NextTuesday && action.DelayOption != NextWednesday &&
		action.DelayOption != NextThursday && action.DelayOption != NextFriday && action.DelayOption != NextSaturday && action.DelayOption != NextSunday {
		return errors.New("delay option should be one of for-duration, next-time, next-weekday, next-monday, next-tuesday, next-wednesday, next-thursday, next-friday, next-saturday, next-sunday")
	}
	if action.DelayOption == ForDuration {
		if action.Duration == nil {
			return errors.New("delayAction duration cannot be empty")
		}
		err := ValidateDuration(action.Duration)
		if err != nil {
			return err
		}
	}
	if action.DelayOption != ForDuration && ((*action.UntilHour < 0 || *action.UntilHour > 23) || (*action.UntilMinute < 0 || *action.UntilMinute > 59)) {
		return errors.New("delayAction's UntilHour or UntilMinute is not valid")
	}
	return nil
}

func ValidateResponders(responders *[]alert.Responder) error {
	for _, responder := range *responders {
		if responder.Type != alert.UserResponder && responder.Type != alert.TeamResponder {
			return errors.New("responder type for alert policy should be one of team or user")
		}
		if responder.Id == "" {
			return errors.New("responder id should be provided")
		}
	}
	return nil
}

func ValidatePolicyIdentifier(policyType string, id string, teamId string) error {
	if id == "" {
		return errors.New("policy id should be provided")
	}
	if "notification" == policyType && teamId == "" {
		return errors.New("policy team id should be provided")
	}
	return nil
}
