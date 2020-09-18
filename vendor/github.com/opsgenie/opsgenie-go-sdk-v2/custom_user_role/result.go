package custom_user_role

import "github.com/opsgenie/opsgenie-go-sdk-v2/client"

type CreateResult struct {
	client.ResultMetadata
	Result string `json:"result"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

type GetResult struct {
	client.ResultMetadata
	Id               string       `json:"id"`
	Name             string       `json:"name"`
	ExtendedRole     ExtendedRole `json:"extendedRole"`
	GrantedRights    []string     `json:"grantedRights"`
	DisallowedRights []string     `json:"disallowedRights"`
}

type UpdateResult struct {
	client.ResultMetadata
	Result string `json:"result"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type ListResult struct {
	client.ResultMetadata
	CustomUserRoles []CustomUserRole `json:"data"`
}

type CustomUserRole struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
