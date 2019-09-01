package opsgenie

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/integration"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

func resourceOpsgenieApiIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieApiIntegrationCreate,
		Read:   resourceOpsgenieApiIntegrationRead,
		Update: resourceOpsgenieApiIntegrationUpdate,
		Delete: resourceOpsgenieApiIntegrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsgenieIntegrationName,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allow_write_access": {
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceOpsgenieApiIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	allowWriteAccess := d.Get("allow_write_access").(bool)
	ignoreRespondersFromPayload := d.Get("ignore_responders_from_payload").(bool)
	suppressNotifications := d.Get("suppress_notifications").(bool)
	ownerTeam := d.Get("owner_team_id").(string)
	enabled := d.Get("enabled").(bool)

	createRequest := &integration.APIBasedIntegrationRequest{
		Name:                        name,
		Type:                        ApiIntegrationType,
		AllowWriteAccess:            allowWriteAccess,
		IgnoreRespondersFromPayload: ignoreRespondersFromPayload,
		SuppressNotifications:       suppressNotifications,
		Responders:                  expandOpsgenieIntegrationResponders(d),
	}
	if ownerTeam != "" {
		createRequest.OwnerTeam = &og.OwnerTeam{
			Id: ownerTeam,
		}
	}

	log.Printf("[INFO] Creating OpsGenie api integration '%s'", name)

	result, err := client.CreateApiBased(context.Background(), createRequest)
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
		log.Printf("[INFO] Enabled OpsGenie api integration '%s'", name)

	}

	return resourceOpsgenieApiIntegrationRead(d, meta)
}

func resourceOpsgenieApiIntegrationRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("name", result.Data["name"])
	d.Set("id", result.Data["id"])
	d.Set("responders", result.Data["responders"])
	d.Set("ignore_responders_from_payload", result.Data["ignoreRespondersFromPayload"])
	d.Set("suppress_notifications", result.Data["suppressNotifications"])

	return nil
}

func resourceOpsgenieApiIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	ignoreRespondersFromPayload := d.Get("ignore_responders_from_payload").(bool)
	suppressNotifications := d.Get("suppress_notifications").(bool)
	enabled := d.Get("enabled").(bool)

	updateRequest := &integration.UpdateIntegrationRequest{
		Id:                          d.Id(),
		Name:                        name,
		Type:                        ApiIntegrationType,
		IgnoreRespondersFromPayload: ignoreRespondersFromPayload,
		SuppressNotifications:       suppressNotifications,
		Responders:                  expandOpsgenieIntegrationResponders(d),
		Enabled:                     enabled,
	}

	log.Printf("[INFO] Updating OpsGenie api based integration '%s'", name)

	_, err = client.ForceUpdateAllFields(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieApiIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie api integration '%s'", d.Get("name").(string))
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
