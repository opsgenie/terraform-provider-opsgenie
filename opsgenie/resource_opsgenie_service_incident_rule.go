package opsgenie

import (
	"context"
	"fmt"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/service"
)

func resourceOpsGenieServiceIncidentRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieServiceIncidentRuleCreate,
		Read:   handleNonExistentResource(resourceOpsGenieServiceIncidentRuleRead),
		Update: resourceOpsGenieServiceIncidentRuleUpdate,
		Delete: resourceOpsGenieServiceIncidentRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected service_id/service_incident_rule_id", d.Id())
				}
				d.Set("service_id", idParts[0])
				d.SetId(idParts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 130),
			},
			"incident_rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition_match_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "match-all",
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
										ValidateFunc: validation.StringInSlice([]string{
											"message", "description", "tags",
											"extra-properties", "recipients", "teams", "priority",
										}, false),
									},
									"operation": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"matches", "contains", "starts-with", "ends-with", "equals", "contains-key",
											"contains-value", "greater-than", "less-than", "is-empty", "equals-ignore-whitespace",
										}, false),
									},
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If 'field' is set as 'extra-properties', key could be used for key-value pair",
									},
									"not": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicates behaviour of the given operation. Default value is false",
										Default:     false,
									},
									"expected_value": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "User defined value that will be compared with alert field according to the operation. Default value is empty string",
										ValidateFunc: validation.StringLenBetween(1, 15000),
									},
								},
							},
						},
						"incident_properties": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 130),
									},
									"tags": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"details": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "{{description}}",
										ValidateFunc: validation.StringLenBetween(1, 10000),
									},
									"priority": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"P1", "P2", "P3", "P4", "P5"}, false),
									},
									"stakeholder_properties": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  true,
												},
												"message": {
													Type:     schema.TypeString,
													Required: true,
												},
												"description": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringLenBetween(1, 15000),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceOpsGenieServiceIncidentRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	service_id := d.Get("service_id").(string)
	createRequest := &service.CreateIncidentRuleRequest{
		ServiceId: service_id,
	}

	incident_rule := d.Get("incident_rule").([]interface{})
	for _, v := range incident_rule {
		config := v.(map[string]interface{})
		createRequest.ConditionMatchType = og.ConditionMatchType(config["condition_match_type"].(string))
		createRequest.Conditions = expandOpsGenieServiceIncidentRuleConditions(config["conditions"].([]interface{}))
		createRequest.IncidentProperties = expandOpsGenieServiceIncidentRuleIncidentProperties(config["incident_properties"].([]interface{}))
	}

	log.Printf("[INFO] Creating OpsGenie Service Incident Rule for service '%s'", d.Get("service_id").(string))
	result, err := client.CreateIncidentRule(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsGenieServiceIncidentRuleRead(d, meta)
}

func resourceOpsGenieServiceIncidentRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	service_id := d.Get("service_id").(string)
	incident_rule_id := d.Id()

	log.Printf("[INFO] Reading OpsGenie Service Incident Rule for service: '%s' for rule ID: '%s'", service_id, incident_rule_id)

	incident_rule_res, err := client.GetIncidentRules(context.Background(), &service.GetIncidentRulesRequest{
		ServiceId: service_id,
	})
	if err != nil {
		x := err.(*ogClient.ApiError)
		if x.StatusCode == 404 {
			log.Printf("[WARN] Removing Service Incident Rule because it's gone %s", service_id)
			d.SetId("")
			return nil
		}
	}

	for _, v := range incident_rule_res.IncidentRule {
		if v.Id == incident_rule_id {
			d.Set("service_id", service_id)
			d.Set("incident_rule", flattenOpsGenieServiceIncidentRules(v))
		}
	}
	return nil
}

func resourceOpsGenieServiceIncidentRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	service_id := d.Get("service_id").(string)
	incident_rule_id := d.Id()

	updateRequest := &service.UpdateIncidentRuleRequest{
		ServiceId:      service_id,
		IncidentRuleId: incident_rule_id,
	}

	incident_rule := d.Get("incident_rule").([]interface{})
	for _, v := range incident_rule {
		config := v.(map[string]interface{})
		updateRequest.ConditionMatchType = og.ConditionMatchType(config["condition_match_type"].(string))
		updateRequest.Conditions = expandOpsGenieServiceIncidentRuleConditions(config["conditions"].([]interface{}))
		updateRequest.IncidentProperties = expandOpsGenieServiceIncidentRuleIncidentProperties(config["incident_properties"].([]interface{}))
	}

	log.Printf("[INFO] Updating Service Incident Rule for service: '%s' for rule ID: '%s'", service_id, incident_rule_id)
	_, err = client.UpdateIncidentRule(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieServiceIncidentRuleDelete(d *schema.ResourceData, meta interface{}) error {
	service_id := d.Get("service_id").(string)
	incident_rule_id := d.Id()

	log.Printf("[INFO] Deleting OpsGenie ervice Incident Rule for service: '%s' for rule ID: '%s'", service_id, incident_rule_id)
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &service.DeleteIncidentRuleRequest{
		ServiceId:      service_id,
		IncidentRuleId: incident_rule_id,
	}

	_, err = client.DeleteIncidentRule(context.Background(), deleteRequest)
	if err != nil {
		return err
	}
	return nil
}

func flattenOpsGenieServiceIncidentRules(input service.IncidentRuleResult) []map[string]interface{} {
	incident_rule := make(map[string]interface{})
	conditions_list := make([]map[string]interface{}, 0, len(input.Conditions))

	for _, v := range input.Conditions {
		conditions := flattenOpsGenieServiceIncidentRuleConditions(v)
		conditions_list = append(conditions_list, conditions)
	}
	incident_rule["conditions"] = conditions_list
	incident_rule["condition_match_type"] = input.ConditionMatchType
	incident_rule["incident_properties"] = flattenOpsGenieServiceIncidentRuleIncidentProperties(input.IncidentProperties)

	return []map[string]interface{}{incident_rule}

}

func expandOpsGenieServiceIncidentRuleConditions(input []interface{}) []og.Condition {
	conditions := make([]og.Condition, 0, len(input))
	if input == nil {
		return conditions
	}

	for _, v := range input {
		condition := og.Condition{}
		config := v.(map[string]interface{})
		not_value := config["not"].(bool)
		condition.Field = og.ConditionFieldType(config["field"].(string))
		condition.Operation = og.ConditionOperation(config["operation"].(string))
		condition.IsNot = &not_value
		condition.ExpectedValue = config["expected_value"].(string)
		if condition.Field == og.ExtraProperties {
			key := config["key"].(string)
			if key != "" {
				condition.Key = config["key"].(string)
			}
		}
		conditions = append(conditions, condition)
	}

	return conditions
}

func expandOpsGenieServiceIncidentRuleIncidentProperties(input []interface{}) service.IncidentProperties {
	incident_properties := service.IncidentProperties{}

	for _, v := range input {
		config := v.(map[string]interface{})
		incident_properties.Message = config["message"].(string)

		if config["tags"].(*schema.Set).Len() > 0 {
			incident_properties.Tags = flattenOpsgenieServiceIncidentRuleRequestTags(config["tags"].(*schema.Set))
		}
		if config["details"].(map[string]interface{}) != nil {
			incident_properties.Details = flattenOpsgenieServiceIncidentRuleRequestDetails(config["details"].(map[string]interface{}))
		}

		incident_properties.Description = config["description"].(string)
		incident_properties.Priority = alert.Priority(config["priority"].(string))
		incident_properties.StakeholderProperties = expandOpsGenieServiceIncidentRuleStakeholderProperties(config["stakeholder_properties"].([]interface{}))
	}
	return incident_properties
}

func expandOpsGenieServiceIncidentRuleStakeholderProperties(input []interface{}) service.StakeholderProperties {
	stakeholder_properties := service.StakeholderProperties{}
	if input == nil {
		return stakeholder_properties
	}

	for _, v := range input {
		config := v.(map[string]interface{})
		enable := config["enable"].(bool)
		stakeholder_properties.Enable = &enable
		stakeholder_properties.Message = config["message"].(string)
		stakeholder_properties.Description = config["description"].(string)
	}

	return stakeholder_properties
}

func flattenOpsgenieServiceIncidentRuleRequestTags(input *schema.Set) []string {
	tags := make([]string, len(input.List()))
	if input == nil {
		return tags
	}

	for k, v := range input.List() {
		tags[k] = v.(string)
	}
	return tags
}

func flattenOpsgenieServiceIncidentRuleRequestDetails(input map[string]interface{}) map[string]string {
	details := make(map[string]string)

	if input == nil {
		return details
	}

	for k, v := range input {
		details[k] = v.(string)
	}
	return details
}

func flattenOpsGenieServiceIncidentRuleConditions(input og.Condition) map[string]interface{} {
	condition := make(map[string]interface{})

	condition["field"] = input.Field
	condition["operation"] = input.Operation
	condition["not"] = *input.IsNot
	condition["expected_value"] = input.ExpectedValue
	if input.Key != "" {
		condition["key"] = input.Key
	}

	return condition
}

func flattenOpsGenieServiceIncidentRuleIncidentProperties(input service.IncidentProperties) []map[string]interface{} {
	incident_properties := make(map[string]interface{})
	if len(input.Tags) > 0 {
		incident_properties["tags"] = schema.NewSet(schema.HashString, convertStringSliceToInterfaceSlice(input.Tags))
	}
	if len(input.Details) > 0 {
		incident_properties["details"] = convertStringMapToInterfaceMap(input.Details)
	}
	if input.Description != "" {
		incident_properties["description"] = input.Description
	}
	incident_properties["message"] = input.Message
	incident_properties["priority"] = input.Priority
	incident_properties["stakeholder_properties"] = flattenOpsGenieServiceIncidentRuleStakeholderProperties(input.StakeholderProperties)
	return []map[string]interface{}{incident_properties}
}

func flattenOpsGenieServiceIncidentRuleStakeholderProperties(input service.StakeholderProperties) []map[string]interface{} {
	stakeholders_properties := make(map[string]interface{})

	stakeholders_properties["enable"] = *input.Enable
	stakeholders_properties["message"] = input.Message
	if input.Description != "" {
		stakeholders_properties["description"] = input.Description
	}
	return []map[string]interface{}{stakeholders_properties}
}
