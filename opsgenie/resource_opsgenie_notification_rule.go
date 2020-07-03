package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/notification"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

func resourceOpsGenieNotificationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieNotificationRuleCreate,
		Read:   resourceOpsGenieNotificationRuleRead,
		Update: resourceOpsGenieNotificationRuleUpdate,
		Delete: resourceOpsGenieNotificationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected username/notification_rule_id", d.Id())
				}
				d.Set("username", idParts[0])
				d.SetId(idParts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
				ForceNew:     true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"create-alert", "acknowledged-alert", "closed-alert", "assigned-alert", "add-note",
					"schedule-start", "schedule-end", "incoming-call-routing",
				}, false),
				ForceNew: true,
			},
			"notification_time": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"just-before", "15-minutes-ago", "1-hour-ago", "1-day-ago",
					}, false),
				},
			},
			"steps": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"send_after": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"contact": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"email", "sms", "voice", "mobile"}, false),
									},
									"to": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"repeat": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"loop_after": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
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
							Type:     schema.TypeSet,
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

func resourceOpsGenieNotificationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := notification.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	enabled := d.Get("enabled").(bool)
	timeRestriction := d.Get("time_restriction").([]interface{})

	createRequest := &notification.CreateRuleRequest{
		UserIdentifier:   d.Get("username").(string),
		Name:             d.Get("name").(string),
		ActionType:       notification.ActionType(d.Get("action_type").(string)),
		NotificationTime: expandOpsGenieNotificationRuleNotificationTime(d.Get("notification_time").(*schema.Set)),
		Enabled:          &enabled,
		Order:            uint32(d.Get("order").(int)),
		Repeat:           expandOpsGenieNotificationRuleRepeat(d),
	}

	if len(d.Get("steps").([]interface{})) > 0 {
		createRequest.Steps = expandOpsGenieNotificationRuleSteps(d.Get("steps").([]interface{}))
	}

	if len(timeRestriction) > 0 {
		createRequest.TimeRestriction = expandTimeRestrictions(timeRestriction)
	}

	log.Printf("[INFO] Creating Notification Rule '%s' for User: '%s'", d.Get("name").(string), d.Get("username").(string))
	result, err := client.CreateRule(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.SimpleNotificationRule.Id)

	return resourceOpsGenieNotificationRuleRead(d, meta)
}

func resourceOpsGenieNotificationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := notification.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	username := d.Get("username").(string)

	log.Printf("[INFO] Reading OpsGenie Notification Rule '%s' for user '%s'", name, username)

	rule, err := client.GetRule(context.Background(), &notification.GetRuleRequest{
		UserIdentifier: username,
		RuleId:         d.Id(),
	})
	if err != nil {
		x := err.(*ogClient.ApiError)
		if x.StatusCode == 404 {
			log.Printf("[WARN] Removing Notification Rule because it's gone %s", name)
			d.SetId("")
			return nil
		}
	}
	d.Set("name", rule.Name)
	d.Set("action_type", rule.ActionType)
	d.Set("notification_time", rule.NotificationTime)
	d.Set("enabled", rule.Enabled)
	d.Set("order", rule.Order)
	d.Set("time_restriction", flattenOpsgenieNotificationRuleTimeRestriction(rule.TimeRestriction))

	if rule.Steps != nil {
		d.Set("steps", flattenOpsGenieNotificationRuleSteps(rule.Steps))
	} else {
		d.Set("steps", nil)
	}

	return nil
}

func resourceOpsGenieNotificationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := notification.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	enabled := d.Get("enabled").(bool)
	timeRestriction := d.Get("time_restriction").([]interface{})

	updateRequest := &notification.UpdateRuleRequest{
		UserIdentifier:   d.Get("username").(string),
		RuleId:           d.Id(),
		NotificationTime: expandOpsGenieNotificationRuleNotificationTime(d.Get("notification_time").(*schema.Set)),
		Enabled:          &enabled,
		Order:            uint32(d.Get("order").(int)),
		Repeat:           expandOpsGenieNotificationRuleRepeat(d),
	}

	if len(d.Get("steps").([]interface{})) > 0 {
		updateRequest.Steps = expandOpsGenieNotificationRuleSteps(d.Get("steps").([]interface{}))
	}

	if len(timeRestriction) > 0 {
		updateRequest.TimeRestriction = expandTimeRestrictions(timeRestriction)
	}

	log.Printf("[INFO] Updating Notification Rule '%s' for User: '%s'", d.Get("name").(string), d.Get("username").(string))
	result, err := client.UpdateRule(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	d.SetId(result.SimpleNotificationRule.Id)

	return resourceOpsGenieNotificationRuleRead(d, meta)
}

func resourceOpsGenieNotificationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie Notification Rule '%s' for user '%s'", d.Get("name").(string), d.Get("username").(string))
	client, err := notification.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &notification.DeleteRuleRequest{
		UserIdentifier: d.Get("username").(string),
		RuleId:         d.Id(),
	}

	_, err = client.DeleteRule(context.Background(), deleteRequest)
	if err != nil {
		return err
	}
	return nil
}

func expandOpsGenieNotificationRuleNotificationTime(input *schema.Set) []notification.NotificationTimeType {
	output := make([]notification.NotificationTimeType, 0)

	if input == nil {
		return output
	}

	for _, v := range input.List() {
		output = append(output, notification.NotificationTimeType(v.(string)))
	}
	return output
}

func expandOpsGenieNotificationRuleSteps(input []interface{}) []*og.Step {
	output := make([]*og.Step, 0)
	if input == nil {
		return output
	}
	for _, v := range input {
		config := v.(map[string]interface{})
		enabled := config["enabled"].(bool)
		element := og.Step{}
		element.Enabled = &enabled
		element.Contact = expandOpsGenieNotificationRuleStepsContact(config["contact"].([]interface{}))
		if config["send_after"].(int) > 0 {
			element.SendAfter = &og.SendAfter{
				TimeUnit:   "minute",
				TimeAmount: uint32(config["send_after"].(int)),
			}
		}
		output = append(output, &element)
	}
	return output
}

func expandOpsGenieNotificationRuleStepsContact(input []interface{}) og.Contact {
	output := og.Contact{}
	if input == nil {
		return output
	}
	for _, v := range input {
		config := v.(map[string]interface{})
		output.To = config["to"].(string)
		output.MethodOfContact = og.MethodType(config["method"].(string))
	}
	return output
}

func flattenOpsGenieNotificationRuleSteps(input []*notification.StepResult) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, len(input))
	for _, v := range input {
		element := make(map[string]interface{})
		element["enabled"] = v.Enabled
		element["contact"] = flattenOpsGenieNotificationRuleStepsContact(v.Contact)
		if v.SendAfter != nil {
			element["send_after"] = v.SendAfter.TimeAmount
		}
		output = append(output, element)
	}
	return output
}

func flattenOpsGenieNotificationRuleStepsContact(input og.Contact) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, 1)
	element := make(map[string]interface{})
	element["to"] = input.To
	element["method"] = input.MethodOfContact
	output = append(output, element)
	return output
}

func expandOpsGenieNotificationRuleRepeat(d *schema.ResourceData) *notification.Repeat {
	input := d.Get("repeat").([]interface{})
	repeat := notification.Repeat{}
	for _, r := range input {
		repeatMap := r.(map[string]interface{})
		repeatEnabled := repeatMap["enabled"].(bool)
		repeat.LoopAfter = uint32(repeatMap["loop_after"].(int))
		repeat.Enabled = &repeatEnabled
	}
	return &repeat
}

func flattenOpsgenieNotificationRuleTimeRestriction(input *og.TimeRestriction) []map[string]interface{} {
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
		return restrictions
	} else {
		restriction := make(map[string]interface{})
		//IF RESTRICTION
		restriction["end_day"] = input.Restriction.EndDay
		restriction["end_hour"] = input.Restriction.EndHour
		restriction["end_min"] = input.Restriction.EndMin
		restriction["start_day"] = input.Restriction.StartDay
		restriction["start_hour"] = input.Restriction.StartHour
		restriction["start_min"] = input.Restriction.StartMin

		//IF restrictions
		rules = append(rules, restriction)
		return rules
	}
}
