package opsgenie

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"

	"fmt"

	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpsgenieSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieScheduleCreate,
		Read:   handleNonExistentResource(resourceOpsgenieScheduleRead),
		Update: resourceOpsgenieScheduleUpdate,
		Delete: resourceOpsgenieScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsgenieScheduleName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOpsgenieScheduleDescription,
			},
			"timezone": {
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
		},
	}
}

func resourceOpsgenieScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	timezone := d.Get("timezone").(string)
	enabled := d.Get("enabled").(bool)
	ownerTeam := d.Get("owner_team_id").(string)

	createRequest := &schedule.CreateRequest{
		Name:        name,
		Enabled:     &enabled,
		Description: description,
		Timezone:    timezone,
	}
	if ownerTeam != "" {
		createRequest.OwnerTeam = &og.OwnerTeam{
			Id: ownerTeam,
		}
	}

	log.Printf("[INFO] Creating OpsGenie schedule '%s'", name)

	result, err := client.Create(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsgenieScheduleRead(d, meta)
}

func resourceOpsgenieScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	getRequest := &schedule.GetRequest{
		IdentifierType:  schedule.Id,
		IdentifierValue: d.Id(),
	}

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.Set("name", getResponse.Schedule.Name)
	if getResponse.Schedule.OwnerTeam != nil {
		d.Set("owner_team_id", getResponse.Schedule.OwnerTeam.Id)
	}
	d.Set("timezone", getResponse.Schedule.Timezone)
	d.Set("description", getResponse.Schedule.Description)
	d.Set("enabled", getResponse.Schedule.Enabled)

	return nil
}

func resourceOpsgenieScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	timezone := d.Get("timezone").(string)
	enabled := d.Get("enabled").(bool)
	ownerTeam := d.Get("owner_team_id").(string)

	updateRequest := &schedule.UpdateRequest{
		IdentifierType:  schedule.Id,
		IdentifierValue: d.Id(),
		Name:            name,
		Enabled:         &enabled,
		Description:     description,
		Timezone:        timezone,
	}

	if ownerTeam != "" {
		updateRequest.OwnerTeam = &og.OwnerTeam{
			Id: ownerTeam,
		}
	} else {
		updateRequest.OwnerTeam = nil
	}
	log.Printf("[INFO] Updating OpsGenie schedule '%s'", name)

	_, err = client.Update(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie schedule '%s'", d.Get("name").(string))
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &schedule.DeleteRequest{
		IdentifierType:  schedule.Id,
		IdentifierValue: d.Id(),
	}

	_, err = client.Delete(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func validateOpsgenieScheduleName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z 0-9_-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alpha numeric characters and underscores are allowed in %q: %q", k, value))
	}

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}
func validateOpsgenieScheduleDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 10000 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}
