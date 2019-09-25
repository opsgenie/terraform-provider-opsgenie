package contact

import "github.com/opsgenie/opsgenie-go-sdk-v2/client"

type Contact struct {
	Id              string `json:"id"`
	MethodOfContact string `json:"method"`
	To              string `json:"to,omitempty"`
	Status          Status `json:"status,omitempty"`
	ApplyOrder      uint32 `json:"applyOrder,omitempty"`
}

type Status struct {
	Enabled        bool   `json:"enabled"`
	DisabledReason string `json:"disabledReason"`
}

type CreateResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type GetResult struct {
	client.ResultMetadata
	Id              string `json:"id"`
	MethodOfContact string `json:"method"`
	To              string `json:"to,omitempty"`
	Status          Status `json:"status,omitempty"`
}

type UpdateResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result string `json:"result,omitempty"`
}

type ListResult struct {
	client.ResultMetadata
	Contact []Contact `json:"data,omitempty"`
}

type EnableResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type DisableResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}
