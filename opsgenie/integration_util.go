package opsgenie

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/integration"
)

func expandOpsgenieIntegrationResponders(d *schema.ResourceData) []integration.Responder {
	input := d.Get("responders").([]interface{})
	responders := make([]integration.Responder, 0, len(input))

	if input == nil {
		return responders
	}

	for _, v := range input {
		config := v.(map[string]interface{})
		responderID := config["id"].(string)
		responder := integration.Responder{
			Type: integration.ResponderType(config["type"].(string)),
			Id:   responderID,
		}

		responders = append(responders, responder)
	}

	return responders
}

func flattenIntegrationResponders(r []interface{}) []map[string]interface{} {
	responders := []map[string]interface{}{}
	for _, i := range r {
		c := i.(map[string]interface{})
		responders = append(responders, map[string]interface{}{
			"type": c["type"],
			"id":   c["id"],
		})
	}
	return responders
}

func validateResponderType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"user":       true,
		"team":       true,
		"escalation": true,
		"schedule":   true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("it can only be one of these 'user', 'schedule', 'team', 'escalation'"))
	}
	return
}

const (
	ApiIntegrationType     = "API"
	EmailIntegrationType   = "Email"
	WebhookIntegrationType = "Webhook"
)
