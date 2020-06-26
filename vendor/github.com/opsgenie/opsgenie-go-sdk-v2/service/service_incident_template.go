package service

import (
	"context"
)

func (c *Client) CreateIncidentTemplate(context context.Context, request *CreateIncidentTemplateRequest) (*CreateIncidentTemplateResult, error) {
	result := &CreateIncidentTemplateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetIncidentTemplates(context context.Context, request *GetIncidentTemplatesRequest) (*GetIncidentTemplatesResult, error) {
	result := &GetIncidentTemplatesResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteIncidentTemplate(context context.Context, request *DeleteIncidentTemplateRequest) (*DeleteIncidentTemplateResult, error) {
	result := &DeleteIncidentTemplateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateIncidentTemplate(context context.Context, request *UpdateIncidentTemplateRequest) (*UpdateIncidentTemplateResult, error) {
	result := &UpdateIncidentTemplateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
