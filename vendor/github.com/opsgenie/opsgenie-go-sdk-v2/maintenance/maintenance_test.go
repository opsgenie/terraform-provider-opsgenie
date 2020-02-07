package maintenance

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateRequest_Validate(t *testing.T) {
	request := &CreateRequest{}
	err := request.Validate()
	assert.Equal(t, err.Error(), errors.New("Time.Type should be one of these: "+
		"'For5Minutes', 'For30Minutes', 'For1Hour', 'Indefinitely' and 'Schedule'").Error())
	time := Time{
		Type: For1Hour,
	}
	request.Time = time
	err = request.Validate()
	assert.Equal(t, err.Error(), errors.New("There should be at least one rule.").Error())
	rules := []Rule{
		{
			State: Enabled,
			Entity: Entity{
				Id:   "e14dda76-488e-4e98-a1c7-78cda1900e27",
				Type: Policy,
			},
		},
	}
	request.Rules = rules
	err = request.Validate()
	assert.Nil(t, err)
}

func TestGetRequest_Validate(t *testing.T) {
	request := &GetRequest{}
	err := request.Validate()
	assert.Equal(t, err.Error(), errors.New("Maintenance ID cannot be blank.").Error())
	request.Id = "e14dda76-488e-4e98-a1c7-78cda1900e27"
	err = request.Validate()
	assert.Nil(t, err)
}

func TestUpdateRequest_Validate(t *testing.T) {
	now := time.Now()
	after := now.Add(time.Minute)
	request := &UpdateRequest{}
	err := request.Validate()
	assert.Equal(t, err.Error(), errors.New("Maintenance ID cannot be blank.").Error())
	request.Id = "e14dda76-488e-4e98-a1c7-78cda1900e27"
	err = request.Validate()
	assert.Equal(t, err.Error(), errors.New("Time.Type should be one of these: "+
		"'For5Minutes', 'For30Minutes', 'For1Hour', 'Indefinitely' and 'Schedule'").Error())
	request.Time = Time{
		Type: Schedule,
	}
	err = request.Validate()
	assert.Equal(t, err.Error(), errors.New("For 'Schedule' type both 'StartDate' and 'EndDate' fields cannot be empty.").Error())
	request.Time.StartDate = &now
	request.Time.EndDate = &after
	err = request.Validate()
	assert.Equal(t, err.Error(), errors.New("There should be at least one rule.").Error())
	rules := []Rule{
		{
			State: Disabled,
			Entity: Entity{
				Id:   "e14dda76-488e-4e98-a1c7-78cda1900e27",
				Type: Policy,
			},
		},
	}
	request.Rules = rules
	err = request.Validate()
	assert.Nil(t, err)
}

func TestDeleteRequest_Validate(t *testing.T) {
	request := &DeleteRequest{}
	err := request.Validate()
	assert.Equal(t, err.Error(), errors.New("Maintenance ID cannot be blank.").Error())
	request.Id = "e14dda76-488e-4e98-a1c7-78cda1900e27"
	err = request.Validate()
	assert.Nil(t, err)
}

func TestListRequest_Validate(t *testing.T) {
	request := &ListRequest{}
	err := request.Validate()
	assert.Nil(t, err)
}

func TestCancelRequest_Validate(t *testing.T) {
	request := &CancelRequest{}
	err := request.Validate()
	assert.Equal(t, err.Error(), errors.New("Maintenance ID cannot be blank.").Error())
	request.Id = "e14dda76-488e-4e98-a1c7-78cda1900e27"
	err = request.Validate()
	assert.Nil(t, err)
}

func TestTime_Validate(t *testing.T) {
	now := time.Now()
	before := time.Now().Add(-time.Hour)
	after := time.Now().Add(time.Hour)
	tf := Time{
		Type: "cem",
	}
	err := validateTime(tf)
	assert.Equal(t, err.Error(), errors.New("Time.Type should be one of these: 'For5Minutes',"+
		" 'For30Minutes', 'For1Hour', 'Indefinitely' and 'Schedule'").Error())
	tf.Type = Indefinitely
	err = validateTime(tf)
	assert.Nil(t, err)
	tf.Type = Schedule
	err = validateTime(tf)
	assert.Equal(t, err.Error(), errors.New("For 'Schedule' type both 'StartDate' and"+
		" 'EndDate' fields cannot be empty.").Error())
	tf.StartDate = &now
	err = validateTime(tf)
	assert.Equal(t, err.Error(), errors.New("For 'Schedule' type both 'StartDate' and"+
		" 'EndDate' fields cannot be empty.").Error())
	tf.EndDate = &before
	err = validateTime(tf)
	assert.Equal(t, err.Error(), errors.New("EndDate should be after the StartDate.").Error())
	tf.EndDate = &after
	err = validateTime(tf)
	assert.Nil(t, err)
}

func TestRules_Validate(t *testing.T) {
	rules := []Rule{
		{},
	}
	err := validateRules(rules)
	assert.Equal(t, err.Error(), errors.New("Rule.Entity.Id should be one of these: 'Policy', 'Integration'.").Error())
	rules = []Rule{
		{
			Entity: Entity{
				Type: Policy,
			},
		},
	}
	err = validateRules(rules)
	assert.Equal(t, err.Error(), errors.New("Rule.State field cannot be empty when the"+
		" Rule.Entity.Type is not 'Integration'.").Error())
	rules[0].State = Enabled
	err = validateRules(rules)
	assert.Nil(t, err)
	rules = append(rules,
		Rule{
			Entity: Entity{
				Type: "cem",
			},
			State: "cem",
		},
	)
	err = validateRules(rules)
	assert.Equal(t, err.Error(), errors.New("Rule.Entity.Id should be one of these:"+
		" 'Policy', 'Integration'.").Error())
	rules[1].Entity.Type = Integration
	err = validateRules(rules)
	assert.Equal(t, err.Error(), errors.New("Rule.State field should be one of these: "+
		"'Enabled' and 'Disabled' or empty.").Error())
	rules[1].State = Enabled
	err = validateRules(rules)
	assert.Equal(t, err.Error(), errors.New("Rule.State field cannot be 'Enabled' when "+
		"the Rule.Entity.Type is 'Integration'.").Error())
	rules[1].State = ""
	err = validateRules(rules)
	assert.Nil(t, err)
}

func TestStatusType_Validate(t *testing.T) {
	err := validateStatusType("cem")
	assert.Equal(t, err.Error(), errors.New("Priority should be one of these: "+
		"'All', 'NonExpired' and 'Past' or empty.").Error())
	err = validateStatusType("")
	assert.Nil(t, err)
	err = validateStatusType(All)
	assert.Nil(t, err)
	err = validateStatusType(NonExpired)
	assert.Nil(t, err)
	err = validateStatusType(Past)
	assert.Nil(t, err)
}
