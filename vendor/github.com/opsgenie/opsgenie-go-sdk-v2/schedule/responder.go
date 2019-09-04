package schedule

type ResponderType string

const (
	UserResponderType       ResponderType = "user"
	TeamResponderType       ResponderType = "team"
	EscalationResponderType ResponderType = "escalation"
	ScheduleResponderType   ResponderType = "schedule"
)

type Responder struct {
	Type     ResponderType `json:"type, omitempty"`
	Name     string        `json:"name,omitempty"`
	Id       string        `json:"id,omitempty"`
	Username string        `json:"username, omitempty"`
}

type TeamResponder struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserResponder struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type EscalationResponder struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ScheduleResponder struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (s *ScheduleResponder) SetID(id string) {
	s.ID = id
}

func (s *ScheduleResponder) SetUsername(name string) {
	s.Name = name
}
