package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/service"
	"log"
	"strconv"
	"strings"
	"time"
)

func dataSourceOpsGenieService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieServiceRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieServiceName,
			},
			"team_id": {
				Type:         schema.TypeString,
				Optional:     true,
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
	// OpsGenie async call to create service might take a bit of time to take affect.
	// This sleep will make sure we are not hitting 404 error if hit get/list service API before creation could happen.
	time.Sleep(5 * time.Second)

	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	log.Printf("[INFO] Reading OpsGenie service '%s'", name)

	breakFlag := false
	offset := 0

	for {
		res, err := client.List(context.Background(), &service.ListRequest{
			Limit:  100,
			Offset: offset,
		})
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Searching for service name: '%s' in your account", name)
		for _, srvObj := range res.Services {
			if name == srvObj.Name {
				log.Printf("[DEBUG] Found the service")
				log.Printf("[DEBUG] Service ID: '%s'", srvObj.Id)
				d.Set("name", srvObj.Name)
				d.Set("team_id", srvObj.TeamId)
				d.Set("description", srvObj.Description)
				d.SetId(srvObj.Id)
				breakFlag = true
				break
			}
		}

		if res.Paging.Next == "" || breakFlag {
			break
		}

		offsetString := strings.Split(res.Paging.Next, string('&'))[2]
		offsetString = strings.Split(offsetString, string('='))[1]
		offset, err = strconv.Atoi(offsetString)

		if err != nil {
			return err
		}
	}
	return nil
}
