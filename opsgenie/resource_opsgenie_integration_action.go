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
		Update: resourceOpsgenieIntegrationActionUpdate,
		Delete: resourceOpsgenieIntegrationActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"integration_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "create",
						},
						"order": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
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
							Default:  true,
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"responders": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id": {
										Type:     schema.TypeString,
										Required: true,
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"close": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "close",
						},
						"order": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
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
					},
				},
			},
			"acknowledge": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "acknowledge",
						},
						"order": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
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
					},
				},
			},
			"add_note": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "AddNote",
						},
						"order": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
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
					},
				},
			},
		},
	}
}

func convertInterfaceSliceToString(input []interface{}) []string {
	result := make([]string, 0)
	for _, item := range input {
		result = append(result, item.(string))
	}
	return result
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

func expandOpsgenieFilter(input []interface{}) integration.Filter {
	filter := integration.Filter{}
	for _, r := range input {
		inputMap := r.(map[string]interface{})
		conditions := expandOpsgenieConditions(inputMap["conditions"].([]interface{}))
		filter.Conditions = conditions
		filter.ConditionMatchType = og.ConditionMatchType(inputMap["type"].(string))
	}
	return filter
}

func expandOpsgenieIntegrationActions(input interface{}) []integration.IntegrationAction {

	actions := make([]integration.IntegrationAction, 0)

	if input == nil {
		return actions
	}

	for _, v := range input.([]interface{}) {
		inputMap := v.(map[string]interface{})
		action := integration.IntegrationAction{}

		action.Type = integration.ActionType(inputMap["type"].(string))
		action.Name = inputMap["name"].(string)
		action.Alias = inputMap["alias"].(string)
		action.Order = inputMap["order"].(int)
		action.User = inputMap["user"].(string)
		action.Note = inputMap["note"].(string)
		filters := expandOpsgenieFilter(inputMap["filter"].([]interface{}))
		action.Filter = &filters

		if action.Type == integration.Create {
			action.Source = inputMap["source"].(string)
			action.Message = inputMap["message"].(string)
			action.Description = inputMap["description"].(string)
			action.Entity = inputMap["entity"].(string)
			action.AlertActions = convertInterfaceSliceToString(inputMap["alert_actions"].([]interface{}))
			action.Tags = convertInterfaceSliceToString(inputMap["tags"].([]interface{}))
			if extraProperties := inputMap["extra_properties"]; extraProperties != nil {
				action.ExtraProperties = extraProperties.(map[string]string)
			}

			appendAttachment := inputMap["append_attachments"].(bool)
			ignoreAlertActionsFromPayload := inputMap["ignore_alert_actions_from_payload"].(bool)
			ignoreRespondersFromPayload := inputMap["ignore_responders_from_payload"].(bool)
			ignoreTagsFromPayload := inputMap["ignore_tags_from_payload"].(bool)
			ignoreExtraPropertiesFromPayload := inputMap["ignore_extra_properties_from_payload"].(bool)

			action.AppendAttachments = &appendAttachment
			action.IgnoreTagsFromPayload = &ignoreTagsFromPayload
			action.IgnoreRespondersFromPayload = &ignoreRespondersFromPayload
			action.IgnoreAlertActionsFromPayload = &ignoreAlertActionsFromPayload
			action.IgnoreExtraPropertiesFromPayload = &ignoreExtraPropertiesFromPayload
			action.Responders = expandOpsgenieActionResponders(inputMap["responders"].([]interface{}))
		}

		actions = append(actions, action)
	}
	return actions
}

func flattenOpsgenieFilter(input integration.FilterResult) []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, 1)
	out := make(map[string]interface{})
	out["type"] = input.ConditionMatchType
	conditions := make([]map[string]interface{}, 0, len(input.Conditions))
	for _, r := range input.Conditions {
		conditionMap := make(map[string]interface{})
		conditionMap["order"] = r.Order
		if r.Key != "" {
			conditionMap["key"] = r.Key
		}
		conditionMap["expected_value"] = r.ExpectedValue
		conditionMap["operation"] = r.Operation
		conditionMap["field"] = r.Field
		conditionMap["not"] = r.IsNot
		conditions = append(conditions, conditionMap)
	}
	out["conditions"] = conditions
	rules = append(rules, out)
	return rules
}

func flattenOpsgenieIntegrationCloseActions(input []integration.CloseAction) []map[string]interface{} {
	actions := make([]map[string]interface{}, 0)
	for _, action := range input {
		actionMap := make(map[string]interface{})
		actionMap["type"] = integration.Close
		actionMap["name"] = action.Name
		actionMap["alias"] = action.Alias
		actionMap["order"] = action.Order
		actionMap["note"] = action.Note
		actionMap["user"] = action.User
		actionMap["filter"] = flattenOpsgenieFilter(action.Filter)

		actions = append(actions, actionMap)
	}
	return actions
}

