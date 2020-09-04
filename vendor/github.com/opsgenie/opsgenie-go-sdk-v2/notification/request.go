package notification

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
)

type CreateRuleStepRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
	Contact        og.Contact    `json:"contact"`
	SendAfter      *og.SendAfter `json:"sendAfter,omitempty"`
	Enabled        *bool         `json:"enabled,omitempty"`
}

func (r *CreateRuleStepRequest) Validate() error {
	err := validateRuleIdentifier(r.UserIdentifier, r.RuleId)
	if err != nil {
		return err
	}

	err = validateContact(&r.Contact)
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateRuleStepRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/steps"
}

func (r *CreateRuleStepRequest) Method() string {
	return http.MethodPost
}

type GetRuleStepRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
	RuleStepId     string
}

func (r *GetRuleStepRequest) Validate() error {
	err := validateRuleStepIdentifier(r.UserIdentifier, r.RuleId, r.RuleStepId)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetRuleStepRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/steps/" + r.RuleStepId
}

func (r *GetRuleStepRequest) Method() string {
	return http.MethodGet
}

type UpdateRuleStepRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
	RuleStepId     string
	Contact        *og.Contact   `json:"contact,omitempty"`
	SendAfter      *og.SendAfter `json:"sendAfter,omitempty"`
	Enabled        *bool         `json:"enabled,omitempty"`
}

func (r *UpdateRuleStepRequest) Validate() error {
	err := validateRuleStepIdentifier(r.UserIdentifier, r.RuleId, r.RuleStepId)
	if err != nil {
		return err
	}
	if r.Contact != nil {
		err = validateContact(r.Contact)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UpdateRuleStepRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/steps/" + r.RuleStepId
}

func (r *UpdateRuleStepRequest) Method() string {
	return http.MethodPatch
}

type DeleteRuleStepRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
	RuleStepId     string
}

func (r *DeleteRuleStepRequest) Validate() error {
	err := validateRuleStepIdentifier(r.UserIdentifier, r.RuleId, r.RuleStepId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeleteRuleStepRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/steps/" + r.RuleStepId
}

func (r *DeleteRuleStepRequest) Method() string {
	return http.MethodDelete
}

type ListRuleStepsRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
}

func (r *ListRuleStepsRequest) Validate() error {
	err := validateRuleIdentifier(r.UserIdentifier, r.RuleId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ListRuleStepsRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/steps"
}

func (r *ListRuleStepsRequest) Method() string {
	return http.MethodGet
}

type EnableRuleStepRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
	RuleStepId     string
}

func (r *EnableRuleStepRequest) Validate() error {
	err := validateRuleStepIdentifier(r.UserIdentifier, r.RuleId, r.RuleStepId)
	if err != nil {
		return err
	}
	return nil
}

func (r *EnableRuleStepRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/steps/" + r.RuleStepId + "/enable"
}

func (r *EnableRuleStepRequest) Method() string {
	return http.MethodPost
}

type DisableRuleStepRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
	RuleStepId     string
}

func (r *DisableRuleStepRequest) Validate() error {
	err := validateRuleStepIdentifier(r.UserIdentifier, r.RuleId, r.RuleStepId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DisableRuleStepRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/steps/" + r.RuleStepId + "/disable"
}

func (r *DisableRuleStepRequest) Method() string {
	return http.MethodPost
}

type CreateRuleRequest struct {
	client.BaseRequest
	UserIdentifier   string
	Name             string                 `json:"name"`
	ActionType       ActionType             `json:"actionType"`
	Criteria         *og.Criteria           `json:"criteria,omitempty"`
	NotificationTime []NotificationTimeType `json:"notificationTime,omitempty"`
	TimeRestriction  *og.TimeRestriction    `json:"timeRestriction,omitempty"`
	Schedules        []Schedule             `json:"schedules,omitempty"`
	Steps            []*og.Step             `json:"steps,omitempty"`
	Order            uint32                 `json:"order,omitempty"`
	Repeat           *Repeat                `json:"repeat,omitempty"`
	Enabled          *bool                  `json:"enabled,omitempty"`
}

func (r *CreateRuleRequest) Validate() error {
	if r.UserIdentifier == "" {
		return errors.New("User identifier cannot be empty.")
	}
	if r.Name == "" {
		return errors.New("Name cannot be empty.")
	}
	if r.ActionType == "" {
		return errors.New("Action type cannot be empty.")
	}
	if (r.ActionType == ScheduleStart || r.ActionType == ScheduleEnd) && len(r.NotificationTime) == 0 {
		return errors.New("Notification time cannot be empty.")
	}
	if len(r.Schedules) != 0 {
		for _, schedule := range r.Schedules {
			err := validateSchedule(schedule)
			if err != nil {
				return err
			}
		}
	}
	if len(r.Steps) != 0 {
		for _, step := range r.Steps {
			err := validateStep(step, r.ActionType)
			if err != nil {
				return err
			}
		}
	}
	if r.Criteria != nil {
		err := og.ValidateCriteria(*r.Criteria)
		if err != nil {
			return err
		}
	}

	if r.TimeRestriction != nil {
		err := og.ValidateRestrictions(r.TimeRestriction)
		if err != nil {
			return err
		}
	}

	if r.Repeat != nil && r.Repeat.LoopAfter <= 0 {
		return errors.New("Loop after must have a positive integer value.")
	}

	return nil
}

func (r *CreateRuleRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules"
}

func (r *CreateRuleRequest) Method() string {
	return http.MethodPost
}

type GetRuleRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
}

