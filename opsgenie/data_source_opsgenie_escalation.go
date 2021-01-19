package opsgenie

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/escalation"
)

func dataSourceOpsgenieEscalation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsgenieEscalationRead,
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
				Optional: true,
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
							Default:  0,
						},
						"count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  20,
						},
						"reset_recipient_states": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"close_alert_after_all": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func dataSourceOpsgenieEscalationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := escalation.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	escalationName := d.Get("name").(string)

	getRequest := &escalation.GetRequest{
		IdentifierType: escalation.Name,
		Identifier:     escalationName,
	}

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.SetId(getResponse.Id)
	d.Set("name", getResponse.Name)
	d.Set("description", getResponse.Description)
	d.Set("rules", flattenOpsgenieEscalationRules(getResponse.Rules))
	d.Set("repeat", flattenOpsgenieEscalationRepeat(getResponse.Repeat))
	if getResponse.OwnerTeam != nil {
		d.Set("owner_team_id", getResponse.OwnerTeam.Id)
	}

	return nil
}
