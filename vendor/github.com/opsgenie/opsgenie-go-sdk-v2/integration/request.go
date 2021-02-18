package integration

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
)

type GetRequest struct {
	client.BaseRequest
	Id string
}

func (r *GetRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Integration ID cannot be blank.")
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {
	return "/v2/integrations/" + r.Id
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

type listRequest struct {
	client.BaseRequest
}

func (r *listRequest) Validate() error {
	return nil
}

func (r *listRequest) ResourcePath() string {
	return "/v2/integrations"
}

func (r *listRequest) Method() string {
	return http.MethodGet
}

type APIBasedIntegrationRequest struct {
	client.BaseRequest
	Name                        string        `json:"name"`
	Type                        string        `json:"type"`
	AllowWriteAccess            *bool         `json:"allowWriteAccess"`
	IgnoreRespondersFromPayload *bool         `json:"ignoreRespondersFromPayload"`
	SuppressNotifications       *bool         `json:"suppressNotifications"`
	OwnerTeam                   *og.OwnerTeam `json:"ownerTeam,omitempty"`
	Responders                  []Responder   `json:"responders,omitempty"`
}

func (r *APIBasedIntegrationRequest) Validate() error {
	if r.Name == "" || r.Type == "" {
		return errors.New("Name and Type fields cannot be empty.")
	}
	err := validateResponders(r.Responders)
	if err != nil {
		return err
	}
	return nil
}

func (r *APIBasedIntegrationRequest) ResourcePath() string {
	return "/v2/integrations"
}

func (r *APIBasedIntegrationRequest) Method() string {
	return http.MethodPost
}

type WebhookIntegrationRequest struct {
	client.BaseRequest
	Name                  string            `json:"name"`
	Type                  string            `json:"type"`
	AllowWriteAccess      *bool             `json:"allowWriteAccess"`
	SuppressNotifications *bool             `json:"suppressNotifications"`
	OwnerTeam             *og.OwnerTeam     `json:"ownerTeam,omitempty"`
	Responders            []Responder       `json:"responders,omitempty"`
	WebhookUrl            string            `json:"url"`
	AddAlertDescription   *bool             `json:"addAlertDescription"`
	AddAlertDetails       *bool             `json:"addAlertDetails"`
	Headers               map[string]string `json:"headers,omitempty"`
}

func (r *WebhookIntegrationRequest) Validate() error {
	if r.Name == "" || r.Type == "" || r.WebhookUrl == "" {
		return errors.New("Name, Type and WebhookUrl fields cannot be empty.")
	}
	if r.Type != "Webhook" {
		return errors.New("Type has to be [Webhook] for Webhook integration.")
	}
	err := validateResponders(r.Responders)
	if err != nil {
		return err
	}
	return nil
}

func (r *WebhookIntegrationRequest) ResourcePath() string {
	return "/v2/integrations"
}

func (r *WebhookIntegrationRequest) Method() string {
	return http.MethodPost
}

type EmailBasedIntegrationRequest struct {
	client.BaseRequest
	Name                        string        `json:"name"`
	Type                        string        `json:"type"`
	EmailUsername               string        `json:"emailUsername"`
	IgnoreRespondersFromPayload *bool         `json:"ignoreRespondersFromPayload,omitempty"`
	SuppressNotifications       *bool         `json:"suppressNotifications,omitempty"`
	OwnerTeam                   *og.OwnerTeam `json:"ownerTeam,omitempty"`
	Responders                  []Responder   `json:"responders,omitempty"`
}

func (r *EmailBasedIntegrationRequest) Validate() error {
	if r.Name == "" || r.Type == "" || r.EmailUsername == "" {
		return errors.New("Name, Type and EmailUsername fields cannot be empty.")
	}
	err := validateResponders(r.Responders)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmailBasedIntegrationRequest) ResourcePath() string {
	return "/v2/integrations"
}

func (r *EmailBasedIntegrationRequest) Method() string {
	return http.MethodPost
}

type UpdateIntegrationRequest struct {
	client.BaseRequest
	Id                          string
	Name                        string
	Type                        string
	EmailUsername               string
	WebhookUrl                  string
	Enabled                     *bool
	IgnoreRespondersFromPayload *bool
	SuppressNotifications       *bool
	Responders                  []Responder
	AddAlertDescription         *bool
	AddAlertDetails             *bool
	Headers                     map[string]string
	OtherFields
}

type OtherFields map[string]interface{}

func (r OtherFields) Validate() error {

	if _, ok := r["id"]; !ok {
		return errors.New("Integration ID cannot be blank.")
	}
	if _, ok := r["name"]; !ok {
		return errors.New("Name field cannot be empty.")
	}
	if _, ok := r["type"]; !ok {
		return errors.New("Type field cannot be empty.")
	}
	if r["type"] == "Webhook" {
		if _, ok := r["url"]; !ok {
			return errors.New("[url] cannot be empty for type Webhook.")
		}
	}
	err := validateResponders(r["responders"].([]Responder))
	if err != nil {
		return err
	}
	return nil
}

func (r OtherFields) ResourcePath() string {
	return "/v2/integrations/" + r["id"].(string)
}

func (r OtherFields) Method() string {
	return http.MethodPut
}

func (r OtherFields) RequestParams() map[string]string {
	return nil
}

func (r OtherFields) Metadata(apiRequest client.ApiRequest) map[string]interface{} {
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json; charset=utf-8"

	return headers
}

type DeleteIntegrationRequest struct {
	client.BaseRequest
	Id string
}

func (r *DeleteIntegrationRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Integration ID cannot be blank.")
	}
	return nil
}

func (r *DeleteIntegrationRequest) ResourcePath() string {
	return "/v2/integrations/" + r.Id
}

func (r *DeleteIntegrationRequest) Method() string {
	return http.MethodDelete
}

type EnableIntegrationRequest struct {
	client.BaseRequest
	Id string
}

func (r *EnableIntegrationRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Integration ID cannot be blank.")
	}
	return nil
}

