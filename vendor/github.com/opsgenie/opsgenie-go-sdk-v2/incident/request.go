package incident

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

type RequestStatusRequest struct {
	client.BaseRequest
	Id string
}

func (r *RequestStatusRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	return nil
}

func (r *RequestStatusRequest) ResourcePath() string {
	return "/v1/incidents/requests/" + r.Id
}

func (r *RequestStatusRequest) Method() string {
	return http.MethodGet
}

type CreateRequest struct {
	client.BaseRequest
	Message            string            `json:"message"`
	Description        string            `json:"description,omitempty"`
	Responders         []Responder       `json:"responders,omitempty"`
	Tags               []string          `json:"tags,omitempty"`
	Details            map[string]string `json:"details,omitempty"`
	Priority           Priority          `json:"priority,omitempty"`
	Note               string            `json:"note,omitempty"`
	ServiceId          string            `json:"serviceId"`
	StatusPageEntity   *StatusPageEntity `json:"statusPageEntry,omitempty"`
	NotifyStakeholders *bool             `json:"notifyStakeholders,omitempty"`
}

type StatusPageEntity struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.Message == "" || r.ServiceId == "" {
		return errors.New("Message and ServiceId fields cannot be blank.")
	}
	if r.StatusPageEntity != nil {
		if r.StatusPageEntity.Title == "" {
			return errors.New("StatusPageEntity.Title cannot be blank.")
		}
	}
	err := ValidatePriority(r.Priority)
	if err != nil {
		return err
	}
	err = validateResponders(r.Responders)
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateRequest) ResourcePath() string {
	return "/v1/incidents/create"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type DeleteRequest struct {
	client.BaseRequest
	Id         string
	Identifier IdentifierType
}

func (r *DeleteRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *DeleteRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type GetRequest struct {
	client.BaseRequest
	Id         string
	Identifier IdentifierType
}

func (r *GetRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

func (r *GetRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}
	return params
}

type ListRequest struct {
	client.BaseRequest
	Limit  int
	Sort   SortField
	Offset int
	Order  Order
	Query  string
}

func (r *ListRequest) Validate() error {
	if r.Query == "" {
		return errors.New("Query field cannot be empty.")
	}
	return nil
}

func (r *ListRequest) ResourcePath() string {
	return "/v1/incidents"
}

func (r *ListRequest) Method() string {
	return http.MethodGet
}

func (r *ListRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Limit != 0 {
		params["limit"] = strconv.Itoa(r.Limit)
	}
	if r.Sort != "" {
		params["sort"] = string(r.Sort)
	}
	if r.Offset != 0 {
		params["offset"] = strconv.Itoa(r.Offset)
	}
	if r.Query != "" {
		params["query"] = r.Query

	}
	if r.Order != "" {
		params["order"] = string(r.Order)

	}
	return params
}

type CloseRequest struct {
	client.BaseRequest
	Id         string
	Identifier IdentifierType
	Note       string `json:"note,omitempty"`
}

func (r *CloseRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *CloseRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/close"
}

func (r *CloseRequest) Method() string {
	return http.MethodPost
}

func (r *CloseRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type AddNoteRequest struct {
	client.BaseRequest
	Id         string
	Identifier IdentifierType
	Note       string `json:"note"`
}

func (r *AddNoteRequest) Validate() error {
	if r.Id == "" || r.Note == "" {
		return errors.New("Incident ID and Note fields cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *AddNoteRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/notes"

}

func (r *AddNoteRequest) Method() string {
	return http.MethodPost
}

func (r *AddNoteRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type AddResponderRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string      `json:"incidentId"`
	Note       string      `json:"note"`
	Responders []Responder `json:"responder"`
}

func (r *AddResponderRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if len(r.Responders) == 0 {
		return errors.New("Responders field cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	err := validateResponders(r.Responders)
	if err != nil {
		return err
	}
	return nil
}

func (r *AddResponderRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/responders"

}

func (r *AddResponderRequest) Method() string {
	return http.MethodPost
}

func (r *AddResponderRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type AddTagsRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Note       string   `json:"note"`
	Tags       []string `json:"tags"`
}

func (r *AddTagsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if len(r.Tags) == 0 {
		return errors.New("Tags field cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *AddTagsRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/tags"
}

func (r *AddTagsRequest) Method() string {
	return http.MethodPost
}

func (r *AddTagsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type RemoveTagsRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Note       string
	Tags       []string
}

func (r *RemoveTagsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if len(r.Tags) == 0 {
		return errors.New("Tags field cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *RemoveTagsRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/tags"
}

func (r *RemoveTagsRequest) Method() string {
	return http.MethodDelete
}

func (r *RemoveTagsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	//comma separated tag list
	params["tags"] = strings.Join(r.Tags[:], ",")

	if r.Note != "" {
		params["note"] = r.Note

	}
	return params
}

type AddDetailsRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Note       string            `json:"note,omitempty"`
	Details    map[string]string `json:"details"`
}

func (r *AddDetailsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if len(r.Details) == 0 {
		return errors.New("Details field cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *AddDetailsRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/details"
}

func (r *AddDetailsRequest) Method() string {
	return http.MethodPost
}

func (r *AddDetailsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type RemoveDetailsRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Note       string
	Keys       []string
}

func (r *RemoveDetailsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if len(r.Keys) == 0 {
		return errors.New("Details field cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *RemoveDetailsRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/details"
}

func (r *RemoveDetailsRequest) Method() string {
	return http.MethodDelete
}

func (r *RemoveDetailsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	//comma separated key list
	params["keys"] = strings.Join(r.Keys[:], ",")

	if r.Note != "" {
		params["note"] = r.Note

	}
	return params
}

type UpdatePriorityRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Priority   Priority `json:"priority"`
}

func (r *UpdatePriorityRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	err := ValidatePriority(r.Priority)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdatePriorityRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/priority"

}

func (r *UpdatePriorityRequest) Method() string {
	return http.MethodPut
}

func (r *UpdatePriorityRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type UpdateMessageRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Message    string `json:"message"`
}

func (r *UpdateMessageRequest) Validate() error {
	if r.Id == "" || r.Message == "" {
		return errors.New("Incident ID and Message fields cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *UpdateMessageRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/message"

}

func (r *UpdateMessageRequest) Method() string {
	return http.MethodPost
}

func (r *UpdateMessageRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type UpdateDescriptionRequest struct {
	client.BaseRequest
	Identifier  IdentifierType
	Id          string
	Description string `json:"description"`
}

func (r *UpdateDescriptionRequest) Validate() error {
	if r.Id == "" || r.Description == "" {
		return errors.New("Incident ID and Description fields cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *UpdateDescriptionRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/description"
}

func (r *UpdateDescriptionRequest) Method() string {
	return http.MethodPost
}

func (r *UpdateDescriptionRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type ListLogsRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Limit      int
	Offset     int
	Order      Order
	Direction  string
}

func (r *ListLogsRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *ListLogsRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/logs"
}

func (r *ListLogsRequest) Method() string {
	return http.MethodGet
}

func (r *ListLogsRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	if r.Limit != 0 {
		params["limit"] = strconv.Itoa(r.Limit)
	}
	if r.Offset != 0 {
		params["offset"] = strconv.Itoa(r.Offset)
	}
	if r.Direction != "" {
		params["direction"] = r.Direction

	}
	if r.Order != "" {
		params["order"] = string(r.Order)
	}

	return params
}

type ListNotesRequest struct {
	client.BaseRequest
	Identifier IdentifierType
	Id         string
	Limit      int
	Offset     int
	Order      Order
	Direction  string
}

func (r *ListNotesRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != Id && r.Identifier != Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *ListNotesRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/notes"

}

func (r *ListNotesRequest) Method() string {
	return http.MethodGet
}

func (r *ListNotesRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Identifier == Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	if r.Limit != 0 {
		params["limit"] = strconv.Itoa(r.Limit)
	}
	if r.Offset != 0 {
		params["offset"] = strconv.Itoa(r.Offset)
	}
	if r.Direction != "" {
		params["direction"] = r.Direction

	}
	if r.Order != "" {
		params["order"] = string(r.Order)
	}

	return params
}

type IdentifierType string
type ResponderType string
type Priority string
type Order string
type SortField string

const (
	Id   IdentifierType = "id"
	Tiny IdentifierType = "tiny"

	User ResponderType = "user"
	Team ResponderType = "team"

	P1 Priority = "P1"
	P2 Priority = "P2"
	P3 Priority = "P3"
	P4 Priority = "P4"
	P5 Priority = "P5"

	Asc  Order = "asc"
	Desc Order = "desc"

	CreatedAt SortField = "createdAt"
	TinyId    SortField = "tinyId"
	Message   SortField = "message"
	Status    SortField = "status"
	IsSeen    SortField = "isSeen"
	Owner     SortField = "owner"
)

type Responder struct {
	Type ResponderType `json:"type,omitempty"`
	Name string        `json:"name,omitempty"`
	Id   string        `json:"id,omitempty"`
}

func validateResponders(responders []Responder) error {
	for _, responder := range responders {
		if responder.Type == "" {
			return errors.New("Responder type cannot be empty.")
		}
		if !(responder.Type == User || responder.Type == Team) {
			return errors.New("Responder type should be one of these: 'User', 'Team'.")
		}
		if responder.Name == "" && responder.Id == "" {
			return errors.New("For responder either name or id must be provided.")
		}
	}
	return nil
}

func ValidatePriority(priority Priority) error {
	switch priority {
	case P1, P2, P3, P4, P5, "":
		return nil
	}
	return errors.New("Priority should be one of these: " +
		"'P1', 'P2', 'P3', 'P4' and 'P5' or empty")
}
