package notification

import (
	"net/http"
	"testing"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequest_Validate(t *testing.T) {
	var err error
	createRequest := &CreateRuleStepRequest{}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	createRequest.UserIdentifier = "123"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	createRequest.RuleId = "123"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("To cannot be empty.").Error())

	createRequest.Contact.To = "ab@email.com"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Method cannot be empty.").Error())

	createRequest.Contact.MethodOfContact = og.Email
	err = createRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, createRequest.ResourcePath(), "/v2/users/123/notification-rules/123/steps")
	assert.Equal(t, createRequest.Method(), http.MethodPost)

}

func TestGetRequest_Validate(t *testing.T) {
	var err error
	getRequest := &GetRuleStepRequest{}
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	getRequest.UserIdentifier = "123"
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	getRequest.RuleId = "123"
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule Step identifier cannot be empty.").Error())

	getRequest.RuleStepId = "1234"
	err = getRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, getRequest.ResourcePath(), "/v2/users/123/notification-rules/123/steps/1234")
	assert.Equal(t, getRequest.Method(), http.MethodGet)

}

func TestUpdateRequest_Validate(t *testing.T) {
	var err error
	updateRequest := &UpdateRuleStepRequest{}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	updateRequest.UserIdentifier = "123"
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	updateRequest.RuleId = "123"
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule Step identifier cannot be empty.").Error())

	updateRequest.RuleStepId = "1234"
	err = updateRequest.Validate()
	assert.Nil(t, err)

	updateRequest.Contact = &og.Contact{}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("To cannot be empty.").Error())

	updateRequest.Contact = &og.Contact{To: "abc"}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Method cannot be empty.").Error())

	updateRequest.Contact = &og.Contact{To: "abc", MethodOfContact: "email"}
	err = updateRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, updateRequest.ResourcePath(), "/v2/users/123/notification-rules/123/steps/1234")
	assert.Equal(t, updateRequest.Method(), http.MethodPatch)

}

func TestDeleteRequest_Validate(t *testing.T) {
	var err error
	deleteRequest := &DeleteRuleStepRequest{}
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	deleteRequest.UserIdentifier = "123"
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	deleteRequest.RuleId = "123"
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule Step identifier cannot be empty.").Error())

	deleteRequest.RuleStepId = "1234"
	err = deleteRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, deleteRequest.ResourcePath(), "/v2/users/123/notification-rules/123/steps/1234")
	assert.Equal(t, deleteRequest.Method(), http.MethodDelete)

}

func TestListRequest_Validate(t *testing.T) {
	var err error
	listRequest := &ListRuleStepsRequest{}
	err = listRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	listRequest.UserIdentifier = "123"
	err = listRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	listRequest.RuleId = "1234"
	err = listRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, listRequest.ResourcePath(), "/v2/users/123/notification-rules/1234/steps")
	assert.Equal(t, listRequest.Method(), http.MethodGet)

}

func TestEnableRequest_Validate(t *testing.T) {
	var err error
	enableRequest := &EnableRuleStepRequest{}
	err = enableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	enableRequest.UserIdentifier = "123"
	err = enableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	enableRequest.RuleId = "123"
	err = enableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule Step identifier cannot be empty.").Error())

	enableRequest.RuleStepId = "1234"
	err = enableRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, enableRequest.ResourcePath(), "/v2/users/123/notification-rules/123/steps/1234/enable")
	assert.Equal(t, enableRequest.Method(), http.MethodPost)

}

func TestDisableRequest_Validate(t *testing.T) {
	var err error
	disableRuleStepRequest := &DisableRuleStepRequest{}
	err = disableRuleStepRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	disableRuleStepRequest.UserIdentifier = "123"
	err = disableRuleStepRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	disableRuleStepRequest.RuleId = "123"
	err = disableRuleStepRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule Step identifier cannot be empty.").Error())

	disableRuleStepRequest.RuleStepId = "1234"
	err = disableRuleStepRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, disableRuleStepRequest.ResourcePath(), "/v2/users/123/notification-rules/123/steps/1234/disable")
	assert.Equal(t, disableRuleStepRequest.Method(), http.MethodPost)
}

