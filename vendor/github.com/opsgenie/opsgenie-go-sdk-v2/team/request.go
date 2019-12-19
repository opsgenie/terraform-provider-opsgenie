package team

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type Identifier uint32

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type Member struct {
	User User   `json:"user,omitempty"`
	Role string `json:"role,omitempty"`
}

type CreateTeamRequest struct {
	client.BaseRequest
	Description string   `json:"description,omitempty"`
	Name        string   `json:"name,omitempty"`
	Members     []Member `json:"members,omitempty"`
}

func (r *CreateTeamRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name can not be empty")
	}

	return nil
}

func (r *CreateTeamRequest) ResourcePath() string {

	return "/v2/teams"
}

func (r *CreateTeamRequest) Method() string {
	return http.MethodPost
}

type ListTeamRequest struct {
	client.BaseRequest
}

func (r *ListTeamRequest) Validate() error {

	return nil
}

func (r *ListTeamRequest) ResourcePath() string {

	return "/v2/teams"
}

func (r *ListTeamRequest) Method() string {
	return http.MethodGet
}

type DeleteTeamRequest struct {
	client.BaseRequest
	IdentifierType  Identifier
	IdentifierValue string
}

func (r *DeleteTeamRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeleteTeamRequest) ResourcePath() string {

	return "/v2/teams/" + r.IdentifierValue
}

