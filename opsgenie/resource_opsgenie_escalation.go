package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/escalation"
)

func resourceOpsgenieEscalation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieEscalationCreate,
		Read:   handleNonExistentResource(resourceOpsgenieEscalationRead),
		Update: resourceOpsgenieEscalationUpdate,
		Delete: resourceOpsgenieEscalationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateOpsgenieEscalationRulesCondition,
						},
						"notify_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateOpsgenieEscalationRulesNotifyType,
						},
						"recipient": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateEscalationParticipantType,
									},
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"delay": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"owner_team_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"wait_interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"reset_recipient_states": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"close_alert_after_all": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpsgenieEscalationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := escalation.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	ownerTeam := d.Get("owner_team_id").(string)

	createRequest := &escalation.CreateRequest{
		Name:        name,
		Description: description,
		Rules:       expandOpsgenieEscalationRules(d),
		Repeat:      expandOpsgenieEscalationRepeat(d),
	}

	if ownerTeam != "" {
		createRequest.OwnerTeam = &og.OwnerTeam{
			Id: ownerTeam,
		}
	}

	log.Printf("[INFO] Creating OpsGenie escalation '%s'", name)

	result, err := client.Create(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsgenieEscalationRead(d, meta)
}

func resourceOpsgenieEscalationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := escalation.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	getRequest := &escalation.GetRequest{
		IdentifierType: escalation.Id,
		Identifier:     d.Id(),
	}

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.Set("name", getResponse.Name)
	d.Set("description", getResponse.Description)
	d.Set("rules", flattenOpsgenieEscalationRules(getResponse.Rules))
	repeat := d.Get("repeat").([]interface{})
	if len(repeat) > 0 {
		d.Set("repeat", flattenOpsgenieEscalationRepeat(getResponse.Repeat))
	}
	if getResponse.OwnerTeam != nil {
		d.Set("owner_team_id", getResponse.OwnerTeam.Id)
	}

	return nil
}

func resourceOpsgenieEscalationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := escalation.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	ownerTeam := d.Get("owner_team_id").(string)

	updateRequest := &escalation.UpdateRequest{
		IdentifierType: escalation.Id,
		Identifier:     d.Id(),
		Name:           name,
		Description:    description,
		Rules:          expandOpsgenieEscalationRules(d),
		Repeat:         expandOpsgenieEscalationRepeat(d),
	}
	if ownerTeam != "" {
		updateRequest.OwnerTeam = &og.OwnerTeam{
			Id: ownerTeam,
		}
	} else {
		updateRequest.OwnerTeam = nil
	}
	log.Printf("[INFO] Updating OpsGenie escalation '%s'", name)

	_, err = client.Update(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieEscalationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie escalation '%s'", d.Get("name").(string))
	client, err := escalation.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &escalation.DeleteRequest{
		IdentifierType: escalation.Id,
		Identifier:     d.Id(),
	}

	_, err = client.Delete(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func flattenOpsgenieEscalationRules(input []escalation.Rule) []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, len(input))
	for _, rule := range input {
		out := make(map[string]interface{})
		out["notify_type"] = rule.NotifyType
		out["condition"] = rule.Condition
		out["delay"] = rule.Delay.TimeAmount
		recipientArr := make([]map[string]interface{}, 0, 1)
		recipient := make(map[string]interface{})
		recipient["id"] = rule.Recipient.Id
		recipient["type"] = rule.Recipient.Type
		recipientArr = append(recipientArr, recipient)
		out["recipient"] = recipientArr
		rules = append(rules, out)
	}

	return rules
}

func flattenOpsgenieEscalationRepeat(input escalation.Repeat) []map[string]interface{} {
	repeats := make([]map[string]interface{}, 0, 1)
	out := make(map[string]interface{})
	out["count"] = input.Count
	out["wait_interval"] = input.WaitInterval
	out["close_alert_after_all"] = input.CloseAlertAfterAll
	out["reset_recipient_states"] = input.ResetRecipientStates
	repeats = append(repeats, out)
	return repeats
}

func expandOpsgenieEscalationRules(d *schema.ResourceData) []escalation.RuleRequest {
	input := d.Get("rules").([]interface{})

	rules := make([]escalation.RuleRequest, 0, len(input))
	if input == nil {
		return rules
	}
	for _, v := range input {
		config := v.(map[string]interface{})

		condition := config["condition"].(string)
		notifyType := config["notify_type"].(string)
		recipient := config["recipient"].([]interface{})
		delay := config["delay"].(int)

		rule := escalation.RuleRequest{
			Condition:  og.EscalationCondition(condition),
			NotifyType: og.NotifyType(notifyType),
			Recipient:  expandOpsgenieEscalationRuleParticipant(recipient),
			Delay: escalation.EscalationDelayRequest{
				TimeAmount: uint32(delay),
			},
		}

		rules = append(rules, rule)
	}

	return rules
}

func expandOpsgenieEscalationRuleParticipant(pp []interface{}) og.Participant {
	participant := og.Participant{}
	for _, b := range pp {
		p := b.(map[string]interface{})
		participant.Type = og.ParticipantType(p["type"].(string))

		id := p["id"]
		participant.Id = id.(string)
	}
	return participant
}

func expandOpsgenieEscalationRepeat(d *schema.ResourceData) *escalation.RepeatRequest {
	input := d.Get("repeat").([]interface{})
	repeat := escalation.RepeatRequest{}
	for _, r := range input {
		repeatMap := r.(map[string]interface{})
		resetRecipientState := repeatMap["reset_recipient_states"].(bool)
		closeAlertAfterAll := repeatMap["close_alert_after_all"].(bool)
		repeat.WaitInterval = uint32(repeatMap["wait_interval"].(int))
		repeat.Count = uint32(repeatMap["count"].(int))
		repeat.ResetRecipientStates = &resetRecipientState
		repeat.CloseAlertAfterAll = &closeAlertAfterAll
	}

	return &repeat
}

func validateOpsgenieEscalationRulesNotifyType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"default":  true,
		"next":     true,
		"previous": true,
		"users":    true,
		"admins":   true,
		"random":   true,
		"all":      true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("It can only be one of these 'user', 'schedule', 'team'"))
	}
	return
}

func validateOpsgenieEscalationRulesCondition(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"if-not-acked":  true,
		"if-not-closed": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("It can only be one of these 'user', 'schedule', 'team'"))
	}
	return
}

func validateEscalationParticipantType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"user":     true,
		"team":     true,
		"schedule": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("It can only be one of these 'user', 'schedule', 'team'"))
	}
	return
}
