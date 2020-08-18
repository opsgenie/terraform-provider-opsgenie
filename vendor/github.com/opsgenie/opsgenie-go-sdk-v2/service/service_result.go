package service

import "github.com/opsgenie/opsgenie-go-sdk-v2/client"

type Service struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Visibility  Visibility `json:"visibility"`
	TeamId      string     `json:"teamId"`
	Tags        []string   `json:"tags,omitempty"`
}

type CreateResult struct {
	client.ResultMetadata
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UpdateResult struct {
	client.ResultMetadata
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type GetResult struct {
	client.ResultMetadata
	Service Service `json:"data"`
}

type ListResult struct {
	client.ResultMetadata
	Services []Service `json:"data"`
	Paging   Paging    `json:"paging"`
}

type Paging struct {
	Next  string `json:"next"`
	First string `json:"first"`
	Last  string `json:"last"`
}
