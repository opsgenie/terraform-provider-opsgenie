package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
)

func dataSourceOpsgenieHeartbeat() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsgenieHeartbeatRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsgenieHeartbeat,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"interval_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"owner_team_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"alert_priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceOpsgenieHeartbeatRead(d *schema.ResourceData, meta interface{}) error {
	client, err := heartbeat.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	heartbeatName := d.Get("name").(string)

	result, err := client.Get(context.Background(), heartbeatName)
	if err != nil {
		return err
	}

	d.SetId(result.Name)
	err = d.Set("name", result.Name)
	err = d.Set("enabled", result.Enabled)
	err = d.Set("description", result.Description)
	err = d.Set("interval", result.Interval)
	err = d.Set("interval_unit", result.IntervalUnit)
	err = d.Set("owner_team_id", result.OwnerTeam.Id)
	err = d.Set("alert_priority", result.AlertPriority)
	err = d.Set("alert_tags", result.AlertTags)
	err = d.Set("alert_message", result.AlertMessage)

	return nil
}