func flattenOpsgenieIntegrationAcknowledgeActions(input []integration.AcknowledgeAction) []map[string]interface{} {
	actions := make([]map[string]interface{}, 0)
	for _, action := range input {
		actionMap := make(map[string]interface{})
		actionMap["type"] = integration.Acknowledge
		actionMap["name"] = action.Name
		actionMap["alias"] = action.Alias
		actionMap["order"] = action.Order
		actionMap["note"] = action.Note
		actionMap["user"] = action.User
		actionMap["filter"] = flattenOpsgenieFilter(action.Filter)

		actions = append(actions, actionMap)
	}
	return actions
}

func flattenOpsgenieIntegrationAddNoteActions(input []integration.AddNoteAction) []map[string]interface{} {
	actions := make([]map[string]interface{}, 0)
	for _, action := range input {
		actionMap := make(map[string]interface{})
		actionMap["type"] = integration.AddNote
		actionMap["name"] = action.Name
		actionMap["alias"] = action.Alias
		actionMap["order"] = action.Order
		actionMap["note"] = action.Note
		actionMap["user"] = action.User
		actionMap["filter"] = flattenOpsgenieFilter(action.Filter)

		actions = append(actions, actionMap)
	}
	return actions
}

func flattenOpsgenieIntegrationCreateActions(input []integration.CreateAction) []map[string]interface{} {

	actions := make([]map[string]interface{}, 0)
	for _, action := range input {
		actionMap := make(map[string]interface{})
		actionMap["type"] = integration.Create
		actionMap["name"] = action.Name
		actionMap["alias"] = action.Alias
		actionMap["order"] = action.Order
		actionMap["note"] = action.Note
		actionMap["user"] = action.User
		actionMap["filter"] = flattenOpsgenieFilter(action.Filter)
		actionMap["source"] = action.Source
		actionMap["message"] = action.Message
		actionMap["description"] = action.Description
		actionMap["entity"] = action.Entity
		actionMap["append_attachments"] = action.AppendAttachments
		actionMap["alert_actions"] = action.AlertActions
		actionMap["ignore_alert_actions_from_payload"] = action.IgnoreAlertActionsFromPayload
		actionMap["ignore_responders_from_payload"] = action.IgnoreRespondersFromPayload
		actionMap["ignore_tags_from_payload"] = action.IgnoreTagsFromPayload
		actionMap["ignore_extra_properties_from_payload"] = action.IgnoreExtraPropertiesFromPayload

		responders := make([]map[string]string, 0)
		for _, responder := range action.Responders {
			responders = append(responders, map[string]string{
				"type":     string(responder.Type),
				"name":     responder.Name,
				"id":       responder.Id,
				"username": responder.Username,
			})
		}
		actionMap["responders"] = responders
		actionMap["tags"] = action.Tags
		actionMap["extra_properties"] = action.ExtraProperties

		actions = append(actions, actionMap)
	}
	return actions
}

func resourceOpsgenieIntegrationActionCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	integrationId := d.Get("integration_id").(string)
	updateRequest := &integration.UpdateAllIntegrationActionsRequest{
		Id:          integrationId,
		Create:      expandOpsgenieIntegrationActions(d.Get("create")),
		Close:       expandOpsgenieIntegrationActions(d.Get("close")),
		Acknowledge: expandOpsgenieIntegrationActions(d.Get("acknowledge")),
		AddNote:     expandOpsgenieIntegrationActions(d.Get("add_note")),
	}

	log.Printf("[INFO] Creating OpsGenie integration actions for '%s'", integrationId)
	result, err := client.UpdateAllActions(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Parent.Id)
	d.Set("integration_id", result.Parent.Id)

	return resourceOpsgenieIntegrationActionRead(d, meta)
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

	d.SetId(result.Parent.Id)
	d.Set("integration_id", result.Parent.Id)
	d.Set("create", flattenOpsgenieIntegrationCreateActions(result.Create))
	d.Set("close", flattenOpsgenieIntegrationCloseActions(result.Close))
	d.Set("acknowledge", flattenOpsgenieIntegrationAcknowledgeActions(result.Acknowledge))
	d.Set("add_note", flattenOpsgenieIntegrationAddNoteActions(result.AddNote))

	return nil
}

func resourceOpsgenieIntegrationActionUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceOpsgenieIntegrationActionCreate(d, meta)
}

func resourceOpsgenieIntegrationActionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie api integration actions for '%s'", d.Get("integration_id").(string))
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	deleteRequest := &integration.UpdateAllIntegrationActionsRequest{
		Id:          d.Get("integration_id").(string),
		Create:      []integration.IntegrationAction{},
		Close:       []integration.IntegrationAction{},
		Acknowledge: []integration.IntegrationAction{},
		AddNote:     []integration.IntegrationAction{},
	}

	_, err = client.UpdateAllActions(context.Background(), deleteRequest)
	if err != nil {
		apiError := err.(*ogClient.ApiError)
		if apiError.StatusCode != 404 {
			return err
		}
	}

	return nil
}