func (r *GetRuleRequest) Validate() error {
	err := validateRuleIdentifier(r.UserIdentifier, r.RuleId)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetRuleRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId
}

func (r *GetRuleRequest) Method() string {
	return http.MethodGet
}

type UpdateRuleRequest struct {
	client.BaseRequest
	UserIdentifier   string
	RuleId           string
	Criteria         *og.Criteria           `json:"criteria,omitempty"`
	NotificationTime []NotificationTimeType `json:"notificationTime,omitempty"`
	TimeRestriction  *og.TimeRestriction    `json:"timeRestriction,omitempty"`
	Schedules        []Schedule             `json:"schedules,omitempty"`
	Steps            []*og.Step             `json:"steps,omitempty"`
	Order            uint32                 `json:"order,omitempty"`
	Repeat           *Repeat                `json:"repeat,omitempty"`
	Enabled          *bool                  `json:"enabled,omitempty"`
}

func (r *UpdateRuleRequest) Validate() error {
	err := validateRuleIdentifier(r.UserIdentifier, r.RuleId)
	if err != nil {
		return err
	}
	if len(r.Schedules) != 0 {
		for _, schedule := range r.Schedules {
			err := validateSchedule(schedule)
			if err != nil {
				return err
			}
		}
	}

	if len(r.Steps) != 0 {
		for _, step := range r.Steps {
			err := validateStepWithoutActionTypeInfo(step)
			if err != nil {
				return err
			}
		}
	}
	if r.Criteria != nil {
		err := og.ValidateCriteria(*r.Criteria)
		if err != nil {
			return err
		}
	}

	if r.TimeRestriction != nil {
		err := og.ValidateRestrictions(r.TimeRestriction)
		if err != nil {
			return err
		}
	}

	if r.Repeat != nil && r.Repeat.LoopAfter <= 0 {
		return errors.New("Loop after must have a positive integer value.")
	}
	return nil
}

func (r *UpdateRuleRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId
}

func (r *UpdateRuleRequest) Method() string {
	return http.MethodPatch
}

type DeleteRuleRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
}

func (r *DeleteRuleRequest) Validate() error {
	err := validateRuleIdentifier(r.UserIdentifier, r.RuleId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeleteRuleRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId
}

func (r *DeleteRuleRequest) Method() string {
	return http.MethodDelete
}

type ListRuleRequest struct {
	client.BaseRequest
	UserIdentifier string
}

func (r *ListRuleRequest) Validate() error {
	if r.UserIdentifier == "" {
		return errors.New("User identifier cannot be empty.")
	}
	return nil
}

func (r *ListRuleRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules"
}

func (r *ListRuleRequest) Method() string {
	return http.MethodGet
}

type EnableRuleRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
}

func (r *EnableRuleRequest) Validate() error {
	err := validateRuleIdentifier(r.UserIdentifier, r.RuleId)
	if err != nil {
		return err
	}
	return nil
}

func (r *EnableRuleRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/enable"
}

func (r *EnableRuleRequest) Method() string {
	return http.MethodPost
}

type DisableRuleRequest struct {
	client.BaseRequest
	UserIdentifier string
	RuleId         string
}

