package opsgenie

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {

	p := &schema.Provider{
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
			"api_retry_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"api_retry_wait_min": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"api_retry_wait_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
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
			"opsgenie_saved_search":          resourceOpsgenieSavedSearch(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opsgenie_team":       dataSourceOpsGenieTeam(),
			"opsgenie_user":       dataSourceOpsGenieUser(),
			"opsgenie_escalation": dataSourceOpsgenieEscalation(),
			"opsgenie_schedule":   dataSourceOpsgenieSchedule(),
			"opsgenie_heartbeat":  dataSourceOpsgenieHeartbeat(),
			"opsgenie_service":    dataSourceOpsGenieService(),
		},
	}
	p.ConfigureContextFunc = providerConfigure

	return p

}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Println("[INFO] Initializing OpsGenie client")

	config := Config{
		ApiKey:          data.Get("api_key").(string),
		ApiUrl:          data.Get("api_url").(string),
		ApiRetryCount:   data.Get("api_retry_count").(int),
		ApiRetryWaitMin: data.Get("api_retry_wait_min").(int),
		ApiRetryWaitMax: data.Get("api_retry_wait_max").(int),
	}
	cli, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return cli, nil
}
