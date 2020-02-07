package policy

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAlertPolicy_Validate(t *testing.T) {
	req := &CreateAlertPolicyRequest{}
	req.MainFields = MainFields{PolicyType: "alert"}
	err := req.Validate()
	assert.Equal(t, err.Error(), errors.New("policy name cannot be empty").Error())

	req.Name = "a policy"
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("alert message cannot be empty").Error())

	req.Message = "message"
	err = req.Validate()
	assert.Nil(t, err)

	req.Filter = &og.Filter{
		ConditionMatchType: "invalid type",
		Conditions:         nil,
	}
	err = req.Validate()
	assert.Contains(t, err.Error(), errors.New("filter condition type should be one of").Error())

	req.Filter.ConditionMatchType = og.MatchAllConditions
	req.Filter.Conditions = []og.Condition{}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("filter conditions cannot be empty").Error())

	isNot := false
	req.Filter.Conditions = []og.Condition{
		{
			Field:         "random field",
			IsNot:         &isNot,
			Operation:     "",
			ExpectedValue: "",
			Key:           "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("condition field should be one of message, alias, description, source, entity, tags, actions, details, extra-properties, recipients, teams or priority").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Message,
			IsNot:         &isNot,
			Operation:     og.Contains,
			ExpectedValue: "",
			Key:           "asd",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("condition key is only valid for extra-properties field").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Actions,
			IsNot:         &isNot,
			Operation:     og.LessThan,
			ExpectedValue: "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("less-than is not valid operation for actions").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Actions,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("greater-than is not valid operation for actions").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Message,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("greater-than is not valid operation for message").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Details,
			IsNot:         &isNot,
			Operation:     og.EqualsIgnoreWhitespcae,
			ExpectedValue: "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("equals-ignore-whitespace is not valid operation for details").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Priority,
			IsNot:         &isNot,
			Operation:     og.Matches,
			ExpectedValue: "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("matches is not valid operation for priority").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Priority,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "critical",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("for field priority expected value should be one of P1, P2, P3, P4, P5").Error())

	req.Filter.Conditions = []og.Condition{
		{
			Field:         og.Priority,
			IsNot:         &isNot,
			Operation:     og.GreaterThan,
			ExpectedValue: "P4",
		},
	}
	err = req.Validate()
	assert.Nil(t, err)

	req.Responders = &[]alert.Responder{
		{
			Type:     alert.ScheduleResponder,
			Name:     "",
			Id:       "",
			Username: "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("responder type for alert policy should be one of team or user").Error())

	req.Responders = &[]alert.Responder{
		{
			Type:     alert.UserResponder,
			Name:     "",
			Id:       "",
			Username: "",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("responder id should be provided").Error())

	req.Responders = &[]alert.Responder{
		{
			Type:     alert.UserResponder,
			Name:     "",
			Id:       "userid",
			Username: "",
		},
	}
	err = req.Validate()
	assert.Nil(t, err)

	req.Responders = &[]alert.Responder{
		{
			Type:     alert.TeamResponder,
			Name:     "",
			Id:       "",
			Username: "user1",
		},
	}
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("responder id should be provided").Error())

	req.Responders = &[]alert.Responder{
		{
			Type:     alert.TeamResponder,
			Name:     "",
			Id:       "teamId",
			Username: "",
		},
	}
	err = req.Validate()
	assert.Nil(t, err)

	req.Priority = "asd"
	err = req.Validate()
	assert.Equal(t, err.Error(), errors.New("Priority should be one of these: 'P1', 'P2', 'P3', 'P4' and 'P5'").Error())

	req.Priority = alert.P1
	err = req.Validate()
	assert.Nil(t, err)
}

func TestCreateNotificationPolicy_Validate(t *testing.T) {
	request := CreateNotificationPolicyRequest{}
	request.MainFields = MainFields{
		PolicyType: "notification",
		Name:       "name",
	}
	err := request.Validate()
	assert.Equal(t, err.Error(), "policy team id should be provided")

	request.TeamId = "teamId"
	request.AutoRestartAction = &AutoRestartAction{}
	err = request.Validate()
	assert.Equal(t, err.Error(), "autoRestart action duration cannot be empty")

	request.AutoRestartAction.Duration = &Duration{}
	err = request.Validate()
	assert.Equal(t, err.Error(), "duration timeAmount should be greater than zero")

	request.AutoRestartAction.Duration = &Duration{TimeUnit: "asd"}
	err = request.Validate()
	assert.Equal(t, err.Error(), "timeUnit provided for duration is not valid")

	request.AutoRestartAction.Duration = &Duration{TimeAmount: 5}
	err = request.Validate()
	assert.Equal(t, og.Minutes, request.AutoRestartAction.Duration.TimeUnit)

	request.AutoRestartAction.MaxRepeatCount = -1
	err = request.Validate()
	assert.Equal(t, err.Error(), "autoRestart maxRepeatCount is not valid")

	request.AutoRestartAction.MaxRepeatCount = 5
	err = request.Validate()
	assert.Nil(t, err)

	request.AutoCloseAction = &AutoCloseAction{}
	err = request.Validate()
	assert.Equal(t, "autoClose action duration cannot be empty", err.Error())

	request.AutoCloseAction = &AutoCloseAction{Duration: &Duration{}}
	err = request.Validate()
	assert.Equal(t, "duration timeAmount should be greater than zero", err.Error())

	request.AutoCloseAction = &AutoCloseAction{Duration: &Duration{TimeAmount: 1}}
	err = request.Validate()
	assert.Nil(t, err)

	request.DeDuplicationAction = &DeDuplicationAction{}
	err = request.Validate()
	assert.Equal(t, "deDuplication action type should be one of value-based or frequency-based", err.Error())

	request.DeDuplicationAction = &DeDuplicationAction{DeDuplicationActionType: FrequencyBased, Count: -4}
	err = request.Validate()
	assert.Equal(t, "deDuplication count is not valid", err.Error())

	request.DeDuplicationAction = &DeDuplicationAction{DeDuplicationActionType: FrequencyBased, Count: 4}
	err = request.Validate()
	assert.Nil(t, err)

	request.DeDuplicationAction = &DeDuplicationAction{DeDuplicationActionType: FrequencyBased, Count: 4, Duration: &Duration{}}
	err = request.Validate()
	assert.Equal(t, "duration timeAmount should be greater than zero", err.Error())

	request.DeDuplicationAction = &DeDuplicationAction{DeDuplicationActionType: FrequencyBased, Count: 4, Duration: &Duration{TimeAmount: 11}}
	err = request.Validate()
	assert.Nil(t, err)

	request.DelayAction = &DelayAction{}
	err = request.Validate()
	assert.Equal(t, "delay option should be one of for-duration, next-time, next-weekday, next-monday, next-tuesday, next-wednesday, next-thursday, next-friday, next-saturday, next-sunday", err.Error())

	request.DelayAction = &DelayAction{DelayOption: ForDuration}
	err = request.Validate()
	assert.Equal(t, "delayAction duration cannot be empty", err.Error())

	request.DelayAction = &DelayAction{DelayOption: ForDuration, Duration: &Duration{}}
	err = request.Validate()
	assert.Equal(t, "duration timeAmount should be greater than zero", err.Error())

	request.DelayAction = &DelayAction{DelayOption: ForDuration, Duration: &Duration{TimeAmount: 5}}
	err = request.Validate()
	assert.Nil(t, err)

	request.DelayAction = &DelayAction{DelayOption: NextSunday, UntilHour: -3}
	err = request.Validate()
	assert.Equal(t, "delayAction's UntilHour or UntilMinute is not valid", err.Error())

	request.DelayAction = &DelayAction{DelayOption: NextSunday, UntilHour: 5, UntilMinute: 60}
	err = request.Validate()
	assert.Equal(t, "delayAction's UntilHour or UntilMinute is not valid", err.Error())

	request.DelayAction = &DelayAction{DelayOption: NextSunday, UntilHour: 5, UntilMinute: 55}
	err = request.Validate()
	assert.Nil(t, err)

}

func TestGetAlertPolicy_Validate(t *testing.T) {
	request := GetAlertPolicyRequest{}
	err := request.Validate()
	assert.Equal(t, "policy id should be provided", err.Error())

	request = GetAlertPolicyRequest{Id: "id"}
	err = request.Validate()
	assert.Nil(t, err)
}

func TestGetNotificationPolicy_Validate(t *testing.T) {
	request := GetNotificationPolicyRequest{}
	err := request.Validate()
	assert.Equal(t, "policy id should be provided", err.Error())

	request.Id = "id"
	err = request.Validate()
	assert.Equal(t, "policy team id should be provided", err.Error())

	request.TeamId = "teamId"
	err = request.Validate()
	assert.Nil(t, err)
}

func TestUpdateAlertPolicy_Validate(t *testing.T) {
	request := UpdateAlertPolicyRequest{}
	err := request.Validate()
	assert.Equal(t, "policy id should be provided", err.Error())

	request = UpdateAlertPolicyRequest{Id: "id"}
	request.PolicyType = "alert"
	request.Name = "asd"
	request.Message = "asdasd"
	err = request.Validate()
	assert.Nil(t, err)
}

func TestUpdateNotificationPolicy_Validate(t *testing.T) {
	request := UpdateNotificationPolicyRequest{}
	err := request.Validate()
	assert.Equal(t, "policy id should be provided", err.Error())

	request = UpdateNotificationPolicyRequest{Id: "id"}
	request.PolicyType = "notification"
	request.Name = "asd"
	err = request.Validate()
	assert.Equal(t, "policy team id should be provided", err.Error())

	request.TeamId = "asdads"
	err = request.Validate()
	assert.Nil(t, err)
}

func TestDeletePolicy_Validate(t *testing.T) {
	request := &DeletePolicyRequest{}
	err := request.Validate()
	assert.Equal(t, "policy type should be one of alert or notification", err.Error())

	request.Type = NotificationPolicy
	err = request.Validate()
	assert.Equal(t, "policy id should be provided", err.Error())

	request.Id = "asd"
	err = request.Validate()
	assert.Equal(t, "policy team id should be provided", err.Error())

	request.TeamId = "asda"
	err = request.Validate()
	assert.Nil(t, err)
}

func TestChangeOrder_Validate(t *testing.T) {
	request := &ChangeOrderRequest{
		Id:          "asd",
		TeamId:      "asd",
		Type:        NotificationPolicy,
		TargetIndex: -1,
	}
	err := request.Validate()
	assert.Equal(t, "target index should be at least 0", err.Error())

	request.TargetIndex = 0
	err = request.Validate()
	assert.Nil(t, err)
}

func TestListNotificationPolicies_Validate(t *testing.T) {
	request := &ListNotificationPoliciesRequest{}
	err := request.Validate()
	assert.Equal(t, "team id should be provided", err.Error())

	request.TeamId = "asdads"
	err = request.Validate()
	assert.Nil(t, err)
}
