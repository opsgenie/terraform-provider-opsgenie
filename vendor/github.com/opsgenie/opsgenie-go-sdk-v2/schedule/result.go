package schedule

import (
	"errors"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"io/ioutil"
	"net/http"
	"time"
)

type Schedule struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Timezone    string        `json:"timezone,omitempty"`
	Enabled     bool          `json:"enabled"`
	OwnerTeam   *og.OwnerTeam `json:"ownerTeam,omitempty"`
	Rotations   []og.Rotation `json:"rotations,omitempty"`
}
type CreateResult struct {
	client.ResultMetadata
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type GetResult struct {
	client.ResultMetadata
	Schedule Schedule `json:"data,omitempty"`
}

type UpdateResult struct {
	client.ResultMetadata
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result string `json:"result,omitempty"`
}

type ListResult struct {
	client.ResultMetadata
	Schedule         []Schedule `json:"data,omitempty"`
	ExpandableFields []string   `json:"expandable,omitempty"`
}

type TimelineResult struct {
	client.ResultMetadata
	ScheduleInfo       Info         `json:"_parent"`
	Description        string       `json:"description"`
	OwnerTeam          og.OwnerTeam `json:"ownerTeam,omitempty"`
	StartDate          time.Time    `json:"startDate,omitempty"`
	EndDate            time.Time    `json:"endDate,omitempty"`
	FinalTimeline      Timeline     `json:"finalTimeline,omitempty"`
	BaseTimeline       Timeline     `json:"baseTimeline,omitempty"`
	OverrideTimeline   Timeline     `json:"overrideTimeline,omitempty"`
	ForwardingTimeline Timeline     `json:"forwardingTimeline,omitempty"`
	ExpandableFields   []string     `json:"expandable,omitempty"`
}

type Timeline struct {
	Rotations []TimelineRotation `json:"rotations,omitempty"`
}

type TimelineRotation struct {
	Id      string   `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Order   float32  `json:"order,omitempty"`
	Periods []Period `json:"periods,omitempty"`
}

type exportScheduleResult struct {
	client.ResultMetadata
	FileContent []byte
}

func (rm *exportScheduleResult) Parse(response *http.Response, result client.ApiResult) error {

	if response == nil {
		return errors.New("No response received")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	result.(*exportScheduleResult).FileContent = body

	return nil
}

type Info struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type Period struct {
	StartDate time.Time      `json:"startDate,omitempty"`
	EndDate   time.Time      `json:"endDate,omitempty"`
	Type      string         `json:"type,omitempty"`
	Recipient og.Participant `json:"recipient,omitempty"`
}

type Rotation struct {
	Id              string              `json:"id,omitempty"`
	Name            string              `json:"name,omitempty"`
	StartDate       *time.Time          `json:"startDate,omitempty"`
	EndDate         *time.Time          `json:"endDate,omitempty"`
	Type            og.RotationType     `json:"type,omitempty"`
	Length          uint32              `json:"length,omitempty"`
	Participants    []og.Participant    `json:"participants,omitempty"`
	TimeRestriction *og.TimeRestriction `json:"timeRestriction,omitempty"`
}

type CreateRotationResult struct {
	client.ResultMetadata
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GetRotationResult struct {
	client.ResultMetadata
	Rotation
	Info `json:"_parent,omitempty"`
}

type UpdateRotationResult struct {
	client.ResultMetadata
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ListRotationsResult struct {
	client.ResultMetadata
	Rotations []Rotation `json:"data,omitempty"`
}
