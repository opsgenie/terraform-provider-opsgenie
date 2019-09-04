package og

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
	"time"
)

type OwnerTeam struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Rotation struct {
	client.BaseRequest
	Name            string           `json:"name,omitempty"`
	StartDate       *time.Time       `json:"startDate,omitempty"`
	EndDate         *time.Time       `json:"endDate,omitempty"`
	Type            RotationType     `json:"type,omitempty"`
	Length          uint32           `json:"length,omitempty"`
	Participants    []Participant    `json:"participants,omitempty"`
	TimeRestriction *TimeRestriction `json:"timeRestriction,omitempty"`
}

func (r *Rotation) Validate() error {

	if r.Type == "" {
		return errors.New("Rotation type cannot be empty.")
	}
	if r.StartDate == nil {
		return errors.New("Rotation start date cannot be empty.")
	}
	if r.EndDate != nil && !r.StartDate.Before(*r.EndDate) {
		return errors.New("Rotation end time should be later than start time.")
	}
	if len(r.Participants) == 0 {
		return errors.New("Rotation participants cannot be empty.")
	}
	err := validateParticipants(r)
	if err != nil {
		return err
	}
	if r.TimeRestriction != nil {
		err := ValidateRestrictions(r.TimeRestriction)
		if err != nil {
			return err
		}
	}

	return nil

}

func ValidateRotations(rotations []Rotation) error {
	for _, rotation := range rotations {

		err := rotation.Validate()

		if err != nil {
			return err
		}
	}
	return nil
}

func validateParticipants(rotation *Rotation) error {
	for _, participant := range rotation.Participants {
		if participant.Type == "" {
			return errors.New("Participant type cannot be empty.")
		}
		if !(participant.Type == User || participant.Type == Team) {
			return errors.New("Participant type should be one of these: 'User', 'Team'")
		}
		if participant.Type == User && participant.Username == "" && participant.Id == "" {
			return errors.New("For participant type user either username or id must be provided.")
		}
		if participant.Type == Team && participant.Name == "" && participant.Id == "" {
			return errors.New("For participant type team either team name or id must be provided.")
		}
	}
	return nil
}

func (r *Rotation) WithParticipant(participant Participant) *Rotation {
	r.Participants = append(r.Participants, participant)
	return r
}

func (r *Rotation) WithParticipants(participant ...Participant) *Rotation {
	r.Participants = participant
	return r
}

func (r *Rotation) WithTimeRestriction(timeRestriction TimeRestriction) *Rotation {
	r.TimeRestriction = &timeRestriction
	return r
}

func (tr *TimeRestriction) WithRestrictions(restrictions ...Restriction) *TimeRestriction {
	tr.RestrictionList = restrictions
	return tr
}

