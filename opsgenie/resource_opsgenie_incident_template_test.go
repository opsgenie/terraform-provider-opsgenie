package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/incident"
)

func init() {
	resource.AddTestSweepers("opsgenie_incident_template", &resource.Sweeper{
		Name: "opsgenie_incident_template",
		F:    testSweepIncidentTemplate,
	})
}

func testSweepIncidentTemplate(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}
	client, err := incident.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	result, err := client.GetIncidentTemplate(context.Background(), &incident.GetIncidentTemplateRequest{})
	if err != nil {
		return err
	}
	if result != nil {
		for _, value := range result.IncidentTemplates["incidentTemplates"] {
			if strings.HasPrefix(value.Name, "genietest-incident-template-") {
				log.Printf("Destroying incident template %s", value.Name)
				deleteRequest := incident.DeleteIncidentTemplateRequest{IncidentTemplateId: value.IncidentTemplateId}
				if _, err := client.DeleteIncidentTemplate(context.Background(), &deleteRequest); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func TestAccOpsGenieIncidentTemplate_basic(t *testing.T) {
	config := testAccOpsGenieIncidentTemplate_basic(acctest.RandString(6))
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieIncidentTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieIncidentTemplateExists(),
				),
			},
		},
	})
}

func testCheckOpsGenieIncidentTemplateDestroy(s *terraform.State) error {
	client, err := incident.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_incident_template" {
			continue
		}
		result, err := client.GetIncidentTemplate(context.Background(), &incident.GetIncidentTemplateRequest{})
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Incident template still exists: %s", x.Error()))
			}
		} else if result != nil {
			for _, value := range result.IncidentTemplates["incidentTemplates"] {
				if strings.HasPrefix(value.Name, "genietest-incident-template-") {
					return fmt.Errorf("incident template still exists(it shouldn't exist)")
				}
			}
		}
	}
	return nil
}

func testCheckOpsGenieIncidentTemplateExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := incident.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		result, err := client.GetIncidentTemplate(context.Background(), &incident.GetIncidentTemplateRequest{})
		if err != nil && result != nil {
			for _, value := range result.IncidentTemplates["incidentTemplates"] {
				if strings.HasPrefix(value.Name, "genietest-incident-template-") {
					log.Printf("Incident template found.")
					return nil
				}
			}
			return fmt.Errorf("incident template does not exist (and it should)")
		} else {
			return err
		}
	}
}

func testAccOpsGenieIncidentTemplate_basic(randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genietest-team-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test" {
  name  = "genietest-service-%s"
  team_id = opsgenie_team.test.id
}
resource "opsgenie_incident_template" "test" {
  name = "genietest-incident-template-%s"
  message = "Incident Message"
  priority = "P2"
  stakeholder_properties {
    enable = true
    message = "Stakeholder Message"
    description = "Stakeholder Description"
  }
  tags = ["tag1", "tag2"]
  description = "Incident Description"
  details = {
    key1 = "value1"
    key2 = "value2"
  }
  impacted_services = [ 
    opsgenie_service.test.id
  ]
}
`, randomName, randomName, randomName)
}
