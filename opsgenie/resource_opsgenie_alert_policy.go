package opsgenie

import (
	"context"
	"fmt"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"strconv"

	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
)

func resourceOpsGenieAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieAlertPolicyCreate,
		Read:   handleNonExistentResource(resourceOpsGenieAlertPolicyRead),
		Update: resourceOpsGenieAlertPolicyUpdate,
		Delete: resourceOpsGenieAlertPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected team_id/notification_policy_id", d.Id())
				}
				d.Set("team_id", idParts[0])
				d.SetId(idParts[1])
				return []*schema.ResourceData{d}, nil
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
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
											"message", "alias", "description", "source", "entity", "tags",
											"actions", "details", "extra-properties", "recipients", "teams", "priority",
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
			"time_restriction": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"time-of-day", "weekday-and-time-of-day"}, false),
						},
						"restrictions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_day": {
										Type:     schema.TypeString,
										Required: true,
									},
									"end_day": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start_hour": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"start_min": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"end_hour": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"end_min": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"restriction": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_hour": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"start_min": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"end_hour": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"end_min": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
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
			},
			"alert_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ignore_original_actions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				//Set: schema.HashString,
			},
			"ignore_original_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"details": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				//Set: schema.HashString,
			},
			"ignore_original_responders": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"responders": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "team"}, false),
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				//Set: schema.HashString,
			},
			"priority": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"P1", "P2", "P3", "P4", "P5"}, false),
			},
		},
	}
}