func TestCreateRuleRequest_Validate(t *testing.T) {
	var err error
	createRequest := &CreateRuleRequest{}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	createRequest.UserIdentifier = "123"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Name cannot be empty.").Error())

	createRequest.Name = "test"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Action type cannot be empty.").Error())

	createRequest.ActionType = ScheduleEnd
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Notification time cannot be empty.").Error())

	notificationTypes := make([]NotificationTimeType, 1)
	notificationTypes[0] = JustBefore
	createRequest.NotificationTime = notificationTypes
	err = createRequest.Validate()
	assert.Nil(t, err)

	schedules := make([]Schedule, 1)
	schedules[0] = Schedule{TypeOfSchedule: ""}
	createRequest.Schedules = schedules
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Type of schedule must be schedule.").Error())

	schedules[0] = Schedule{TypeOfSchedule: "schedule"}
	createRequest.Schedules = schedules
	err = createRequest.Validate()
	assert.Nil(t, err)

	createRequest.ActionType = CreateAlert
	steps := make([]*og.Step, 1)
	contact := og.Contact{}
	steps[0] = &og.Step{Contact: contact}
	createRequest.Steps = steps
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("To cannot be empty.").Error())

	contact = og.Contact{To: "abc@a.com"}
	steps[0] = &og.Step{Contact: contact}
	createRequest.Steps = steps
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Method cannot be empty.").Error())

	contact = og.Contact{To: "abc@a.com", MethodOfContact: og.Email}
	steps[0] = &og.Step{Contact: contact}
	createRequest.Steps = steps
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("SendAfter cannot be empty.").Error())

	sendAfter := &og.SendAfter{TimeAmount: 3, TimeUnit: "minutes"}
	steps[0] = &og.Step{Contact: contact, SendAfter: sendAfter}
	createRequest.Steps = steps
	err = createRequest.Validate()
	assert.Nil(t, err)

	createRequest.Criteria = &og.Filter{}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("filter condition type should be one of match-all, match-any-condition or match-all-conditions").Error())

	createRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("filter conditions cannot be empty").Error())

	conditions := make([]og.Condition, 1)
	conditions[0] = og.Condition{Operation: og.Contains}
	createRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition, Conditions: conditions}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("condition field should be one of message, alias, description, source, entity, tags, actions, details, extra-properties, recipients, teams or priority").Error())

	conditions[0] = og.Condition{Field: og.Alias, Operation: "a"}
	createRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition, Conditions: conditions}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("a is not valid operation for alias").Error())

	conditions[0] = og.Condition{Field: og.Alias, Operation: og.Contains}
	createRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition, Conditions: conditions}
	err = createRequest.Validate()
	assert.Nil(t, err)

	createRequest.TimeRestriction = &og.TimeRestriction{}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Type of time restriction must be time-of-day or weekday-and-time-of-day.").Error())

	createRequest.TimeRestriction = &og.TimeRestriction{Type: og.TimeOfDay}
	restrictions := make([]og.Restriction, 1)
	restrictions[0] = og.Restriction{}
	createRequest.TimeRestriction = &og.TimeRestriction{Type: og.TimeOfDay, Restriction: og.Restriction{}}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("startHour, startMin, endHour, endMin cannot be empty.").Error())

	restrictions[0] = og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(18)}
	createRequest.TimeRestriction = &og.TimeRestriction{Type: og.TimeOfDay, Restriction: og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(18)}}
	err = createRequest.Validate()
	assert.Nil(t, err)

	restrictions[0] = og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(18)}
	createRequest.TimeRestriction = &og.TimeRestriction{Type: og.WeekdayAndTimeOfDay, RestrictionList: restrictions}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("startDay, endDay cannot be empty.").Error())

	restrictions[0] = og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(18), StartDay: og.Wednesday, EndDay: og.Friday}
	createRequest.TimeRestriction = &og.TimeRestriction{Type: og.WeekdayAndTimeOfDay, RestrictionList: restrictions}
	err = createRequest.Validate()
	assert.Nil(t, err)

	createRequest.Repeat = &Repeat{}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Loop after must have a positive integer value.").Error())

	createRequest.Repeat = &Repeat{LoopAfter: 2}
	err = createRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, createRequest.ResourcePath(), "/v2/users/123/notification-rules")
	assert.Equal(t, createRequest.Method(), http.MethodPost)

}