func (r *EnableIntegrationRequest) ResourcePath() string {
	return "/v2/integrations/" + r.Id + "/enable"
}

func (r *EnableIntegrationRequest) Method() string {
	return http.MethodPost
}

type DisableIntegrationRequest struct {
	client.BaseRequest
	Id string
}

func (r *DisableIntegrationRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Integration ID cannot be blank.")
	}
	return nil
}

func (r *DisableIntegrationRequest) ResourcePath() string {
	return "/v2/integrations/" + r.Id + "/disable"
}

func (r *DisableIntegrationRequest) Method() string {
	return http.MethodPost
}

type AuthenticateIntegrationRequest struct {
	client.BaseRequest
	Type string `json:"type"`
}

func (r *AuthenticateIntegrationRequest) Validate() error {
	if r.Type == "" {
		return errors.New("Type cannot be blank.")
	}
	return nil
}

func (r *AuthenticateIntegrationRequest) ResourcePath() string {
	return "/v2/integrations/authenticate"
}

func (r *AuthenticateIntegrationRequest) Method() string {
	return http.MethodPost
}

type GetIntegrationActionsRequest struct {
	client.BaseRequest
	Id string
}

func (r *GetIntegrationActionsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Type cannot be blank.")
	}
	return nil
}

func (r *GetIntegrationActionsRequest) ResourcePath() string {
	return "/v2/integrations/" + r.Id + "/actions"
}

func (r *GetIntegrationActionsRequest) Method() string {
	return http.MethodGet
}

type Filter struct {
	ConditionMatchType og.ConditionMatchType `json:"conditionMatchType,omitempty"`
	Conditions         []og.Condition        `json:"conditions,omitempty"`
}

type CreateIntegrationActionsRequest struct {
	client.BaseRequest
	Id                               string
	Type                             ActionType        `json:"type"`
	Name                             string            `json:"name"`
	Alias                            string            `json:"alias"`
	Order                            int               `json:"order,omitempty"`
	User                             string            `json:"user,omitempty"`
	Note                             string            `json:"note,omitempty"`
	Filter                           *Filter           `json:"filter,omitempty"`
	Source                           string            `json:"source,omitempty"`
	Message                          string            `json:"message,omitempty"`
	Description                      string            `json:"description,omitempty"`
	Entity                           string            `json:"entity,omitempty"`
	AppendAttachments                *bool             `json:"appendAttachments,omitempty"`
	AlertActions                     []string          `json:"alertActions,omitempty"`
	IgnoreAlertActionsFromPayload    *bool             `json:"ignoreAlertActionsFromPayload,omitempty"`
	IgnoreRespondersFromPayload      *bool             `json:"ignoreRespondersFromPayload,omitempty"`
	IgnoreTagsFromPayload            *bool             `json:"ignoreTagsFromPayload,omitempty"`
	IgnoreExtraPropertiesFromPayload *bool             `json:"ignoreExtraPropertiesFromPayload,omitempty"`
	Responders                       []Responder       `json:"responders,omitempty"`
	Tags                             []string          `json:"tags,omitempty"`
	ExtraProperties                  map[string]string `json:"extraProperties,omitempty"`
}

func (r *CreateIntegrationActionsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Integration ID cannot be blank.")
	}
	if r.Name == "" || r.Type == "" || r.Alias == "" {
		return errors.New("Name, Type and Alias fields cannot be empty.")
	}
	err := validateActionType(r.Type)
	if err != nil {
		return err
	}
	if r.Filter != nil {
		err = validateConditionMatchType(r.Filter.ConditionMatchType)
		if err != nil {
			return err
		}
		err = og.ValidateFilter(og.Filter(*r.Filter))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CreateIntegrationActionsRequest) ResourcePath() string {
	return "/v2/integrations/" + r.Id + "/actions"
}

