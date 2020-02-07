package team

import (
	"errors"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRequest_Validate(t *testing.T) {
	createRequest := &CreateTeamRequest{}
	err := createRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("name can not be empty").Error())
}

func TestGetRequest_Validate(t *testing.T) {
	getRequest := &GetTeamRequest{}
	err := getRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())
}

func TestDeleteRequest_Validate(t *testing.T) {
	deleteRequest := &DeleteTeamRequest{}
	err := deleteRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())
}

func TestUpdateRequest_Validate(t *testing.T) {
	updateRequest := &UpdateTeamRequest{}
	err := updateRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("team id can not be empty").Error())
}

func TestCreateTeamRoleRequest_Validate(t *testing.T) {
	createTeamRoleRequest := &CreateTeamRoleRequest{TeamIdentifierValue: "xx", TeamIdentifierType: Name}
	err := createTeamRoleRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("name can not be empty").Error())
}

func TestGetTeamRoleRequest_Validate(t *testing.T) {
	getTeamRoleRequest := &GetTeamRoleRequest{}
	err := getTeamRoleRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())
}

func TestAddTeamMemberRequest_Validate(t *testing.T) {
	addTeamMemberRequest := &AddTeamMemberRequest{}
	err := addTeamMemberRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	addTeamMemberRequest.TeamIdentifierValue = "test"
	err = addTeamMemberRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("user can not be empty").Error())

	addTeamMemberRequest.User = User{Username: "test0@gmail.com"}

	err = addTeamMemberRequest.Validate()
	assert.Nil(t, err)

}

func TestRemoveTeamMemberRequest_Validate(t *testing.T) {
	removeTeamMemberRequest := &RemoveTeamMemberRequest{}
	err := removeTeamMemberRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	removeTeamMemberRequest.TeamIdentifierValue = "test"
	err = removeTeamMemberRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("member identifier cannot be empty").Error())

	removeTeamMemberRequest.MemberIdentifierType = Name
	removeTeamMemberRequest.MemberIdentifierValue = "test2"
	err = removeTeamMemberRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("member identifier must be id or username").Error())

	removeTeamMemberRequest.MemberIdentifierType = Username
	err = removeTeamMemberRequest.Validate()
	assert.Nil(t, err)

}

func TestCreateRoutingRuleRequest_Validate(t *testing.T) {
	createRoutingRuleRequest := &CreateRoutingRuleRequest{}
	err := createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	createRoutingRuleRequest.TeamIdentifierValue = "test"
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("notify can not be empty").Error())

	createRoutingRuleRequest.Notify = &Notify{Type: None}
	err = createRoutingRuleRequest.Validate()
	assert.Nil(t, err)

	createRoutingRuleRequest.Criteria = &og.Filter{
		ConditionMatchType: "invalid type",
		Conditions:         nil,
	}
	err = createRoutingRuleRequest.Validate()
	assert.Contains(t, err.Error(), errors.New("filter condition type should be one of").Error())

	createRoutingRuleRequest.Criteria.ConditionMatchType = og.MatchAllConditions
	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("filter conditions cannot be empty").Error())

	isNot := false
	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         "random field",
			IsNot:         &isNot,
			Operation:     "",
			ExpectedValue: "",
			Key:           "",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("condition field should be one of message, alias, description, source, entity, tags, actions, details, extra-properties, recipients, teams or priority").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Message,
			IsNot:         &isNot,
			Operation:     og.Contains,
			ExpectedValue: "",
			Key:           "asd",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("condition key is only valid for extra-properties field").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Actions,
			IsNot:         &isNot,
			Operation:     og.LessThan,
			ExpectedValue: "",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("less-than is not valid operation for actions").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Actions,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("greater-than is not valid operation for actions").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Message,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("greater-than is not valid operation for message").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Details,
			IsNot:         &isNot,
			Operation:     og.EqualsIgnoreWhitespcae,
			ExpectedValue: "",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("equals-ignore-whitespace is not valid operation for details").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Priority,
			IsNot:         &isNot,
			Operation:     og.Matches,
			ExpectedValue: "",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("matches is not valid operation for priority").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Priority,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "critical",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("for field priority expected value should be one of P1, P2, P3, P4, P5").Error())

	createRoutingRuleRequest.Criteria.Conditions = []og.Condition{
		{
			Field:         og.Priority,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "P4",
		},
	}
	err = createRoutingRuleRequest.Validate()
	assert.Nil(t, err)

}

func TestGetRoutingRuleRequest_Validate(t *testing.T) {
	getRoutingRuleRequest := &GetRoutingRuleRequest{}
	err := getRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	getRoutingRuleRequest.TeamIdentifierValue = "test"
	err = getRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("routing rule id can not be empty").Error())

	getRoutingRuleRequest.RoutingRuleId = "test"
	err = getRoutingRuleRequest.Validate()
	assert.Nil(t, err)

}

func TestUpdateRoutingRuleRequest_Validate(t *testing.T) {
	updateRoutingRuleRequest := &UpdateRoutingRuleRequest{}
	err := updateRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	updateRoutingRuleRequest.TeamIdentifierValue = "test"
	err = updateRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("routing rule id can not be empty").Error())

	updateRoutingRuleRequest.RoutingRuleId = "test"
	err = updateRoutingRuleRequest.Validate()
	assert.Nil(t, err)

}

func TestDeleteRoutingRuleRequest_Validate(t *testing.T) {
	deleteRoutingRuleRequest := &DeleteRoutingRuleRequest{}
	err := deleteRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	deleteRoutingRuleRequest.TeamIdentifierValue = "test"
	err = deleteRoutingRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("routing rule id can not be empty").Error())

	deleteRoutingRuleRequest.RoutingRuleId = "test"
	err = deleteRoutingRuleRequest.Validate()
	assert.Nil(t, err)

}

func TestListRoutingRulesRequest_Validate(t *testing.T) {
	listRoutingRulesRequest := &ListRoutingRulesRequest{}
	err := listRoutingRulesRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	listRoutingRulesRequest.TeamIdentifierValue = "test"
	err = listRoutingRulesRequest.Validate()
	assert.Nil(t, err)

}

func TestChangeRoutingRuleOrderRequest_Validate(t *testing.T) {
	changeRoutingRuleOrderRequest := &ChangeRoutingRuleOrderRequest{}
	err := changeRoutingRuleOrderRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("team identifier can not be empty").Error())

	changeRoutingRuleOrderRequest.TeamIdentifierValue = "test"
	err = changeRoutingRuleOrderRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("routing rule id can not be empty").Error())

	changeRoutingRuleOrderRequest.RoutingRuleId = "test"
	err = changeRoutingRuleOrderRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("order can not be empty").Error())

	order := new(int)
	*order = 0
	changeRoutingRuleOrderRequest.Order = order
	err = changeRoutingRuleOrderRequest.Validate()
	assert.Nil(t, err)

}

func TestNotifyType_Validate(t *testing.T) {
	err := validateNotifyType("test")
	assert.Equal(t, err.Error(), errors.New("Notify type should be one of these: "+
		"'EscalationNotifyType','ScheduleNotifyType','None'").Error())

	err = validateNotifyType(ScheduleNotifyType)
	assert.Nil(t, err)

	err = validateNotifyType(EscalationNotifyType)
	assert.Nil(t, err)

	err = validateNotifyType(None)
	assert.Nil(t, err)

}