func TestGetRuleRequest_Validate(t *testing.T) {
	var err error
	getRequest := &GetRuleRequest{}
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	getRequest.UserIdentifier = "123"
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	getRequest.RuleId = "123"
	err = getRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, getRequest.ResourcePath(), "/v2/users/123/notification-rules/123")
	assert.Equal(t, getRequest.Method(), http.MethodGet)

}

func TestUpdateRuleRequest_Validate(t *testing.T) {
	var err error
	updateRuleRequest := &UpdateRuleRequest{}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	updateRuleRequest.UserIdentifier = "123"
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	updateRuleRequest.RuleId = "123"
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	notificationTypes := make([]NotificationTimeType, 1)
	notificationTypes[0] = JustBefore
	updateRuleRequest.NotificationTime = notificationTypes
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	schedules := make([]Schedule, 1)
	schedules[0] = Schedule{TypeOfSchedule: ""}
	updateRuleRequest.Schedules = schedules
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Type of schedule must be schedule.").Error())

	schedules[0] = Schedule{TypeOfSchedule: "schedule"}
	updateRuleRequest.Schedules = schedules
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	steps := make([]*og.Step, 1)
	contact := og.Contact{}
	steps[0] = &og.Step{Contact: contact}
	updateRuleRequest.Steps = steps
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("To cannot be empty.").Error())

	contact = og.Contact{To: "abc@a.com"}
	steps[0] = &og.Step{Contact: contact}
	updateRuleRequest.Steps = steps
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Method cannot be empty.").Error())

	contact = og.Contact{To: "abc@a.com", MethodOfContact: og.Email}
	steps[0] = &og.Step{Contact: contact}
	updateRuleRequest.Steps = steps
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	updateRuleRequest.Criteria = &og.Filter{}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("filter condition type should be one of match-all, match-any-condition or match-all-conditions").Error())

	updateRuleRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("filter conditions cannot be empty").Error())

	conditions := make([]og.Condition, 1)
	conditions[0] = og.Condition{Operation: og.Contains}
	updateRuleRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition, Conditions: conditions}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("condition field should be one of message, alias, description, source, entity, tags, actions, details, extra-properties, recipients, teams or priority").Error())

	conditions[0] = og.Condition{Field: og.Alias, Operation: "a"}
	updateRuleRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition, Conditions: conditions}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("a is not valid operation for alias").Error())

	conditions[0] = og.Condition{Field: og.Alias, Operation: og.Contains}
	updateRuleRequest.Criteria = &og.Filter{ConditionMatchType: og.MatchAnyCondition, Conditions: conditions}
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	updateRuleRequest.TimeRestriction = &og.TimeRestriction{}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Type of time restriction must be time-of-day or weekday-and-time-of-day.").Error())

	updateRuleRequest.TimeRestriction = &og.TimeRestriction{Type: og.TimeOfDay}
	restrictions := make([]og.Restriction, 1)
	restrictions[0] = og.Restriction{}
	updateRuleRequest.TimeRestriction = &og.TimeRestriction{Type: og.TimeOfDay, Restriction: og.Restriction{}}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("startHour, startMin, endHour, endMin cannot be empty.").Error())

	restrictions[0] = og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(18)}
	updateRuleRequest.TimeRestriction = &og.TimeRestriction{Type: og.TimeOfDay, Restriction: og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(18)}}
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	restrictions[0] = og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(12)}
	updateRuleRequest.TimeRestriction = &og.TimeRestriction{Type: og.WeekdayAndTimeOfDay, RestrictionList: restrictions}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("startDay, endDay cannot be empty.").Error())

	restrictions[0] = og.Restriction{StartMin: og.Minute(12), StartHour: og.Hour(1), EndHour: og.Hour(1), EndMin: og.Minute(12), StartDay: og.Wednesday, EndDay: og.Friday}
	updateRuleRequest.TimeRestriction = &og.TimeRestriction{Type: og.WeekdayAndTimeOfDay, RestrictionList: restrictions}
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	updateRuleRequest.Repeat = &Repeat{}
	err = updateRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Loop after must have a positive integer value.").Error())

	updateRuleRequest.Repeat = &Repeat{LoopAfter: 2}
	err = updateRuleRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, updateRuleRequest.ResourcePath(), "/v2/users/123/notification-rules/123")
	assert.Equal(t, updateRuleRequest.Method(), http.MethodPatch)
}

