package opsgenie

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"

	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
)

func resourceOpsGenieAlertPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpsGenieAlertPolicyCreate,
		ReadContext:   resourceOpsGenieAlertPolicyRead,
		Update:        resourceOpsGenieAlertPolicyUpdate,
		Delete:        resourceOpsGenieAlertPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) == 1 {
					//if its not team_id/policy_id it can be global policy
					return []*schema.ResourceData{d}, nil
				} else if len(idParts) == 2 && idParts[0] != "" && idParts[1] != "" {
					d.Set("team_id", idParts[0])
					d.SetId(idParts[1])
					return []*schema.ResourceData{d}, nil
				} else {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected team_id/notification_policy_id", d.Id())
				}
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"team_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"policy_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "match-all",
							ValidateFunc: validation.StringInSlice([]string{"match-all", "match-any-condition", "match-all-conditions"}, false),
						},
						"conditions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"message", "alias", "description", "source", "entity", "tags",
											"actions", "details", "extra-properties", "responders", "teams", "priority",
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
										Type:        schema.TypeString,
										Optional:    true,
										Description: "User defined value that will be compared with alert field according to the operation. Default value is empty string",
									},
									"order": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Order of the condition in conditions list",
									},
								},
							},
						},
					},
				},
			},
			"time_restriction": timeRestrictionSchema(),
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"continue_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{{alias}}",
			},
			"alert_description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{{description}}",
			},
			"entity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{{entity}}",
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{{source}}",
			},
			"ignore_original_actions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ignore_original_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"actions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ignore_original_responders": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"responders": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "team", "escalation", "schedule"}, false),
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
			"ignore_original_tags": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"priority": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"P1", "P2", "P3", "P4", "P5"}, false),
			},
		},
	}
}

func resourceOpsGenieAlertPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return diag.FromErr(err)
	}

	message := d.Get("message").(string)
	continue_policy := d.Get("continue_policy").(bool)
	alias := d.Get("alias").(string)
	alert_description := d.Get("alert_description").(string)
	entity := d.Get("entity").(string)
	source := d.Get("source").(string)
	ignore_original_actions := d.Get("ignore_original_actions").(bool)
	ignore_original_details := d.Get("ignore_original_details").(bool)
	ignore_original_responders := d.Get("ignore_original_responders").(bool)
	ignore_original_tags := d.Get("ignore_original_tags").(bool)
	priority := d.Get("priority").(string)

	createRequest := &policy.CreateAlertPolicyRequest{
		MainFields:               *expandOpsGenieAlertPolicyRequestMainFields(d),
		Message:                  message,
		Continue:                 &continue_policy,
		Alias:                    alias,
		AlertDescription:         alert_description,
		Entity:                   entity,
		Source:                   source,
		IgnoreOriginalDetails:    &ignore_original_details,
		IgnoreOriginalActions:    &ignore_original_actions,
		IgnoreOriginalResponders: &ignore_original_responders,
		IgnoreOriginalTags:       &ignore_original_tags,
		Priority:                 alert.Priority(priority),
		Actions:                  flattenOpsgenieAlertPolicyActions(d),
		Tags:                     flattenTags(d, "tags"),
	}

	if d.Get("responders").(*schema.Set).Len() > 0 {
		createRequest.Responders = expandOpsGenieAlertPolicyResponders(d)
	}

	log.Printf("[INFO] Creating Alert Policy '%s'", d.Get("name").(string))
	result, err := client.CreateAlertPolicy(context.Background(), createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Id)

	return resourceOpsGenieAlertPolicyRead(ctx, d, meta)
}

func resourceOpsGenieAlertPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	log.Printf("[INFO] Reading OpsGenie Alert Policy '%s'", name)

	policyRes := &policy.GetAlertPolicyResult{}
	if d.Get("team_id").(string) == "" {
		policyRes, err = client.GetAlertPolicy(context.Background(), &policy.GetAlertPolicyRequest{
			Id: d.Id(),
		})
	} else {
		policyRes, err = client.GetAlertPolicy(context.Background(), &policy.GetAlertPolicyRequest{
			Id:     d.Id(),
			TeamId: d.Get("team_id").(string),
		})
	}
	if err != nil {
		x := err.(*ogClient.ApiError)
		if x.StatusCode == 404 {
			tflog.Warn(ctx, fmt.Sprintf("Removing Alert Policy because it's gone %s", name))
			d.SetId("")
			return nil
		} else if x.StatusCode >= 400 {
			tflog.Error(ctx, fmt.Sprintf("%d: %s", x.StatusCode, x.Message))
			d.SetId("")
			return nil
		}
	}
	d.Set("name", policyRes.Name)
	d.Set("enabled", policyRes.Enabled)
	d.Set("policy_description", policyRes.PolicyDescription)

	d.Set("message", policyRes.Message)
	d.Set("continue_policy", policyRes.Continue)
	d.Set("alias", policyRes.Alias)
	d.Set("alert_description", policyRes.AlertDescription)
	d.Set("entity", policyRes.Entity)
	d.Set("source", policyRes.Source)
	d.Set("ignore_original_actions", policyRes.IgnoreOriginalActions)
	d.Set("ignore_original_details", policyRes.IgnoreOriginalDetails)
	d.Set("ignore_original_responders", policyRes.IgnoreOriginalResponders)
	d.Set("ignore_original_tags", policyRes.IgnoreOriginalTags)
	d.Set("actions", policyRes.Actions)
	d.Set("tags", policyRes.Tags)

	if policyRes.Responders != nil {
		d.Set("responders", flattenOpsGenieAlertPolicyResponders(policyRes.Responders))
	} else {
		d.Set("responders", nil)
	}

	if policyRes.MainFields.Filter != nil {
		d.Set("filter", flattenOpsGenieAlertPolicyFilter(policyRes.MainFields.Filter))
	} else {
		d.Set("filter", nil)
	}

	if policyRes.MainFields.TimeRestriction != nil {
		log.Printf("[DEBUG] 'policy.MainFields.TimeRestriction' is not 'nil'.")
		d.Set("time_restriction", flattenOpsgenieTimeRestriction(policyRes.MainFields.TimeRestriction))
	} else {
		log.Printf("[DEBUG] 'policy.MainFields.TimeRestriction' is 'nil'.")
		d.Set("time_restriction", nil)
	}

	return nil
}

func resourceOpsGenieAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	message := d.Get("message").(string)
	continue_policy := d.Get("continue_policy").(bool)
	alias := d.Get("alias").(string)
	alert_description := d.Get("alert_description").(string)
	entity := d.Get("entity").(string)
	source := d.Get("source").(string)
	ignore_original_actions := d.Get("ignore_original_actions").(bool)
	ignore_original_details := d.Get("ignore_original_details").(bool)
	ignore_original_responders := d.Get("ignore_original_responders").(bool)
	ignore_original_tags := d.Get("ignore_original_tags").(bool)
	priority := d.Get("priority").(string)

	updateRequest := &policy.UpdateAlertPolicyRequest{
		Id:                       d.Id(),
		MainFields:               *expandOpsGenieAlertPolicyRequestMainFields(d),
		Message:                  message,
		Continue:                 &continue_policy,
		Alias:                    alias,
		AlertDescription:         alert_description,
		Entity:                   entity,
		Source:                   source,
		IgnoreOriginalDetails:    &ignore_original_details,
		IgnoreOriginalActions:    &ignore_original_actions,
		IgnoreOriginalResponders: &ignore_original_responders,
		IgnoreOriginalTags:       &ignore_original_tags,
		Priority:                 alert.Priority(priority),
		Actions:                  flattenOpsgenieAlertPolicyActions(d),
		Tags:                     flattenTags(d, "tags"),
	}

	if d.Get("responders").(*schema.Set).Len() > 0 {
		updateRequest.Responders = expandOpsGenieAlertPolicyResponders(d)
	}

	log.Printf("[INFO] Updating Alert Policy '%s'", d.Get("name").(string))
	_, err = client.UpdateAlertPolicy(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie Alert Policy '%s'", d.Get("name").(string))
	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	deleteRequest := &policy.DeletePolicyRequest{}
	if d.Get("team_id").(string) == "" {
		deleteRequest = &policy.DeletePolicyRequest{
			Id:   d.Id(),
			Type: "alert",
		}
	} else {
		deleteRequest = &policy.DeletePolicyRequest{
			Id:     d.Id(),
			TeamId: d.Get("team_id").(string),
			Type:   "alert",
		}

	}

	_, err = client.DeletePolicy(context.Background(), deleteRequest)
	if err != nil {
		return err
	}
	return nil
}

