package opsgenie

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opsgenie/opsgenie-go-sdk-v2/integration"
)

func resourceOpsgenieEmailIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieEmailIntegrationCreate,
		Read:   handleNonExistentResource(resourceOpsgenieEmailIntegrationRead),
		Update: resourceOpsgenieEmailIntegrationUpdate,
		Delete: resourceOpsgenieEmailIntegrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 250),
			},
			"email_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ignore_responders_from_payload": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"suppress_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"owner_team_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"responders": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateResponderType,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpsgenieEmailIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	emailUsername := d.Get("email_username").(string)
	ignoreRespondersFromPayload := d.Get("ignore_responders_from_payload").(bool)
	suppressNotifications := d.Get("suppress_notifications").(bool)
	enabled := d.Get("enabled").(bool)
	ownerTeam := d.Get("owner_team_id").(string)

	createRequest := &integration.EmailBasedIntegrationRequest{
		Name:                        name,
		Type:                        EmailIntegrationType,
		EmailUsername:               emailUsername,
		IgnoreRespondersFromPayload: &ignoreRespondersFromPayload,
		SuppressNotifications:       &suppressNotifications,
		Responders:                  expandOpsgenieIntegrationResponders(d),
	}

	if ownerTeam != "" {
		createRequest.OwnerTeam = &og.OwnerTeam{
			Id: ownerTeam,
		}
	}

	log.Printf("[INFO] Creating OpsGenie email integration '%s'", name)

	result, err := client.CreateEmailBased(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	if enabled {
		_, err = client.Enable(context.Background(), &integration.EnableIntegrationRequest{
			Id: result.Id,
		})
		if err != nil {
			return err
		}
		log.Printf("[INFO] Enabled OpsGenie email integration '%s'", name)

	}

	return resourceOpsgenieEmailIntegrationRead(d, meta)
}

func resourceOpsgenieEmailIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	result, err := client.Get(context.Background(), &integration.GetRequest{
		Id: d.Id(),
	})
	if err != nil {
		return err
	}

	if result.Data["ownerTeam"] != nil {
		ownerTeam := result.Data["ownerTeam"].(map[string]interface{})
		d.Set("owner_team_id", ownerTeam["id"])
	} else if result.Data["responders"] != nil {
		d.Set("responders", flattenIntegrationResponders(result.Data["responders"].([]interface{})))
	}
	d.Set("name", result.Data["name"])
	d.Set("suppress_notifications", result.Data["suppressNotifications"])
	d.Set("email_username", result.Data["emailUsername"])
	d.Set("enabled", result.Data["enabled"])
	return nil
}

func resourceOpsgenieEmailIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	emailUsername := d.Get("email_username").(string)
	ignoreRespondersFromPayload := d.Get("ignore_responders_from_payload").(bool)
	suppressNotifications := d.Get("suppress_notifications").(bool)
	enabled := d.Get("enabled").(bool)

	updateRequest := &integration.UpdateIntegrationRequest{
		Id:                          d.Id(),
		Name:                        name,
		Type:                        EmailIntegrationType,
		EmailUsername:               emailUsername,
		IgnoreRespondersFromPayload: &ignoreRespondersFromPayload,
		SuppressNotifications:       &suppressNotifications,
		Responders:                  expandOpsgenieIntegrationResponders(d),
		Enabled:                     &enabled,
	}

	log.Printf("[INFO] Updating OpsGenie email based integration '%s'", name)

	_, err = client.ForceUpdateAllFields(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieEmailIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie email integration '%s'", d.Get("name").(string))
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &integration.DeleteIntegrationRequest{
		Id: d.Id(),
	}

	_, err = client.Delete(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}
