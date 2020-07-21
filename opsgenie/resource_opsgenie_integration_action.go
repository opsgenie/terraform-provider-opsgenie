package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/integration"
)

func resourceOpsgenieIntegrationAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieIntegrationActionCreate,
		Read:   handleNonExistentResource(resourceOpsgenieIntegrationActionRead),
		Update: resourceOpsgenieIntegrationActionCreate,
		Delete: resourceOpsgenieIntegrationActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"integration_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"email", "api"}, false),
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"create", "close", "acknowledge", "addNote"}, false),
						},
						"filter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"match-all", "match-any-condition", "match-all-conditions"}, false),
									},
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Required: true,
												},
												"key": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"not": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"operation": {
													Type:     schema.TypeString,
													Required: true,
												},
												"expected_value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"order": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"user": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{{user}}",
						},
						"note": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{{note}}",
						},
						"alias": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{{alias}}",
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{{source}}",
						},
						"message": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{{message}}",
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{{description}}",
						},
						"entity": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{{entity}}",
						},
						"append_attachments": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"ignore_alert_actions_from_payload": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"ignore_responders_from_payload": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"ignore_teams_from_payload": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"ignore_tags_from_payload": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"ignore_extra_properties_from_payload": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"alert_actions": {
							Type:     schema.TypeList,
							Optional: true,
						},
						"responders": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"username": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func expandOpsgenieFilter(input []interface{}) *og.Filter {
	filter := og.Filter{}
	for _, r := range input {
		inputMap := r.(map[string]interface{})
		conditions := expandOpsgenieConditions(inputMap["conditions"].([]interface{}))
		filter.Conditions = conditions
		filter.ConditionMatchType = inputMap["type"].(og.ConditionMatchType)
	}
	return &filter
}

func expandOpsgenieActionResponders(input []interface{}) []integration.Responder {

	responders := make([]integration.Responder, 0, len(input))

	if input == nil {
		return responders
	}

	for _, v := range input {
		config := v.(map[string]interface{})
		responderID := config["id"].(string)
		responder := integration.Responder{
			Type: integration.ResponderType(config["type"].(string)),
			Id:   responderID,
		}
		responders = append(responders, responder)
	}
	return responders
}

func expandOpsgenieActions(input []interface{}) *map[integration.ActionType]([]integration.IntegrationAction) {

	actionMap := map[integration.ActionType][]integration.IntegrationAction{
		integration.Create:      make([]integration.IntegrationAction, 0),
		integration.Close:       make([]integration.IntegrationAction, 0),
		integration.Acknowledge: make([]integration.IntegrationAction, 0),
		integration.AddNote:     make([]integration.IntegrationAction, 0),
	}

	for _, v := range input {
		inputMap := v.(map[string]interface{})
		action := integration.IntegrationAction{}

		appendAttachment := inputMap["append_attachments"].(bool)
		ignoreAlertActionsFromPayload := inputMap["ignore_alert_actions_from_payload"].(bool)
		ignoreRespondersFromPayload := inputMap["ignore_responders_from_payload"].(bool)
		ignoreTagsFromPayload := inputMap["ignore_tags_from_payload"].(bool)
		ignoreExtraPropertiesFromPayload := inputMap["ignore_extra_properties_from_payload"].(bool)

		action.Type = inputMap["Type"].(integration.ActionType)
		action.Name = inputMap["Name"].(string)
		action.Alias = inputMap["Alias"].(string)
		action.Order = inputMap["Order"].(int)
		action.User = inputMap["User"].(string)
		action.Note = inputMap["Note"].(string)
		action.Source = inputMap["Source"].(string)
		action.Message = inputMap["Message"].(string)
		action.Description = inputMap["Description"].(string)
		action.Entity = inputMap["Entity"].(string)
		action.AlertActions = inputMap["AlertActions"].([]string)
		action.Tags = inputMap["Tags"].([]string)
		action.ExtraProperties = inputMap["ExtraProperties"].(map[string]string)
		action.AppendAttachments = &appendAttachment
		action.IgnoreTagsFromPayload = &ignoreTagsFromPayload
		action.IgnoreRespondersFromPayload = &ignoreRespondersFromPayload
		action.IgnoreAlertActionsFromPayload = &ignoreAlertActionsFromPayload
		action.IgnoreExtraPropertiesFromPayload = &ignoreExtraPropertiesFromPayload

		action.Filter = expandOpsgenieFilter(inputMap["Filter"].([]interface{}))
		action.Responders = expandOpsgenieActionResponders(inputMap["Responders"].([]interface{}))

		actions := actionMap[action.Type]
		actions = append(actions, action)
	}
	return &actionMap
}

func resourceOpsgenieIntegrationActionCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	integrationId := d.Get("integration_id").(string)
	//integrationType := d.Get("type").(string)
	integrationActions := d.Get("actions").([]interface{})

	actionMap := *expandOpsgenieActions(integrationActions)

	updateRequest := &integration.UpdateAllIntegrationActionsRequest{
		Id:          integrationId,
		Create:      actionMap[integration.Create],
		Close:       actionMap[integration.Close],
		Acknowledge: actionMap[integration.Acknowledge],
		AddNote:     actionMap[integration.AddNote],
	}

	//log.Printf("[INFO] Creating OpsGenie api integration '%s'", name)

	result, err := client.UpdateAllActions(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	d.Set("integration_id", result.Parent.Id)
	return resourceOpsgenieIntegrationActionRead(d, meta)
}

func flattenOpsgenieIntegrationActions(result *integration.ActionsResult) []integration.IntegrationAction {

	actions := make([]integration.IntegrationAction, 0)

	for _, action := range result.Create {
		//actions = append(actions, integration.IntegrationAction(action))
		actions = append(actions, integration.IntegrationAction{
			Type:  integration.Create,
			Name:  action.Name,
			Alias: action.Alias,
			Order: action.Order,
			User:  action.User,
			Note:  action.Note,
			//Filter:                           action.Filter,
			Source:                           action.Source,
			Message:                          action.Message,
			Description:                      action.Description,
			Entity:                           action.Entity,
			AppendAttachments:                &action.AppendAttachments,
			AlertActions:                     action.AlertActions,
			IgnoreAlertActionsFromPayload:    &action.IgnoreAlertActionsFromPayload,
			IgnoreRespondersFromPayload:      &action.IgnoreRespondersFromPayload,
			IgnoreTagsFromPayload:            &action.IgnoreTagsFromPayload,
			IgnoreExtraPropertiesFromPayload: &action.IgnoreExtraPropertiesFromPayload,
			Responders:                       action.Responders,
			Tags:                             action.Tags,
			ExtraProperties:                  action.ExtraProperties,
		})
	}

	return actions
}

func resourceOpsgenieIntegrationActionRead(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	result, err := client.GetActions(context.Background(), &integration.GetIntegrationActionsRequest{
		BaseRequest: ogClient.BaseRequest{},
		Id:          d.Get("integration_id").(string),
	})
	if err != nil {
		return err
	}

	d.Set("integration_id", result.Parent.Id)
	d.Set("integration_type", result.Parent.Type)
	actions := flattenOpsgenieIntegrationActions(result)
	d.Set("actions", actions)

	return nil
}

func resourceOpsgenieIntegrationActionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie api integration '%s'", d.Get("name").(string))
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &integration.UpdateAllIntegrationActionsRequest{
		Id:          d.Get("integration_id").(string),
		Create:      nil,
		Close:       nil,
		Acknowledge: nil,
		AddNote:     nil,
	}

	_, err = client.UpdateAllActions(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}
