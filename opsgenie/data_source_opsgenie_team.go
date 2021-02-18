package opsgenie

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func dataSourceOpsGenieTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieTeamRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieTeamName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user",
						},
					},
				},
			},
		},
	}
}

func dataSourceOpsGenieTeamRead(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	teamName := d.Get("name").(string)

	getRequest := &team.GetTeamRequest{
		IdentifierType:  team.Name,
		IdentifierValue: teamName,
	}

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}
	d.SetId(getResponse.Id)

	d.Set("description", getResponse.Description)
	d.Set("member", flattenOpsGenieTeamMembers(getResponse.Members))

	return nil
}
