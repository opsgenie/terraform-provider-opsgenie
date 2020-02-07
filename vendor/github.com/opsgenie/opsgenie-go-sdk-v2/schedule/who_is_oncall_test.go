package schedule

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOnCallsRequest_Validate(t *testing.T) {
	onCallsRequest := &GetOnCallsRequest{}
	err := onCallsRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Schedule identifier cannot be empty.").Error())

	onCallsRequest = &GetOnCallsRequest{
		ScheduleIdentifierType: Id,
		ScheduleIdentifier:     "id1",
	}
	err = onCallsRequest.Validate()

	assert.Equal(t, err, nil)

}

func TestGetNextOnCallsRequest_Validate(t *testing.T) {
	getNextOnCallsRequest := &GetNextOnCallsRequest{}
	err := getNextOnCallsRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Schedule identifier cannot be empty.").Error())
	getNextOnCallsRequest = &GetNextOnCallsRequest{
		ScheduleIdentifierType: Id,
		ScheduleIdentifier:     "id1",
	}
	err = getNextOnCallsRequest.Validate()

	assert.Equal(t, err, nil)

}

func TestExportOnCallUserRequest_Validate(t *testing.T) {
	nextOnCallsRequest := &ExportOnCallUserRequest{}
	err := nextOnCallsRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())
	nextOnCallsRequest = &ExportOnCallUserRequest{
		UserIdentifier: "user1",
	}
	err = nextOnCallsRequest.Validate()

	assert.Equal(t, err, nil)

}
