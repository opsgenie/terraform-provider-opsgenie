package opsgenie

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func resourceOpsGenieTeamRoutingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieTeamRoutingRuleCreate,
		Read:   handleNonExistentResource(resourceOpsGenieTeamRoutingRuleRead),
		Update: resourceOpsGenieTeamRoutingRuleUpdate,
		Delete: resourceOpsGenieTeamRoutingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected team_id/routing_rule_id", d.Id())
				}
				teamId := idParts[0]
				teamRoutingRuleId := idParts[1]
				d.Set("team_id", teamId)
				d.SetId(teamRoutingRuleId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"criteria": {
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
			"time_restriction": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
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
		},
	}
}

func resourceOpsGenieTeamRoutingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	teamId := d.Get("team_id").(string)
	order := d.Get("order").(int)
	timezone := d.Get("timezone").(string)
	timeRestriction := d.Get("time_restriction").([]interface{})
	criteria := d.Get("criteria").([]interface{})
	notify := d.Get("notify").([]interface{})

	expandedCriteria := expandOpsgenieCriteria(criteria)
	if err := validateOpsgenieCriteria(expandedCriteria); err != nil {
		return err
	}

	createRequest := &team.CreateRoutingRuleRequest{
		TeamIdentifierType:  team.Id,
		TeamIdentifierValue: teamId,
		Name:                name,
		Order:               &order,
		Timezone:            timezone,
		Criteria:            expandedCriteria,
		Notify:              expandOpsgenieNotify(notify),
	}

	if len(timeRestriction) > 0 {
		createRequest.TimeRestriction = expandRoutingRuleTimeRestrictions(timeRestriction)
	}

	log.Printf("[INFO] Creating OpsGenie team routing rule '%s'", name)

	result, err := client.CreateRoutingRule(context.Background(), createRequest)
	if err != nil {
		return err
	}
	d.SetId(result.Id)

	return resourceOpsGenieTeamRoutingRuleRead(d, meta)
}

func resourceOpsGenieTeamRoutingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	getRequest := &team.GetRoutingRuleRequest{
		TeamIdentifierType:  team.Id,
		TeamIdentifierValue: d.Get("team_id").(string),
		RoutingRuleId:       d.Id(),
	}

	result, err := client.GetRoutingRule(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.Set("name", result.Name)
	d.Set("time_restriction", flattenOpsgenieTimeRestriction(result.TimeRestriction))
	d.Set("notify", flattenOpsgenieNotify(result.Notify))
	d.Set("criteria", flattenOpsgenieCriteria(result.Criteria))
	d.Set("timezone", result.Timezone)

	return nil
}

func resourceOpsGenieTeamRoutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	teamId := d.Get("team_id").(string)
	timezone := d.Get("timezone").(string)
	timeRestriction := d.Get("time_restriction").([]interface{})
	criteria := d.Get("criteria").([]interface{})
	notify := d.Get("notify").([]interface{})

	expandedCriteria := expandOpsgenieCriteria(criteria)
	if err := validateOpsgenieCriteria(expandedCriteria); err != nil {
		return err
	}

	updateRequest := &team.UpdateRoutingRuleRequest{
		TeamIdentifierType:  team.Id,
		TeamIdentifierValue: teamId,
		RoutingRuleId:       d.Id(),
		Name:                name,
		Timezone:            timezone,
		Criteria:            expandedCriteria,
		Notify:              expandOpsgenieNotify(notify),
	}
	if len(timeRestriction) > 0 {
		updateRequest.TimeRestriction = expandRoutingRuleTimeRestrictions(timeRestriction)
	}

	log.Printf("[INFO] Updating OpsGenie team routing rule '%s'", name)

	_, err = client.UpdateRoutingRule(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieTeamRoutingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie team routing rule'%s'", d.Get("name").(string))
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &team.DeleteRoutingRuleRequest{
		TeamIdentifierType:  team.Id,
		TeamIdentifierValue: d.Get("team_id").(string),
		RoutingRuleId:       d.Id(),
	}

	_, err = client.DeleteRoutingRule(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func flattenOpsgenieNotify(input team.Notify) []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, 1)
	out := make(map[string]interface{})
	out["type"] = input.Type

	out["id"] = input.Id
	out["name"] = input.Name
	rules = append(rules, out)
	return rules
}

