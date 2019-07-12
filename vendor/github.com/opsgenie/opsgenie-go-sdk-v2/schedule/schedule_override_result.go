package schedule

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"time"
)

type ScheduleMeta struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type ScheduleOverride struct {
	Parent    ScheduleMeta         `json:"_parent,omitempty"`
	Alias     string               `json:"alias"`
	User      Responder            `json:"user"`
	StartDate time.Time            `json:"startDate"`
	EndDate   time.Time            `json:"endDate"`
	Rotations []RotationIdentifier `json:"rotations"`
}

type CreateScheduleOverrideResult struct {
	client.ResultMetadata
	Alias string `json:"alias"`
}

type GetScheduleOverrideResult struct {
	client.ResultMetadata
	ScheduleOverride
}

type ListScheduleOverrideResult struct {
	client.ResultMetadata
	ScheduleOverride []ScheduleOverride `json:"data,omitempty"`
}

type DeleteScheduleOverrideResult struct {
	client.ResultMetadata
	Result string `json:"result,omitempty"`
}

type UpdateScheduleOverrideResult struct {
	client.ResultMetadata
	Alias string `json:"alias"`
}
