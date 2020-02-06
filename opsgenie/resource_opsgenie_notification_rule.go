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
		},
	}
}

func resourceOpsGenieNotificationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := notification.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	enabled := d.Get("enabled").(bool)
	createRequest := &notification.CreateRuleRequest{
		UserIdentifier:   d.Get("username").(string),
		Name:             d.Get("name").(string),
		ActionType:       notification.ActionType(d.Get("action_type").(string)),
		NotificationTime: expandOpsGenieNotificationRuleNotificationTime(d.Get("notification_time").(*schema.Set)),
		Enabled:          &enabled,
	}

	if len(d.Get("steps").([]interface{})) > 0 {
		createRequest.Steps = expandOpsGenieNotificationRuleSteps(d.Get("steps").([]interface{}))
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
	updateRequest := &notification.UpdateRuleRequest{
		UserIdentifier:   d.Get("username").(string),
		RuleId:           d.Id(),
		NotificationTime: expandOpsGenieNotificationRuleNotificationTime(d.Get("notification_time").(*schema.Set)),
		Enabled:          &enabled,
	}

	if len(d.Get("steps").([]interface{})) > 0 {
		updateRequest.Steps = expandOpsGenieNotificationRuleSteps(d.Get("steps").([]interface{}))
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
