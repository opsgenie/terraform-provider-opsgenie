package opsgenie

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/service"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpsGenieService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieServiceCreate,
		Read:   handleNonExistentResource(resourceOpsGenieServiceRead),
		Update: resourceOpsGenieServiceUpdate,
		Delete: resourceOpsGenieServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieServiceName,
			},
			"team_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieServiceTeamId,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOpsGenieServiceDescription,
			},
		},
	}
}

func resourceOpsGenieServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	teamId := d.Get("team_id").(string)
	description := d.Get("description").(string)

	createRequest := &service.CreateRequest{
		Name:        name,
		TeamId:      teamId,
		Description: description,
	}

	log.Printf("[INFO] Creating OpsGenie service '%s'", name)
	result, err := client.Create(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsGenieServiceRead(d, meta)
}

func resourceOpsGenieServiceRead(d *schema.ResourceData, meta interface{}) error {
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	log.Printf("[INFO] Reading OpsGenie service '%s'", name)

	res, err := client.Get(context.Background(), &service.GetRequest{
		Id: d.Id(),
	})
	if err != nil {
		return err
	}

	d.Set("name", res.Service.Name)
	d.Set("team_id", res.Service.TeamId)
	d.Set("description", res.Service.Description)

	return nil
}

func resourceOpsGenieServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	log.Printf("[INFO] Updating OpsGenie service '%s'", name)

	updateRequest := &service.UpdateRequest{
		Id:          d.Id(),
		Name:        name,
		Description: description,
	}

	_, err = client.Update(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieServiceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie service '%s'", d.Get("name").(string))
	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &service.DeleteRequest{
		Id: d.Id(),
	}

	_, err = client.Delete(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func validateOpsGenieServiceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}

func validateOpsGenieServiceTeamId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 512 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 512 characters: %q %d", k, value, len(value)))
	}

	return
}

func validateOpsGenieServiceDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 10000 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 10000 characters: %q %d", k, value, len(value)))
	}

	return
}
