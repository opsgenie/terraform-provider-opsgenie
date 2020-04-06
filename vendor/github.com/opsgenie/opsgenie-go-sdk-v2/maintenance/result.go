package maintenance

import "github.com/opsgenie/opsgenie-go-sdk-v2/client"

type Maintenance struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	Time        Time   `json:"time"`
	Description string `json:"description"`
}

type CreateResult struct {
	client.ResultMetadata
	Maintenance
}

type UpdateResult struct {
	client.ResultMetadata
	Maintenance
}

type ChangeEndDateResult struct {
	client.ResultMetadata
	Maintenance
}

type GetResult struct {
	client.ResultMetadata
	Id          string `json:"id"`
	Status      string `json:"status"`
	Time        Time   `json:"time"`
	Description string `json:"description"`
	Results     []Rule `json:"rules"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type CloseResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type ListResult struct {
	client.ResultMetadata
	Maintenances []Maintenance `json:"data"`
}
