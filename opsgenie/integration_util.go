package opsgenie

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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

func validateOpsgenieIntegrationName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_- ]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alpha numeric characters, underscores, and spaces are allowed in %q: %q", k, value))
	}

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
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
	ApiIntegrationType   = "API"
	EmailIntegrationType = "Email"
)