func expandOpsGenieAlertPolicyRequestMainFields(d *schema.ResourceData) *policy.MainFields {
	enabled := d.Get("enabled").(bool)
	fields := policy.MainFields{
		Name:              d.Get("name").(string),
		Enabled:           &enabled,
		PolicyDescription: d.Get("policy_description").(string),
	}
	if d.Get("team_id").(string) != "" {
		fields.TeamId = d.Get("team_id").(string)
	}
	if len(d.Get("filter").([]interface{})) > 0 {
		fields.Filter = expandOpsGenieAlertPolicyFilter(d.Get("filter").([]interface{}))
	}
	if len(d.Get("time_restriction").([]interface{})) > 0 {
		fields.TimeRestriction = expandOpsGenieTimeRestriction(d.Get("time_restriction").([]interface{}))
	}
	return &fields
}

func expandOpsGenieAlertPolicyResponders(d *schema.ResourceData) *[]alert.Responder {
	input := d.Get("responders").(*schema.Set)
	responders := make([]alert.Responder, 0, input.Len())

	if input == nil {
		return &responders
	}

	for _, v := range input.List() {
		config := v.(map[string]interface{})
		responderID := config["id"].(string)
		name := config["name"].(string)
		username := config["username"].(string)

		responder := alert.Responder{
			Type:     alert.ResponderType(config["type"].(string)),
			Id:       responderID,
			Name:     name,
			Username: username,
		}

		responders = append(responders, responder)
	}

	return &responders
}

func expandOpsGenieAlertPolicyFilter(input []interface{}) *og.Filter {
	filter := og.Filter{}

	if input == nil {
		return &filter
	}

	for _, v := range input {
		config := v.(map[string]interface{})
		filter.ConditionMatchType = og.ConditionMatchType(config["type"].(string))
		filter.Conditions = expandOpsGenieAlertPolicyFilterConditions(config["conditions"].(*schema.Set))
	}
	return &filter
}

func expandOpsGenieAlertPolicyFilterConditions(input *schema.Set) []og.Condition {
	conditions := make([]og.Condition, 0, input.Len())
	condition := og.Condition{}
	if input == nil {
		return conditions
	}

	for _, v := range input.List() {
		config := v.(map[string]interface{})
		not_value := config["not"].(bool)
		order := config["order"].(int)
		condition.Field = og.ConditionFieldType(config["field"].(string))
		condition.Operation = og.ConditionOperation(config["operation"].(string))
		condition.Key = config["key"].(string)
		condition.IsNot = &not_value
		condition.ExpectedValue = config["expected_value"].(string)
		condition.Order = &order
		conditions = append(conditions, condition)
	}
	return conditions
}

func flattenOpsGenieAlertPolicyResponders(input *[]alert.Responder) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, len(*input))
	for _, v := range *input {
		element := make(map[string]interface{})
		element["name"] = v.Name
		element["id"] = v.Id
		element["username"] = v.Username
		element["type"] = v.Type
		output = append(output, element)
	}

	return output
}

func flattenOpsGenieAlertPolicyFilter(input *og.Filter) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, 1)
	element := make(map[string]interface{})
	if input.Conditions != nil {
		element["conditions"] = flattenOpsGenieAlertPolicyFilterConditions(input.Conditions)
	}
	element["type"] = input.ConditionMatchType
	output = append(output, element)

	return output
}

func flattenOpsGenieAlertPolicyFilterConditions(input []og.Condition) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, len(input))
	for _, v := range input {
		element := make(map[string]interface{})
		element["field"] = v.Field
		element["operation"] = v.Operation
		element["key"] = v.Key
		element["not"] = v.IsNot
		element["expected_value"] = v.ExpectedValue
		element["order"] = v.Order
		output = append(output, element)
	}

	return output
}

func flattenOpsgenieAlertPolicyActions(d *schema.ResourceData) []string {
	input := d.Get("actions").(*schema.Set)
	actions := make([]string, len(input.List()))

	if input == nil {
		return actions
	}

	for k, v := range input.List() {
		actions[k] = v.(string)
	}

	return actions
}

func flattenOpsgenieAlertPolicyDetailsCreate(d *schema.ResourceData) []string {
	input := d.Get("details").(*schema.Set)
	details := make([]string, len(input.List()))

	if input == nil {
		return details
	}

	for k, v := range input.List() {
		details[k] = v.(string)
	}

	return details
}

func flattenOpsgenieAlertPolicyDetailsUpdate(d *schema.ResourceData) map[string]interface{} {
	input := d.Get("details").(*schema.Set)
	details := make(map[string]interface{}, len(input.List()))

	if input == nil {
		return details
	}

	for k, v := range input.List() {
		index := strconv.Itoa(k)
		details[index] = v
	}

	return details
}
