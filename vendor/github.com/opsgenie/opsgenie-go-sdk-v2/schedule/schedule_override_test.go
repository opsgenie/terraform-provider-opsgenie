package schedule

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildCreateScheduleOverrideRequest_Validate(t *testing.T) {
	var err error
	createRequest := &CreateScheduleOverrideRequest{}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Schedule identifier cannot be empty.").Error())

	createRequest.ScheduleIdentifier = "id"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User cannot be empty.").Error())

	createRequest.User = Responder{
		Type:     UserResponderType,
		Id:       "id",
		Username: "username",
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Start date cannot be empty.").Error())

	createRequest.StartDate = time.Now()
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("End date cannot be empty.").Error())

	createRequest.EndDate = time.Now()
	err = createRequest.Validate()

}

func TestBuildUpdateScheduleOverrideRequest_Validate(t *testing.T) {
	var err error
	updateRequest := &UpdateScheduleOverrideRequest{}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Schedule identifier cannot be empty.").Error())

	updateRequest.ScheduleIdentifier = "id"
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Alias cannot be empty.").Error())

	updateRequest.Alias = "alias"
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User cannot be empty.").Error())

	updateRequest.User = Responder{
		Type:     UserResponderType,
		Id:       "id",
		Username: "username",
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Start date cannot be empty.").Error())

	updateRequest.StartDate = time.Now()
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("End date cannot be empty.").Error())

	updateRequest.EndDate = time.Now()
	err = updateRequest.Validate()
}
