package opsgenie

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

func resourceOpsgenieHeartbeat() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieHeartbeatCreate,
		Read:   handleNonExistentResource(resourceOpsgenieHeartbeatRead),
		Update: resourceOpsgenieHeartbeatUpdate,
		Delete: resourceOpsgenieHeartbeatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOpsgenieHeartbeat,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"interval_unit": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
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

func resourceOpsgenieHeartbeatCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := heartbeat.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	interval := d.Get("interval").(int)
	intervalUnit := d.Get("interval_unit").(string)
	enabled := d.Get("enabled").(bool)
	ownerTeamId := d.Get("owner_team_id").(string)
	alertMessage := d.Get("alert_message").(string)
	alertPriority := d.Get("alert_priority").(string)

	addRequest := &heartbeat.AddRequest{
		Name:          name,
		Description:   description,
		Interval:      interval,
		IntervalUnit:  heartbeat.Unit(intervalUnit),
		Enabled:       &enabled,
		AlertMessage:  alertMessage,
		AlertTag:      flattenTags(d),
		AlertPriority: alertPriority,
	}
	if ownerTeamId != "" {
		addRequest.OwnerTeam = og.OwnerTeam{
			Id: ownerTeamId,
		}
	}

	result, err := client.Add(context.Background(), addRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Heartbeat.Name)

	return resourceOpsgenieHeartbeatRead(d, meta)
}

func resourceOpsgenieHeartbeatRead(d *schema.ResourceData, meta interface{}) error {
	client, err := heartbeat.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	result, err := client.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}

	d.Set("name", result.Name)
	d.Set("enabled", result.Enabled)
	d.Set("description", result.Description)
	d.Set("interval", result.Interval)
	d.Set("interval_unit", result.IntervalUnit)
	d.Set("owner_team_id", result.OwnerTeam.Id)
	d.Set("alert_priority", result.AlertPriority)
	d.Set("alert_tags", result.AlertTags)
	d.Set("alert_message", result.AlertMessage)

	return nil
}

func resourceOpsgenieHeartbeatUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := heartbeat.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	interval := d.Get("interval").(int)
	intervalUnit := d.Get("interval_unit").(string)
	enabled := d.Get("enabled").(bool)
	ownerTeamId := d.Get("owner_team_id").(string)
	alertMessage := d.Get("alert_message").(string)
	alertPriority := d.Get("alert_priority").(string)

	updateRequest := &heartbeat.UpdateRequest{
		Name:          name,
		Description:   description,
		Interval:      interval,
		IntervalUnit:  heartbeat.Unit(intervalUnit),
		Enabled:       &enabled,
		AlertMessage:  alertMessage,
		AlertTag:      flattenTags(d),
		AlertPriority: alertPriority,
	}
	if ownerTeamId != "" {
		updateRequest.OwnerTeam = og.OwnerTeam{
			Id: ownerTeamId,
		}
	}

	_, err = client.Update(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieHeartbeatDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := heartbeat.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	_, err = client.Delete(context.Background(), d.Id())
	if err != nil {
		return err
	}

	return nil
}

func flattenTags(d *schema.ResourceData) []string {
	input := d.Get("alert_tags").(*schema.Set)
	tags := make([]string, len(input.List()))

	if input == nil {
		return tags
	}

	for k, v := range input.List() {
		tags[k] = v.(string)
	}

	return tags
}

func validateOpsgenieHeartbeat(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alpha numeric characters and underscores are allowed in %q: %q", k, value))
	}

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}
