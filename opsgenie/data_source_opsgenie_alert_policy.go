package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
	"log"
	"time"
)

func dataSourceOpsGenieAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieAlertPolicyRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"team_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceOpsGenieAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {

	// OpsGenie async call to create service might take a bit of time to take affect.
	// This sleep will make sure we are not hitting 404 error if hit get/list service API before creation could happen.
	time.Sleep(5 * time.Second)

	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	team_id := d.Get("team_id").(string)

	log.Printf("[INFO] Reading OpsGenie Alert Policy '%s'", name)

	listAlertPoliciesRequest := policy.ListAlertPoliciesRequest{}

	if team_id != "" {
		listAlertPoliciesRequest = policy.ListAlertPoliciesRequest{TeamId: team_id}
	}
	res, err := client.ListAlertPolicies(context.Background(), &listAlertPoliciesRequest)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Searching for Alert Policy name: '%s' in your account", name)
	log.Printf("%v", res.Policies)

	for _, alertObj := range res.Policies {
		if name == alertObj.Name && alertObj.Type == "alert" {
			log.Printf("[DEBUG] Found the Alert Policy")
			log.Printf("[DEBUG] Alert ID: '%s'", alertObj.Id)
			d.Set("name", alertObj.Name)
			d.Set("team_id", team_id)
			d.SetId(alertObj.Id)
			break
		}
	}

	return nil
}
