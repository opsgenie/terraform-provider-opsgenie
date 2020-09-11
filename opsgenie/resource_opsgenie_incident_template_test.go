package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/incident"
)

var randomName = acctest.RandString(6)

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
			incidentName := "Incident Template Name-" + randomName
			if strings.HasPrefix(value.Name, incidentName) {
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
	config := testAccOpsGenieIncidentTemplate_basic(randomName)
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieIncidentTemplateDestroy,
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
			if result != nil {
				for _, value := range result.IncidentTemplates["incidentTemplates"] {
					incidentName := "Incident Template Name-" + randomName
					if strings.HasPrefix(value.Name, incidentName) {
						x := err.(*ogClient.ApiError)
						if x.StatusCode != 404 {
							return errors.New(fmt.Sprintf("Incident template still exists: %s", x.Error()))
						}
					}
				}
			}
		} else {
			return err
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
		if err != nil {
			if result != nil {
				for _, value := range result.IncidentTemplates["incidentTemplates"] {
					incidentName := "Incident Template Name-" + randomName
					if strings.HasPrefix(value.Name, incidentName) {
						log.Printf("Incident template found.")
					} else {
						return fmt.Errorf("incident template does not exist (and it should)")
					}
				}
			} else {
				return fmt.Errorf("incident template does not exist (and it should)")
			}
		} else {
			return err
		}
		return nil
	}
}

func testAccOpsGenieIncidentTemplate_basic(randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "gtest_1-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_team" "test" {
  name        = "gtest_1-%s"
  description = "This team deals with all the things"
  member {
    id        = opsgenie_user.test.id
    role      = "admin"
  }
}
resource "opsgenie_service" "test" {
  name  = "gtest_1-%s"
  team_id = opsgenie_team.test.id
}
resource "opsgenie_incident_template" "test" {
  name = "Incident Template Name-%s"
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
    opsgenie_service.test.id, 
    "75fa8b15-d9a2-4c68-8d45-53b6540d0d09"
  ]
}
`, randomName, randomName, randomName, randomName)
}