func (r *DisableRuleRequest) Validate() error {
	err := validateRuleIdentifier(r.UserIdentifier, r.RuleId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DisableRuleRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/" + r.RuleId + "/disable"
}

func (r *DisableRuleRequest) Method() string {
	return http.MethodPost
}

type CopyNotificationRulesRequest struct {
	client.BaseRequest
	UserIdentifier string
	ToUsers        []string    `json:"toUsers"`
	RuleTypes      []RuleTypes `json:"ruleTypes"`
}

func (r *CopyNotificationRulesRequest) Validate() error {
	if r.UserIdentifier == "" {
		return errors.New("User identifier cannot be empty.")
	}
	if len(r.ToUsers) == 0 {
		return errors.New("You must specify a list of the users which you want to copy the rules to.")
	}
	if len(r.RuleTypes) == 0 {
		return errors.New("Specify a list of the action types you want to copy the rules of.")
	}
	return nil
}

func (r *CopyNotificationRulesRequest) ResourcePath() string {

	return "/v2/users/" + r.UserIdentifier + "/notification-rules/copy-to"
}

func (r *CopyNotificationRulesRequest) Method() string {
	return http.MethodPost
}

func validateRuleIdentifier(userIdentifier string, ruleIdentifier string) error {
	if userIdentifier == "" {
		return errors.New("User identifier cannot be empty.")
	}
	if ruleIdentifier == "" {
		return errors.New("Rule identifier cannot be empty.")

	}
	return nil
}

func validateRuleStepIdentifier(userIdentifier string, ruleIdentifier string, ruleStepId string) error {
	err := validateRuleIdentifier(userIdentifier, ruleIdentifier)
	if err != nil {
		return err
	}
	if ruleStepId == "" {
		return errors.New("Rule Step identifier cannot be empty.")

	}
	return nil
}

func validateContact(contact *og.Contact) error {
	if contact == nil {
		return errors.New("Contact cannot be empty.")

	}
	if contact.To == "" {
		return errors.New("To cannot be empty.")
	}
	if contact.MethodOfContact == "" {
		return errors.New("Method cannot be empty.")

	}
	return nil
}

type ActionType string

const (
	CreateAlert         ActionType = "create-alert"
	AcknowledgedAlert   ActionType = "acknowledged-alert"
	ClosedAlert         ActionType = "closed-alert"
	AssignedAlert       ActionType = "assigned-alert"
	AddNote             ActionType = "add-note"
	ScheduleStart       ActionType = "schedule-start"
	ScheduleEnd         ActionType = "schedule-end"
	IncomingCallRouting ActionType = "incoming-call-routing"
)

type NotificationTimeType string

const (
	JustBefore        NotificationTimeType = "just-before"
	FifteenMinutesAgo NotificationTimeType = "15-minutes-ago"
	OneHourAgo        NotificationTimeType = "1-hour-ago"
	OneDayAgo         NotificationTimeType = "1-day-ago"
)

type Schedule struct {
	TypeOfSchedule string `json:"type"`
	Name           string `json:"name,omitempty"`
	Id             string `json:"id,omitempty"`
}

func validateSchedule(schedule Schedule) error {
	if schedule.TypeOfSchedule != "schedule" {
		return errors.New("Type of schedule must be schedule.")
	}
	return nil
}

type Repeat struct {
	LoopAfter uint32 `json:"loopAfter,omitempty"`
	Enabled   *bool  `json:"enabled,omitempty"`
}

func validateStep(step *og.Step, actionType ActionType) error {
	if step.Contact.To == "" {
		return errors.New("To cannot be empty.")
	}
	if step.Contact.MethodOfContact == "" {
		return errors.New("Method cannot be empty.")
	}
	if (actionType == CreateAlert || actionType == AssignedAlert) && step.SendAfter == nil {
		return errors.New("SendAfter cannot be empty.")
	}

	return nil
}

func validateStepWithoutActionTypeInfo(step *og.Step) error {
	if step.Contact.To == "" {
		return errors.New("To cannot be empty.")
	}
	if step.Contact.MethodOfContact == "" {
		return errors.New("Method cannot be empty.")
	}

	return nil
}

type RuleTypes string

const (
	All                   RuleTypes = "all"
	AcknowledgedAlertRule RuleTypes = "acknowledged-alert"
	RenotifiedAlertRule   RuleTypes = "renotified-alert"
	ClosedAlertRule       RuleTypes = "closed-alert"
	ScheduleStartRule     RuleTypes = "schedule-start"
	AssignedAlertRule     RuleTypes = "assigned-alert"
	AddNoteRule           RuleTypes = "add-note"
	NewAlertRule          RuleTypes = "new-alert"
)
