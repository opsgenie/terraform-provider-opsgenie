package schedule

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type GetOnCallsResult struct {
	client.ResultMetadata
	Parent             ScheduleMeta           `json:"_parent,omitempty"`
	OnCallParticipants []GetOnCallParticipant `json:"onCallParticipants,omitempty"`
	OnCallRecipients   []string               `json:"onCallRecipients,omitempty"`
}

type GetOnCallParticipant struct {
	Type               og.ParticipantType  `json:"type, omitempty"`
	Name               string              `json:"name,omitempty"`
	Id                 string              `json:"id,omitempty"`
	EscalationTime     uint32              `json:"escalationTime,omitempty"`
	NotifyType         og.NotifyType       `json:"notifyType,omitempty"`
	OnCallParticipants []OnCallParticipant `json:"onCallParticipants,omitempty"`
}

type OnCallParticipant struct {
	Type           og.ParticipantType `json:"type, omitempty"`
	Name           string             `json:"name,omitempty"`
	Id             string             `json:"id,omitempty"`
	EscalationTime uint32             `json:"escalationTime,omitempty"`
	NotifyType     og.NotifyType      `json:"notifyType,omitempty"`
}

type GetNextOnCallsResult struct {
	client.ResultMetadata
	Parent                      ScheduleMeta           `json:"_parent,omitempty"`
	NextOnCallRecipients        []NextOnCallRecipients `json:"nextOnCallRecipients,omitempty"`
	ExactNextOnCallRecipients   []NextOnCallRecipients `json:"exactNextOnCallRecipients,omitempty"`
	NextOncallParticipants      []string               `json:"nextOnCallParticipants,omitempty"`
	ExactNextOnCallParticipants []string               `json:"exactNextOnCallParticipants,omitempty"`
}

type NextOnCallRecipients struct {
	Type               og.ParticipantType  `json:"type, omitempty"`
	Name               string              `json:"name,omitempty"`
	Id                 string              `json:"id,omitempty"`
	ForwardedFrom      []og.Participant    `json:"forwardedFrom,omitempty"`
	OnCallParticipants []OnCallParticipant `json:"onCallParticipants,omitempty"`
}

type exportOncallUserResult struct {
	client.ResultMetadata
	FileContent []byte
}

func (rm *exportOncallUserResult) Parse(response *http.Response, result client.ApiResult) error {

	if response == nil {
		return errors.New("No response received")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	result.(*exportOncallUserResult).FileContent = body

	return nil
}
