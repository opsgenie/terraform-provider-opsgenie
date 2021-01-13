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
	"github.com/opsgenie/opsgenie-go-sdk-v2/service"
)

func init() {
	resource.AddTestSweepers("opsgenie_service_incident_rule", &resource.Sweeper{
		Name: "opsgenie_service_incident_rule",
		F:    testSweepServiceIncidentRule,
	})

}

func testSweepServiceIncidentRule(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &service.ListRequest{})
	if err != nil {
		return err
	}

	for _, svc := range resp.Services {
		if strings.HasPrefix(svc.Name, "genietest-") {
			log.Printf("Destroying service %s", svc.Name)

			deleteRequest := service.DeleteRequest{
				Id: svc.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieServiceIncidentRule_basic(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomService := acctest.RandString(6)

	config := testAccOpsGenieServiceIncidentRule_basic(randomTeam, randomService)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieServiceIncidentRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieServiceIncidentRuleExists("opsgenie_service_incident_rule.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieServiceIncidentRule_complete(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomService := acctest.RandString(6)

	config := testAccOpsGenieServiceIncidentRule_complete(randomTeam, randomService)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieServiceIncidentRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieServiceIncidentRuleExists("opsgenie_service_incident_rule.test2"),
				),
			},
		},
	})
}

func testCheckOpsGenieServiceIncidentRuleDestroy(s *terraform.State) error {
	client, err := service.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_service" {
			continue
		}
		req := service.GetRequest{
			Id: rs.Primary.Attributes["id"],
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Service still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieServiceIncidentRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		service_incident_rule_id := rs.Primary.Attributes["id"]
		service_id := rs.Primary.Attributes["service_id"]

		client, err := service.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		incident_rule_res, err := client.GetIncidentRules(context.Background(), &service.GetIncidentRulesRequest{
			ServiceId: service_id,
		})

		if err != nil {
			return fmt.Errorf("Bad: Service ID %q does not exist", service_id)
		} else {
			for _, v := range incident_rule_res.IncidentRule {
				fmt.Printf("checking service incident rule id %s", v.Id)
				if v.Id == service_incident_rule_id {
					log.Printf("Service incident rule found for service:%s ", v.Id)
					return nil
				}
			}

		}

		return fmt.Errorf("Bad: Service incident rule for service id: %q does not exist", service_id)
	}
}

func testAccOpsGenieServiceIncidentRule_basic(randomTeam, randomService string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test" {
  name  = "genietest-service-%s"
  team_id = opsgenie_team.test.id
}
resource "opsgenie_service_incident_rule" "test" {
  service_id = opsgenie_service.test.id
  incident_rule {
	incident_properties {
		message = "This is a test message"
		priority = "P3"
		stakeholder_properties {
			message = "Message for stakeholders"
		}
	}
  }
}
`, randomTeam, randomService)
}

func testAccOpsGenieServiceIncidentRule_complete(randomTeam, randomService string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test2" {
  name  = "genietest-service-%s"
  team_id = opsgenie_team.test.id
}
resource "opsgenie_service_incident_rule" "test2" {
  service_id = opsgenie_service.test2.id
  incident_rule {
	condition_match_type = "match-any-condition"
	conditions {
		field = "message"
		not =  false
		operation = "contains"
		expected_value = "expected1"
	}
	conditions {
		field = "message"
		not =  false
		operation = "contains"
		expected_value = "expected2"
	}
	incident_properties {
		message = "This is a test message"
		priority = "P3"
        tags = ["a","b"]
        details = {
			custom = "parameter"
        }
		stakeholder_properties {
			message = "Message for stakeholders"
			enable = "true"
		}
	}
  }
}
`, randomTeam, randomService)
}