func ValidateFilter(filter Filter) error {
	if filter.ConditionMatchType != MatchAll && filter.ConditionMatchType != MatchAllConditions && filter.ConditionMatchType != MatchAnyCondition {
		return errors.New("filter condition type should be one of match-all, match-any-condition or match-all-conditions")
	}
	if (filter.ConditionMatchType == MatchAllConditions || filter.ConditionMatchType == MatchAnyCondition) && len(filter.Conditions) == 0 {
		return errors.New("filter conditions cannot be empty")
	}
	if len(filter.Conditions) > 0 {
		err := ValidateConditions(filter.Conditions)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateConditions(conditions []Condition) error {
	for _, condition := range conditions {
		if condition.Field != ExtraProperties && condition.Key != "" {
			return errors.New("condition key is only valid for extra-properties field")
		}
		switch condition.Field {
		case Message, Alias, Description, Source, Entity, Tags, Actions, Details, ExtraProperties, Recipients, Teams, Priority:
			break
		default:
			return errors.New("condition field should be one of message, alias, description, source, entity, tags, actions, details, extra-properties, recipients, teams or priority")
		}
		switch condition.Field {
		case Actions, Tags, Recipients:
			if condition.Operation != Contains && condition.Operation != IsEmpty && condition.Operation != Matches {
				return errors.New(string(condition.Operation) + " is not valid operation for " + string(condition.Field))
			}
		case Message, Alias, Description, Source, Entity, Teams:
			if condition.Operation != Contains && condition.Operation != IsEmpty && condition.Operation != Matches &&
				condition.Operation != Equals && condition.Operation != StartsWith && condition.Operation != EndsWith &&
				condition.Operation != EqualsIgnoreWhitespcae {
				return errors.New(string(condition.Operation) + " is not valid operation for " + string(condition.Field))
			}
		case Details:
			if condition.Operation != Contains && condition.Operation != IsEmpty && condition.Operation != ContainsKey &&
				condition.Operation != ContainsValue {
				return errors.New(string(condition.Operation) + " is not valid operation for " + string(condition.Field))
			}
		case Priority:
			if condition.Operation != Equals && condition.Operation != GreaterThan && condition.Operation != LessThan {
				return errors.New(string(condition.Operation) + " is not valid operation for " + string(condition.Field))
			}
			if condition.ExpectedValue != string(alert.P1) && condition.ExpectedValue != string(alert.P2) && condition.ExpectedValue != string(alert.P3) &&
				condition.ExpectedValue != string(alert.P4) && condition.ExpectedValue != string(alert.P5) {
				return errors.New("for field " + string(condition.Field) + " expected value should be one of P1, P2, P3, P4, P5")
			}
		}
	}
	return nil
}

func ValidateRestrictions(timeRestriction *TimeRestriction) error {
	if timeRestriction.Type == "" {
		return errors.New("Type of time restriction must be time-of-day or weekday-and-time-of-day.")
	}
	if timeRestriction.Type == WeekdayAndTimeOfDay && timeRestriction.RestrictionList == nil {
		return errors.New("Restrictions cannot be empty.")
	}
	if len(timeRestriction.RestrictionList) != 0 {
		for _, restriction := range timeRestriction.RestrictionList {
			if timeRestriction.Type == "weekday-and-time-of-day" &&
				(restriction.EndMin < 0 ||
					restriction.EndHour <= 0 ||
					restriction.EndDay == "" ||
					restriction.StartDay == "" ||
					restriction.StartHour <= 0 ||
					restriction.StartMin < 0) {
				return errors.New("startDay, startHour, startMin, endDay, endHour, endMin cannot be empty.")
			}
		}
	}
	if timeRestriction.Type == TimeOfDay &&
		(timeRestriction.Restriction.EndMin < 0 ||
			timeRestriction.Restriction.EndHour <= 0 ||
			timeRestriction.Restriction.StartHour <= 0 ||
			timeRestriction.Restriction.StartMin < 0) {
		return errors.New("startHour, startMin, endHour, endMin cannot be empty.")
	}

	return nil
}

type RotationType string
type ParticipantType string
type Day string
type RestrictionType string
type TimeUnit string

const (
	Daily  RotationType = "daily"
	Weekly RotationType = "weekly"
	Hourly RotationType = "hourly"

	User       ParticipantType = "user"
	Team       ParticipantType = "team"
	Escalation ParticipantType = "escalation"
	Schedule   ParticipantType = "schedule"
	None       ParticipantType = "none"

	Monday    Day = "monday"
	Tuesday   Day = "tuesday"
	Wednesday Day = "wednesday"
	Thursday  Day = "thursday"
	Friday    Day = "friday"
	Saturday  Day = "saturday"
	Sunday    Day = "sunday"

	TimeOfDay           RestrictionType = "time-of-day"
	WeekdayAndTimeOfDay RestrictionType = "weekday-and-time-of-day"

	MatchAll           ConditionMatchType = "match-all"
	MatchAnyCondition  ConditionMatchType = "match-any-condition"
	MatchAllConditions ConditionMatchType = "match-all-conditions"

	Months  TimeUnit = "months"
	Weeks   TimeUnit = "weeks"
	Days    TimeUnit = "days"
	Minutes TimeUnit = "minutes"
	Hours   TimeUnit = "hours"

	Message         ConditionFieldType = "message"
	Alias           ConditionFieldType = "alias"
	Description     ConditionFieldType = "description"
	Source          ConditionFieldType = "source"
	Entity          ConditionFieldType = "entity"
	Tags            ConditionFieldType = "tags"
	Actions         ConditionFieldType = "actions"
	Details         ConditionFieldType = "details"
	ExtraProperties ConditionFieldType = "extra-properties"
	Recipients      ConditionFieldType = "recipients"
	Teams           ConditionFieldType = "teams"
	Priority        ConditionFieldType = "priority"

	Matches                ConditionOperation = "matches"
	Contains               ConditionOperation = "contains"
	StartsWith             ConditionOperation = "starts-with"
	EndsWith               ConditionOperation = "ends-with"
	Equals                 ConditionOperation = "equals"
	ContainsKey            ConditionOperation = "contains-key"
	ContainsValue          ConditionOperation = "contains-value"
	GreaterThan            ConditionOperation = "greater-than"
	LessThan               ConditionOperation = "less-than"
	IsEmpty                ConditionOperation = "is-empty"
	EqualsIgnoreWhitespcae ConditionOperation = "equals-ignore-whitespace"
)

type Identifier interface {
	identifier() string
	identifierType() string
}

type Participant struct {
	Type     ParticipantType `json:"type, omitempty"`
	Name     string          `json:"name,omitempty"`
	Id       string          `json:"id,omitempty"`
	Username string          `json:"username, omitempty"`
}

type TimeRestriction struct {
	Type            RestrictionType `json:"type,omitempty"`
	RestrictionList []Restriction   `json:"restrictions,omitempty"`
	Restriction     Restriction     `json:"restriction,omitempty"`
}

type Restriction struct {
	StartDay  Day    `json:"startDay,omitempty"`
	StartHour uint32 `json:"startHour,omitempty"`
	StartMin  uint32 `json:"startMin,omitempty"`
	EndHour   uint32 `json:"endHour,omitempty"`
	EndDay    Day    `json:"endDay,omitempty"`
	EndMin    uint32 `json:"endMin,omitempty"`
}

type Filter struct {
	ConditionMatchType ConditionMatchType `json:"type,omitempty"`
	Conditions         []Condition        `json:"conditions,omitempty"`
}

type Condition struct {
	Field         ConditionFieldType `json:"field,omitempty"`
	IsNot         bool               `json:"not,omitempty"`
	Operation     ConditionOperation `json:"operation,omitempty"`
	ExpectedValue string             `json:"expectedValue,omitempty"`
	Key           string             `json:"key,omitempty"`
	Order         *int               `json:"order,omitempty"`
}

type ConditionMatchType string
type ConditionFieldType string
type ConditionOperation string

type Contact struct {
	To              string     `json:"to,omitempty"`
	MethodOfContact MethodType `json:"method,omitempty"`
}

type MethodType string

const (
	Sms    MethodType = "sms"
	Email  MethodType = "email"
	Voice  MethodType = "voice"
	Mobile MethodType = "mobile"
)

type SendAfter struct {
	TimeAmount uint32 `json:"timeAmount,omitempty"`
	TimeUnit   string `json:"timeUnit,omitempty"`
}

type Step struct {
	Contact   Contact    `json:"contact,omitempty"`
	SendAfter *SendAfter `json:"sendAfter,omitempty"`
	Enabled   bool       `json:"enabled,omitempty"`
}

type Criteria struct {
	CriteriaType ConditionMatchType `json:"type"`
	Conditions   []Condition        `json:"conditions,omitempty"`
}

type NotifyType string

const (
	Next     NotifyType = "next"
	Previous NotifyType = "previous"
	Default  NotifyType = "default"
	Users    NotifyType = "users"
	Admins   NotifyType = "admins"
	All      NotifyType = "all"
)

type EscalationCondition string

const (
	IfNotAcked  EscalationCondition = "if-not-acked"
	IfNotClosed EscalationCondition = "if-not-closed"
)
