package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/service"
	"log"
)

func dataSourceOpsGenieService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieServiceRead,
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

func dataSourceOpsGenieServiceRead(d *schema.ResourceData, meta interface{}) error {
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
