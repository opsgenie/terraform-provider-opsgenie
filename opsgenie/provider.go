package opsgenie

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSGENIE_API_KEY", nil),
			},
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "api.opsgenie.com",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"opsgenie_team":              resourceOpsGenieTeam(),
			"opsgenie_user":              resourceOpsGenieUser(),
			"opsgenie_user_contact":      resourceOpsGenieUserContact(),
			"opsgenie_escalation":        resourceOpsgenieEscalation(),
			"opsgenie_api_integration":   resourceOpsgenieApiIntegration(),
			"opsgenie_email_integration": resourceOpsgenieEmailIntegration(),
			"opsgenie_schedule":          resourceOpsgenieSchedule(),
			"opsgenie_schedule_rotation": resourceOpsgenieScheduleRotation(),
			"opsgenie_maintenance":       resourceOpsgenieMaintenance(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing OpsGenie client")

	config := Config{
		ApiKey: data.Get("api_key").(string),
		ApiUrl: data.Get("api_url").(string),
	}

	return config.Client()
}
