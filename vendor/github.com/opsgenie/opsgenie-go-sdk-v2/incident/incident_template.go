package incident

import "context"

func (c *Client) CreateIncidentTemplate(context context.Context, request *CreateIncidentTemplateRequest) (*CreateIncidentTemplateResult, error) {
	result := &CreateIncidentTemplateResult{}
	if err := c.client.Exec(context, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateIncidentTemplate(context context.Context, request *UpdateIncidentTemplateRequest) (*UpdateIncidentTemplateResult, error) {
	result := &UpdateIncidentTemplateResult{}
	if err := c.client.Exec(context, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteIncidentTemplate(context context.Context, request *DeleteIncidentTemplateRequest) (*DeleteIncidentTemplateResult, error) {
	result := &DeleteIncidentTemplateResult{}
	if err := c.client.Exec(context, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetIncidentTemplate(context context.Context, request *GetIncidentTemplateRequest) (*GetIncidentTemplateResult, error) {
	result := &GetIncidentTemplateResult{}
	if err := c.client.Exec(context, request, result); err != nil {
		return nil, err
	}
	return result, nil
}