func (r *CreateIntegrationActionsRequest) Method() string {
	return http.MethodPost
}

type UpdateAllIntegrationActionsRequest struct {
	client.BaseRequest
	Id          string
	Create      []IntegrationAction `json:"create"`
	Close       []IntegrationAction `json:"close"`
	Acknowledge []IntegrationAction `json:"acknowledge"`
	AddNote     []IntegrationAction `json:"addNote"`
	Ignore      []IntegrationAction `json:"ignore"`
}

type IntegrationAction struct {
	Type                             ActionType        `json:"type"`
	Name                             string            `json:"name"`
	Alias                            string            `json:"alias"`
	Order                            int               `json:"order,omitempty"`
	User                             string            `json:"user,omitempty"`
	Note                             string            `json:"note,omitempty"`
	Filter                           *Filter           `json:"filter,omitempty"`
	Source                           string            `json:"source,omitempty"`
	Message                          string            `json:"message,omitempty"`
	Description                      string            `json:"description,omitempty"`
	Entity                           string            `json:"entity,omitempty"`
	Priority                         string            `json:"priority,omitempty"`
	CustomPriority                   string            `json:"customPriority,omitempty"`
	AppendAttachments                *bool             `json:"appendAttachments,omitempty"`
	AlertActions                     []string          `json:"alertActions,omitempty"`
	IgnoreAlertActionsFromPayload    *bool             `json:"ignoreAlertActionsFromPayload,omitempty"`
	IgnoreRespondersFromPayload      *bool             `json:"ignoreRespondersFromPayload,omitempty"`
	IgnoreTagsFromPayload            *bool             `json:"ignoreTagsFromPayload,omitempty"`
	IgnoreExtraPropertiesFromPayload *bool             `json:"ignoreExtraPropertiesFromPayload,omitempty"`
	Responders                       []Responder       `json:"responders,omitempty"`
	Tags                             []string          `json:"tags,omitempty"`
	ExtraProperties                  map[string]string `json:"extraProperties,omitempty"`
}

func (r *UpdateAllIntegrationActionsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Integration ID cannot be blank.")
	}
	err := validateActions(r.Create)
	if err != nil {
		return err
	}
	err = validateActions(r.Close)
	if err != nil {
		return err
	}
	err = validateActions(r.AddNote)
	if err != nil {
		return err
	}
	err = validateActions(r.Acknowledge)
	if err != nil {
		return err
	}
	return nil
}

func validateActions(actions []IntegrationAction) error {
	for _, r := range actions {
		err := validateActionType(r.Type)
		if r.Name == "" || r.Type == "" || r.Alias == "" {
			return errors.New("Name, Type and Alias fields cannot be empty.")
		}
		if r.Filter != nil {
			err = validateConditionMatchType(r.Filter.ConditionMatchType)
			if err != nil {
				return err
			}
			err = og.ValidateFilter(og.Filter(*r.Filter))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *UpdateAllIntegrationActionsRequest) ResourcePath() string {
	return "/v2/integrations/" + r.Id + "/actions"
}

func (r *UpdateAllIntegrationActionsRequest) Method() string {
	return http.MethodPut
}

func validateResponders(responders []Responder) error {
	for _, responder := range responders {
		if responder.Type == "" {
			return errors.New("Responder type cannot be empty.")
		}
		if !(responder.Type == User || responder.Type == Team || responder.Type == Schedule || responder.Type == Escalation) {
			return errors.New("Responder type should be one of these: 'User', 'Team', 'Schedule', 'Escalation'")
		}
		if responder.Type == User && responder.Username == "" && responder.Id == "" {
			return errors.New("For responder type user either username or id must be provided.")
		}
		if responder.Type == Team && responder.Name == "" && responder.Id == "" {
			return errors.New("For responder type team either team name or id must be provided.")
		}
		if responder.Type == Schedule && responder.Name == "" && responder.Id == "" {
			return errors.New("For responder type schedule either schedule name or id must be provided.")
		}
		if responder.Type == Escalation && responder.Name == "" && responder.Id == "" {
			return errors.New("For responder type escalation either escalation name or id must be provided.")
		}
	}
	return nil
}

func validateActionType(actionType ActionType) error {
	switch actionType {
	case Create, Close, Acknowledge, AddNote, Ignore:
		return nil
	}
	return errors.New("Action type should be one of these: " +
		"'Create','Close','Acknowledge','AddNote','Ignore'")
}

func validateConditionMatchType(matchType og.ConditionMatchType) error {
	switch matchType {
	case og.MatchAll, og.MatchAllConditions, og.MatchAnyCondition, "":
		return nil
	}
	return errors.New("Action type should be one of these: " +
		"'MatchAll','MatchAllConditions','MatchAnyCondition'")
}