func TestDeleteRuleRequest_Validate(t *testing.T) {
	var err error
	deleteRequest := &DeleteRuleRequest{}
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	deleteRequest.UserIdentifier = "123"
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	deleteRequest.RuleId = "123"
	err = deleteRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, deleteRequest.ResourcePath(), "/v2/users/123/notification-rules/123")
	assert.Equal(t, deleteRequest.Method(), http.MethodDelete)

}

func TestListRuleRequest_Validate(t *testing.T) {
	var err error
	listRequest := &ListRuleRequest{}
	err = listRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	listRequest.UserIdentifier = "123"
	err = listRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, listRequest.ResourcePath(), "/v2/users/123/notification-rules")
	assert.Equal(t, listRequest.Method(), http.MethodGet)

}

func TestEnableRuleRequest_Validate(t *testing.T) {
	var err error
	enableRequest := &EnableRuleRequest{}
	err = enableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	enableRequest.UserIdentifier = "123"
	err = enableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	enableRequest.RuleId = "123"
	err = enableRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, enableRequest.ResourcePath(), "/v2/users/123/notification-rules/123/enable")
	assert.Equal(t, enableRequest.Method(), http.MethodPost)

}

func TestDisableRuleRequest_Validate(t *testing.T) {
	var err error
	disableRuleRequest := &DisableRuleRequest{}
	err = disableRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	disableRuleRequest.UserIdentifier = "123"
	err = disableRuleRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule identifier cannot be empty.").Error())

	disableRuleRequest.RuleId = "123"
	err = disableRuleRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, disableRuleRequest.ResourcePath(), "/v2/users/123/notification-rules/123/disable")
	assert.Equal(t, disableRuleRequest.Method(), http.MethodPost)
}

func TestCopyNotificationRulesRequest_Validate(t *testing.T) {
	var err error
	copyNotificationRulesRequest := &CopyNotificationRulesRequest{}
	err = copyNotificationRulesRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	copyNotificationRulesRequest.UserIdentifier = "123"
	err = copyNotificationRulesRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("You must specify a list of the users which you want to copy the rules to.").Error())

	users := make([]string, 1)
	users[0] = "user1@og.com"
	copyNotificationRulesRequest.ToUsers = users
	err = copyNotificationRulesRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Specify a list of the action types you want to copy the rules of.").Error())

	rules := make([]RuleTypes, 1)
	rules[0] = All
	copyNotificationRulesRequest.RuleTypes = rules
	err = copyNotificationRulesRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, copyNotificationRulesRequest.ResourcePath(), "/v2/users/123/notification-rules/copy-to")
	assert.Equal(t, copyNotificationRulesRequest.Method(), http.MethodPost)
}