func flattenOpsgenieCriteria(input og.Criteria) []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, 1)
	out := make(map[string]interface{})
	out["type"] = input.CriteriaType
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
func expandRoutingRuleTimeRestrictions(d []interface{}) *og.TimeRestriction {
	timeRestriction := og.TimeRestriction{}
	for _, v := range d {
		config := v.(map[string]interface{})

		timeRestrictionType := config["type"].(string)
		timeRestriction.Type = og.RestrictionType(timeRestrictionType)

		if len(config["restrictions"].([]interface{})) > 0 {
			timeRestriction.RestrictionList = expandOpsgenieRestrictions(config["restrictions"].([]interface{}))
		} else {
			timeRestriction.Restriction = expandOpsgenieRestriction(config["restriction"].([]interface{}))
		}
	}

	return &timeRestriction
}

func flattenOpsgenieTimeRestriction(input og.TimeRestriction) []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, 1)
	out := make(map[string]interface{})
	out["type"] = input.Type

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
		out["restrictions"] = restrictions
		rules = append(rules, out)
		return rules
	} else {
		restriction := make(map[string]interface{})
		//IF RESTRICTION
		restriction["end_hour"] = input.Restriction.EndHour
		restriction["end_min"] = input.Restriction.EndMin
		restriction["start_hour"] = input.Restriction.StartHour
		restriction["start_min"] = input.Restriction.StartMin

		//IF restrictions
		out["restriction"] = restriction
		rules = append(rules, out)
		return rules
	}
}

func expandOpsgenieNotify(input []interface{}) *team.Notify {

	notify := team.Notify{}
	for _, r := range input {
		inputMap := r.(map[string]interface{})
		notifyType := inputMap["type"].(string)
		var notifyName string
		if inputMap["name"] != nil {
			notifyName = inputMap["name"].(string)
		}
		if inputMap["id"] != nil {
			notify.Id = inputMap["id"].(string)

		}
		notify.Name = notifyName
		notify.Type = team.NotifyType(notifyType)
	}
	return &notify

}

func validateOpsgenieCriteria(criteria *og.Criteria) error {
	if criteria.CriteriaType == "match-all" && len(criteria.Conditions) > 0 {
		return fmt.Errorf("criteria cannot have conditions set when type is match-all: %v", criteria)
	}
	return nil
}

func expandOpsgenieCriteria(input []interface{}) *og.Criteria {

	criteria := og.Criteria{}
	for _, r := range input {
		inputMap := r.(map[string]interface{})
		criteriaType := inputMap["type"].(string)
		conditions := expandOpsgenieConditions(inputMap["conditions"].([]interface{}))
		criteria.Conditions = conditions
		criteria.CriteriaType = og.ConditionMatchType(criteriaType)
	}
	return &criteria
}

func expandOpsgenieConditions(input []interface{}) []og.Condition {

	conditions := make([]og.Condition, 0, len(input))

	if input == nil {
		return conditions
	}
	for _, v := range input {
		inputMap := v.(map[string]interface{})
		condition := og.Condition{}
		condition.Field = og.ConditionFieldType(inputMap["field"].(string))
		isNot := inputMap["not"].(bool)
		condition.IsNot = &isNot
		condition.Operation = og.ConditionOperation(inputMap["operation"].(string))
		condition.ExpectedValue = inputMap["expected_value"].(string)
		key := inputMap["key"].(string)
		if key != "" {
			condition.Key = inputMap["key"].(string)
		}
		order := inputMap["order"].(int)
		condition.Order = &order
		conditions = append(conditions, condition)
	}

	return conditions
}
