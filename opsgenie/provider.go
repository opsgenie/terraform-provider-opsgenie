package opsgenie

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSGENIE_API_URL", "api.opsgenie.com"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"opsgenie_custom_role":           resourceOpsGenieCustomUserRole(),
			"opsgenie_team":                  resourceOpsGenieTeam(),
			"opsgenie_team_routing_rule":     resourceOpsGenieTeamRoutingRule(),
			"opsgenie_user":                  resourceOpsGenieUser(),
			"opsgenie_user_contact":          resourceOpsGenieUserContact(),
			"opsgenie_notification_policy":   resourceOpsGenieNotificationPolicy(),
			"opsgenie_notification_rule":     resourceOpsGenieNotificationRule(),
			"opsgenie_escalation":            resourceOpsgenieEscalation(),
			"opsgenie_api_integration":       resourceOpsgenieApiIntegration(),
			"opsgenie_email_integration":     resourceOpsgenieEmailIntegration(),
			"opsgenie_integration_action":    resourceOpsgenieIntegrationAction(),
			"opsgenie_service":               resourceOpsGenieService(),
			"opsgenie_schedule":              resourceOpsgenieSchedule(),
			"opsgenie_schedule_rotation":     resourceOpsgenieScheduleRotation(),
			"opsgenie_maintenance":           resourceOpsgenieMaintenance(),
			"opsgenie_heartbeat":             resourceOpsgenieHeartbeat(),
			"opsgenie_alert_policy":          resourceOpsGenieAlertPolicy(),
			"opsgenie_service_incident_rule": resourceOpsGenieServiceIncidentRule(),
			"opsgenie_incident_template":     resourceOpsgenieIncidentTemplate(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opsgenie_team":       dataSourceOpsGenieTeam(),
			"opsgenie_user":       dataSourceOpsGenieUser(),
			"opsgenie_escalation": dataSourceOpsgenieEscalation(),
			"opsgenie_schedule":   dataSourceOpsgenieSchedule(),
			"opsgenie_heartbeat":  dataSourceOpsgenieHeartbeat(),
			"opsgenie_service":    dataSourceOpsGenieService(),
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