func resourceOpsGenieAlertPolicyCreate(d *schema.ResourceData, meta interface{}) error {
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

	createRequest := &policy.CreateAlertPolicyRequest{
		MainFields:               *expandOpsGenieAlertPolicyRequestMainFields(d),
		Message:                  message,
		Continue:                 &continue_policy,
		Alias:                    alias,
		AlertDescription:         alert_description,
		Entity:                   entity,
		Source:                   source,
		IgnoreOriginalDetails:    &ignore_original_actions,
		IgnoreOriginalActions:    &ignore_original_details,
		IgnoreOriginalResponders: &ignore_original_responders,
		IgnoreOriginalTags:       &ignore_original_tags,
		Priority:                 alert.Priority(priority),
	}

	if len(d.Get("responders").([]interface{})) > 0 {
		createRequest.Responders = expandOpsGenieAlertPolicyResponders(d.Get("Responders").([]interface{}))
	}

	if len(d.Get("actions").([]interface{})) > 0 {
		createRequest.Actions = flattenOpsgenieAlertPolicyActions(d)
	}

	if len(d.Get("details").([]interface{})) > 0 {
		createRequest.Details = flattenOpsgenieAlertPolicyDetailsCreate(d)
	}

	if len(d.Get("tags").([]interface{})) > 0 {
		createRequest.Tags = flattenOpsgenieAlertPolicyTags(d)
	}

	log.Printf("[INFO] Creating Alert Policy '%s'", d.Get("name").(string))
	result, err := client.CreateAlertPolicy(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsGenieAlertPolicyRead(d, meta)
}

func resourceOpsGenieAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	log.Printf("[INFO] Reading OpsGenie Alert Policy '%s'", name)

	policy, err := client.GetAlertPolicy(context.Background(), &policy.GetAlertPolicyRequest{
		Id:     d.Id(),
		TeamId: d.Get("team_id").(string),
	})
	if err != nil {
		x := err.(*ogClient.ApiError)
		if x.StatusCode == 404 {
			log.Printf("[WARN] Removing Alert Policy because it's gone %s", name)
			d.SetId("")
			return nil
		}
	}
	d.Set("name", policy.Name)
	d.Set("enabled", policy.Enabled)
	d.Set("policy_description", policy.PolicyDescription)

	d.Set("message", policy.Message)
	d.Set("continue_policy", policy.Continue)
	d.Set("alias", policy.Alias)
	d.Set("alert_description", policy.AlertDescription)
	d.Set("entity", policy.Entity)
	d.Set("source", policy.Source)
	d.Set("ignore_original_actions", policy.IgnoreOriginalActions)
	d.Set("actions", policy.Actions)
	d.Set("ignore_original_details", policy.IgnoreOriginalDetails)
	d.Set("details", policy.Details)
	d.Set("ignore_original_responders", policy.IgnoreOriginalResponders)
	d.Set("ignore_original_tags", policy.IgnoreOriginalTags)
	d.Set("tags", policy.Tags)

	if policy.Responders != nil {
		d.Set("responders", flattenOpsGenieAlertPolicyResponders(policy.Responders))
	} else {
		d.Set("responders", nil)
	}

	if policy.MainFields.Filter != nil {
		d.Set("filter", flattenOpsGenieAlertPolicyFilter(policy.MainFields.Filter))
	} else {
		d.Set("filter", nil)
	}

	if policy.MainFields.TimeRestriction != nil {
		log.Printf("[DEBUG] 'policy.MainFields.TimeRestriction' is not 'nil'.")
		d.Set("time_restriction", flattenOpsgenieAlertPolicyTimeRestriction(policy.MainFields.TimeRestriction))
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
	//actions := d.Get("actions").([]string)
	ignore_original_details := d.Get("ignore_original_details").(bool)
	//details := d.Get("details").(map[string]interface{})
	ignore_original_responders := d.Get("ignore_original_responders").(bool)
	ignore_original_tags := d.Get("ignore_original_tags").(bool)
	//tags := d.Get("tags").([]string)
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
		IgnoreOriginalDetails:    &ignore_original_actions,
		IgnoreOriginalActions:    &ignore_original_details,
		IgnoreOriginalResponders: &ignore_original_responders,
		IgnoreOriginalTags:       &ignore_original_tags,
		Priority:                 alert.Priority(priority),
	}

	if len(d.Get("responders").([]interface{})) > 0 {
		updateRequest.Responders = expandOpsGenieAlertPolicyResponders(d.Get("Responders").([]interface{}))
	}

	if len(d.Get("actions").([]interface{})) > 0 {
		updateRequest.Actions = flattenOpsgenieAlertPolicyActions(d)
	}

	if len(d.Get("details").([]interface{})) > 0 {
		updateRequest.Details = flattenOpsgenieAlertPolicyDetailsUpdate(d)
	}

	if len(d.Get("tags").([]interface{})) > 0 {
		updateRequest.Tags = flattenOpsgenieAlertPolicyTags(d)
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
	deleteRequest := &policy.DeletePolicyRequest{
		Id:     d.Id(),
		TeamId: d.Get("team_id").(string),
		Type:   "alert",
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
		TeamId:            d.Get("team_id").(string),
	}
	if len(d.Get("filter").([]interface{})) > 0 {
		fields.Filter = expandOpsGenieAlertPolicyFilter(d.Get("filter").([]interface{}))
	}
	if len(d.Get("time_restriction").([]interface{})) > 0 {
		fields.TimeRestriction = expandOpsGenieAlertPolicyTimeRestriction(d.Get("time_restriction").([]interface{}))
	}
	return &fields
}

func expandOpsGenieAlertPolicyResponders(input []interface{}) *[]alert.Responder {
	responders := make([]alert.Responder, 0, len(input))

	if input == nil {
		return &responders
	}

	for _, v := range input {
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
		filter.Conditions = expandOpsGenieAlertPolicyFilterConditions(config["conditions"].([]interface{}))
	}
	return &filter
}

func expandOpsGenieAlertPolicyFilterConditions(input []interface{}) []og.Condition {
	conditions := make([]og.Condition, 0, len(input))
	condition := og.Condition{}
	if input == nil {
		return conditions
	}

	for _, v := range input {
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

func expandOpsGenieAlertPolicyTimeRestriction(d []interface{}) *og.TimeRestriction {
	timeRestriction := og.TimeRestriction{}
	for _, v := range d {
		config := v.(map[string]interface{})
		timeRestriction.Type = og.RestrictionType(config["type"].(string))
		if len(config["restrictions"].([]interface{})) > 0 {
			restrictionList := make([]og.Restriction, 0, len(config["restrictions"].([]interface{})))
			for _, v := range config["restrictions"].([]interface{}) {
				config := v.(map[string]interface{})
				startHour := uint32(config["start_hour"].(int))
				startMin := uint32(config["start_min"].(int))
				endHour := uint32(config["end_hour"].(int))
				endMin := uint32(config["end_min"].(int))
				restriction := og.Restriction{
					StartDay:  og.Day(config["start_day"].(string)),
					StartHour: &startHour,
					StartMin:  &startMin,
					EndHour:   &endHour,
					EndDay:    og.Day(config["end_day"].(string)),
					EndMin:    &endMin,
				}
				restrictionList = append(restrictionList, restriction)
			}
			timeRestriction.RestrictionList = restrictionList
		} else {
			restriction := og.Restriction{}
			for _, v := range config["restriction"].([]interface{}) {
				config := v.(map[string]interface{})
				startHour := uint32(config["start_hour"].(int))
				startMin := uint32(config["start_min"].(int))
				endHour := uint32(config["end_hour"].(int))
				endMin := uint32(config["end_min"].(int))
				restriction = og.Restriction{
					StartHour: &startHour,
					StartMin:  &startMin,
					EndHour:   &endHour,
					EndMin:    &endMin,
				}
			}

			timeRestriction.Restriction = restriction
		}
	}
	return &timeRestriction
}

func flattenOpsGenieAlertPolicyDuration(input *policy.Duration) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, 1)
	element := make(map[string]interface{})
	element["time_amount"] = input.TimeAmount
	element["time_unit"] = input.TimeUnit
	output = append(output, element)
	return output
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

func flattenOpsgenieAlertPolicyTimeRestriction(input *og.TimeRestriction) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, 1)
	element := make(map[string]interface{})
	if len(input.RestrictionList) > 0 {
		restrictions := make([]map[string]interface{}, 0, len(input.RestrictionList))
		for _, r := range input.RestrictionList {
			restrictionMap := make(map[string]interface{})
			restrictionMap["start_min"] = r.StartMin
			restrictionMap["start_hour"] = r.StartHour
			restrictionMap["start_day"] = r.StartDay
			restrictionMap["end_min"] = r.EndMin
			restrictionMap["end_hour"] = r.EndHour
			restrictionMap["end_day"] = r.EndDay
			restrictions = append(restrictions, restrictionMap)
		}
		element["restrictions"] = restrictions
	} else {
		restriction := make([]map[string]interface{}, 0, 1)
		restrictionMap := make(map[string]interface{})
		restrictionMap["start_min"] = input.Restriction.StartMin
		restrictionMap["start_hour"] = input.Restriction.StartHour
		restrictionMap["end_min"] = input.Restriction.EndMin
		restrictionMap["end_hour"] = input.Restriction.EndHour
		restriction = append(restriction, restrictionMap)
		element["restriction"] = restriction
	}
	element["type"] = input.Type
	output = append(output, element)
	return output
}

func flattenOpsgenieAlertPolicyTags(d *schema.ResourceData) []string {
	input := d.Get("tags").(*schema.Set)
	tags := make([]string, len(input.List()))

	if input == nil {
		return tags
	}

	for k, v := range input.List() {
		tags[k] = v.(string)
	}

	return tags
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
