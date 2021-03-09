package opsgenie

import (
	"context"

	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpsgenieSchedule() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOpsgenieScheduleRead,
		Schema: opsgenieScheduleSchema,
	}
}

func dataSourceOpsgenieScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	scheduleName := d.Get("name").(string)

	getRequest := &schedule.GetRequest{
		IdentifierType:  schedule.Name,
		IdentifierValue: scheduleName,
	}

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.SetId(getResponse.Schedule.Id)
	d.Set("name", getResponse.Schedule.Name)
	if getResponse.Schedule.OwnerTeam != nil {
		d.Set("owner_team_id", getResponse.Schedule.OwnerTeam.Id)
	}
	d.Set("description", getResponse.Schedule.Description)
	d.Set("timezone", getResponse.Schedule.Timezone)
	d.Set("enabled", getResponse.Schedule.Enabled)

	return nil
}