func (r *DeleteTeamRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteTeamRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type GetTeamRequest struct {
	client.BaseRequest
	IdentifierType  Identifier
	IdentifierValue string
}

func (r *GetTeamRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetTeamRequest) ResourcePath() string {

	return "/v2/teams/" + r.IdentifierValue
}

func (r *GetTeamRequest) Method() string {
	return http.MethodGet
}

func (r *GetTeamRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type UpdateTeamRequest struct {
	client.BaseRequest
	Id          string   `json:"id,omitempty"`
	Description string   `json:"description,omitempty"`
	Name        string   `json:"name,omitempty"`
	Members     []Member `json:"members,omitempty"`
}

func (r *UpdateTeamRequest) Validate() error {
	if r.Id == "" {
		return errors.New("team id can not be empty")
	}
	return nil
}

func (r *UpdateTeamRequest) ResourcePath() string {

	return "/v2/teams/" + r.Id
}

func (r *UpdateTeamRequest) Method() string {
	return http.MethodPatch
}

type ListTeamLogsRequest struct {
	client.BaseRequest
	IdentifierType  Identifier
	IdentifierValue string
	Limit           int    `json:"limit,omitempty"`
	Order           string `json:"order,omitempty"`
	Offset          int    `json:"offset,omitempty"`
}

func (r *ListTeamLogsRequest) Validate() error {
	err := validateIdentifier(r.IdentifierValue)
	if err != nil {
		return err
	}

	return nil
}

func (r *ListTeamLogsRequest) ResourcePath() string {

	return "/v2/teams/" + r.IdentifierValue + "/logs"

}

func (r *ListTeamLogsRequest) Method() string {
	return http.MethodGet
}

func (r *ListTeamLogsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	if r.Limit != 0 {
		params["limit"] = strconv.Itoa(r.Limit)
	}
	if r.Offset != 0 {
		params["offset"] = strconv.Itoa(r.Offset)
	}
	if r.Order != "" {
		params["order"] = string(r.Order)
	}

	return params
}

//team role api
type Right struct {
	Right   string `json:"right"`
	Granted *bool  `json:"granted"`
}

type CreateTeamRoleRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
	Name                string  `json:"name"`
	Rights              []Right `json:"rights"`
}

func (r *CreateTeamRoleRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	if r.Name == "" {
		return errors.New("name can not be empty")
	}

	if r.Rights == nil {
		return errors.New("rights can not be empty")
	}

	return nil
}

func (r *CreateTeamRoleRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/roles"

}

func (r *CreateTeamRoleRequest) Method() string {
	return http.MethodPost
}

func (r *CreateTeamRoleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

type GetTeamRoleRequest struct {
	client.BaseRequest
	TeamID   string
	TeamName string
	RoleID   string
	RoleName string
}

func (r *GetTeamRoleRequest) Validate() error {

	if r.TeamID == "" && r.TeamName == "" {
		return errors.New("team identifier can not be empty")
	}

	if r.RoleID == "" && r.RoleName == "" {
		return errors.New("role identifier can not be empty")
	}

	return nil
}

func (r *GetTeamRoleRequest) ResourcePath() string {

	if r.TeamName != "" {
		if r.RoleName != "" {
			return "/v2/teams/" + r.TeamName + "/roles/" + r.RoleName
		}
		return "/v2/teams/" + r.TeamName + "/roles/" + r.RoleID
	}

	if r.RoleName != "" {
		return "/v2/teams/" + r.TeamID + "/roles/" + r.RoleName
	}
	return "/v2/teams/" + r.TeamID + "/roles/" + r.RoleID

}

func (r *GetTeamRoleRequest) Method() string {
	return http.MethodGet
}

func (r *GetTeamRoleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamName != "" {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	if r.RoleName != "" {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type UpdateTeamRoleRequest struct {
	client.BaseRequest
	TeamID   string
	TeamName string
	RoleID   string
	RoleName string
	Name     string  `json:"name"`
	Rights   []Right `json:"rights"`
}

func (r *UpdateTeamRoleRequest) Validate() error {

	if r.TeamID == "" && r.TeamName == "" {
		return errors.New("team identifier can not be empty")
	}

	if r.RoleID == "" && r.RoleName == "" {
		return errors.New("role identifier can not be empty")
	}

	return nil
}

func (r *UpdateTeamRoleRequest) ResourcePath() string {

	if r.TeamName != "" {
		if r.RoleName != "" {
			return "/v2/teams/" + r.TeamName + "/roles/" + r.RoleName
		}
		return "/v2/teams/" + r.TeamName + "/roles/" + r.RoleID
	}

	if r.RoleName != "" {
		return "/v2/teams/" + r.TeamID + "/roles/" + r.RoleName
	}
	return "/v2/teams/" + r.TeamID + "/roles/" + r.RoleID

}

func (r *UpdateTeamRoleRequest) Method() string {
	return http.MethodPatch
}

func (r *UpdateTeamRoleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamName != "" {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	if r.RoleName != "" {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type DeleteTeamRoleRequest struct {
	client.BaseRequest
	TeamID   string
	TeamName string
	RoleID   string
	RoleName string
}

func (r *DeleteTeamRoleRequest) Validate() error {
	if r.TeamID == "" && r.TeamName == "" {
		return errors.New("team identifier can not be empty")
	}

	if r.RoleID == "" && r.RoleName == "" {
		return errors.New("role identifier can not be empty")
	}

	return nil
}

func (r *DeleteTeamRoleRequest) ResourcePath() string {
	if r.TeamName != "" {
		if r.RoleName != "" {
			return "/v2/teams/" + r.TeamName + "/roles/" + r.RoleName
		}
		return "/v2/teams/" + r.TeamName + "/roles/" + r.RoleID
	}

	if r.RoleName != "" {
		return "/v2/teams/" + r.TeamID + "/roles/" + r.RoleName
	}
	return "/v2/teams/" + r.TeamID + "/roles/" + r.RoleID

}

func (r *DeleteTeamRoleRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteTeamRoleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamName != "" {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	if r.RoleName != "" {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type ListTeamRoleRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
}

func (r *ListTeamRoleRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *ListTeamRoleRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/roles"
}

func (r *ListTeamRoleRequest) Method() string {
	return http.MethodGet
}

func (r *ListTeamRoleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

const (
	Name Identifier = iota
	Id
	Username
)

func validateIdentifier(identifier string) error {
	if identifier == "" {
		return errors.New("team identifier can not be empty")
	}
	return nil
}

//team member api
type AddTeamMemberRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
	User                User   `json:"user,omitempty"`
	Role                string `json:"role,omitempty"`
}

func (r *AddTeamMemberRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	if r.User.ID == "" && r.User.Username == "" {
		return errors.New("user can not be empty")
	}

	return nil
}

func (r *AddTeamMemberRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/members"

}

func (r *AddTeamMemberRequest) Method() string {
	return http.MethodPost
}

func (r *AddTeamMemberRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

type RemoveTeamMemberRequest struct {
	client.BaseRequest
	TeamIdentifierType    Identifier
	TeamIdentifierValue   string
	MemberIdentifierType  Identifier
	MemberIdentifierValue string
}

func (r *RemoveTeamMemberRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	if r.MemberIdentifierValue == "" {
		return errors.New("member identifier cannot be empty")
	}

	if r.MemberIdentifierType != Username && r.MemberIdentifierType != Id {
		return errors.New("member identifier must be id or username")
	}

	return nil
}

func (r *RemoveTeamMemberRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/members/" + r.MemberIdentifierValue

}

func (r *RemoveTeamMemberRequest) Method() string {
	return http.MethodDelete
}

func (r *RemoveTeamMemberRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

//team routing rule api
type NotifyType string

const (
	EscalationNotifyType NotifyType = "escalation"
	ScheduleNotifyType   NotifyType = "schedule"
	None                 NotifyType = "none"
)

type Notify struct {
	Type NotifyType `json:"type, omitempty"`
	Name string     `json:"name,omitempty"`
	Id   string     `json:"id,omitempty"`
}

type CreateRoutingRuleRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
	Name                string              `json:"name,omitempty"`
	Order               *int                `json:"order,omitempty"`
	Timezone            string              `json:"timezone,omitempty"`
	Criteria            *og.Criteria        `json:"criteria,omitempty"`
	TimeRestriction     *og.TimeRestriction `json:"timeRestriction,omitempty"`
	Notify              *Notify             `json:"notify"`
}

func (r *CreateRoutingRuleRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}
	if r.Notify == nil {
		return errors.New("notify can not be empty")
	} else if r.Notify != nil {
		err := validateNotifyType(r.Notify.Type)
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

	if r.Criteria != nil {
		err = og.ValidateCriteria(*r.Criteria)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *CreateRoutingRuleRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/routing-rules"

}

func (r *CreateRoutingRuleRequest) Method() string {
	return http.MethodPost
}

func (r *CreateRoutingRuleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

func (r *CreateRoutingRuleRequest) WithTimeRestriction(timeRestriction og.TimeRestriction) CreateRoutingRuleRequest {
	r.TimeRestriction = &timeRestriction
	return *r
}

type GetRoutingRuleRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
	RoutingRuleId       string
}

func (r *GetRoutingRuleRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	if r.RoutingRuleId == "" {
		return errors.New("routing rule id can not be empty")
	}

	return nil
}

func (r *GetRoutingRuleRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/routing-rules/" + r.RoutingRuleId

}

func (r *GetRoutingRuleRequest) Method() string {
	return http.MethodGet
}

func (r *GetRoutingRuleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

type UpdateRoutingRuleRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
	RoutingRuleId       string
	Name                string              `json:"name,omitempty"`
	Timezone            string              `json:"timezone,omitempty"`
	Criteria            *og.Criteria        `json:"criteria,omitempty"`
	TimeRestriction     *og.TimeRestriction `json:"timeRestriction,omitempty"`
	Notify              *Notify             `json:"notify,omitempty"`
}

func (r *UpdateRoutingRuleRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	if r.RoutingRuleId == "" {
		return errors.New("routing rule id can not be empty")
	}

	if r.TimeRestriction != nil {
		err := og.ValidateRestrictions(r.TimeRestriction)
		if err != nil {
			return err
		}
	}

	if r.Criteria != nil {
		err = og.ValidateCriteria(*r.Criteria)
		if err != nil {
			return err
		}
	}

	if r.Notify != nil {
		err := validateNotifyType(r.Notify.Type)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UpdateRoutingRuleRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/routing-rules/" + r.RoutingRuleId

}

func (r *UpdateRoutingRuleRequest) Method() string {
	return http.MethodPatch
}

func (r *UpdateRoutingRuleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

type DeleteRoutingRuleRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
	RoutingRuleId       string
}

func (r *DeleteRoutingRuleRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	if r.RoutingRuleId == "" {
		return errors.New("routing rule id can not be empty")
	}

	return nil
}

func (r *DeleteRoutingRuleRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/routing-rules/" + r.RoutingRuleId

}

func (r *DeleteRoutingRuleRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteRoutingRuleRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

type ListRoutingRulesRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
}

func (r *ListRoutingRulesRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	return nil
}

func (r *ListRoutingRulesRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/routing-rules"

}

func (r *ListRoutingRulesRequest) Method() string {
	return http.MethodGet
}

func (r *ListRoutingRulesRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

type ChangeRoutingRuleOrderRequest struct {
	client.BaseRequest
	TeamIdentifierType  Identifier
	TeamIdentifierValue string
	RoutingRuleId       string
	Order               *int `json:"order"`
}

func (r *ChangeRoutingRuleOrderRequest) Validate() error {
	err := validateIdentifier(r.TeamIdentifierValue)
	if err != nil {
		return err
	}

	if r.RoutingRuleId == "" {
		return errors.New("routing rule id can not be empty")
	}

	if r.Order == nil {
		return errors.New("order can not be empty")
	}

	return nil
}

func (r *ChangeRoutingRuleOrderRequest) ResourcePath() string {

	return "/v2/teams/" + r.TeamIdentifierValue + "/routing-rules/" + r.RoutingRuleId + "/change-order"

}

func (r *ChangeRoutingRuleOrderRequest) Method() string {
	return http.MethodPost
}

func (r *ChangeRoutingRuleOrderRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.TeamIdentifierType == Name {
		params["teamIdentifierType"] = "name"
	} else {
		params["teamIdentifierType"] = "id"
	}

	return params
}

func validateNotifyType(notifyType NotifyType) error {
	switch notifyType {
	case EscalationNotifyType, ScheduleNotifyType, None:
		return nil
	}
	return errors.New("Notify type should be one of these: " +
		"'EscalationNotifyType','ScheduleNotifyType','None'")
}
