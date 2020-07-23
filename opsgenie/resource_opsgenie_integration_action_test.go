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
	"github.com/opsgenie/opsgenie-go-sdk-v2/integration"
)

func init() {
	resource.AddTestSweepers("opsgenie_integration_action", &resource.Sweeper{
		Name: "opsgenie_integration_action",
		F:    testSweepIntegrationAction,
	})

}

func testSweepIntegrationAction(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background())
	if err != nil {
		return err
	}

	for _, u := range resp.Integrations {
		if strings.HasPrefix(u.Name, "genieintegration-") {
			log.Printf("Destroying integration actions for id: %s", u.Name)

			deleteRequest := integration.UpdateAllIntegrationActionsRequest{
				Id:          u.Id,
				Create:      []integration.IntegrationAction{},
				Close:       []integration.IntegrationAction{},
				Acknowledge: []integration.IntegrationAction{},
				AddNote:     []integration.IntegrationAction{},
			}

			if _, err := client.UpdateAllActions(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieIntegrationAction_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieIntegrationAction_basic(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieIntegrationActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieIntegrationActionExists("opsgenie_integration_action.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieIntegrationAction_complete(t *testing.T) {
	rString := acctest.RandString(6)

	config := testAccOpsGenieIntegrationAction_complete(rString)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieIntegrationActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieIntegrationActionExists("opsgenie_integration_action.test_email"),
					testCheckOpsGenieIntegrationActionExists("opsgenie_integration_action.test_api"),
				),
			},
		},
	})
}

func testCheckOpsGenieIntegrationActionDestroy(s *terraform.State) error {
	client, err := integration.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_integration_action" {
			continue
		}
		Id := rs.Primary.Attributes["id"]

		req := integration.GetIntegrationActionsRequest{
			Id: Id,
		}
		_, err := client.GetActions(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			log.Print(x)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Api Integration still exists: %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieIntegrationActionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := integration.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		Id := rs.Primary.Attributes["id"]

		req := integration.GetIntegrationActionsRequest{
			Id: Id,
		}

		_, err = client.GetActions(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: ApiIntegration with id %q does not exist", Id)
		}
		return nil
	}
}

func testAccOpsGenieIntegrationAction_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_integration_action" "test" {
  integration_id = opsgenie_api_integration.test.id
  close {
    name = "Test close action"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
    }
  }
}
resource "opsgenie_api_integration" "test" {
  type = "API"
  name = "genieintegration-%s"
}
`, rString)
}

func testAccOpsGenieIntegrationAction_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_email_integration" "test" {
  name = "genieintegration-email-%s"
  email_username="example"
  ignore_responders_from_payload = true
  suppress_notifications = true
}
resource "opsgenie_api_integration" "test" {
  name = "genieintegration-api-%s"
  type = "API"
  allow_write_access = false
  ignore_responders_from_payload = true
  suppress_notifications = true
}
resource "opsgenie_integration_action" "test_email" {
  integration_id = opsgenie_email_integration.test.id
  create {
    name = "Filter high prio alerts"
    filter {
      type = "match-any-condition"
      conditions {
        field = "source"
        operation = "equals"
        expected_value = "notifier@opsgenie.com"
      }
      conditions {
        field = "source"
        operation = "equals"
        expected_value = "alert@opsgenie.com"
      }
    }
  }
}
resource "opsgenie_integration_action" "test_api" {
  integration_id = opsgenie_api_integration.test.id
  create {
    name = "Filter high prio alerts"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "greater_than"
        expected_value = "P2"
      }
      conditions {
        field = "message"
        operation = "contains"
        expected_value = "critical"
      }
    }
  }
  acknowledge {
    name = "Ack P5 alerts"
    filter {
      type = "match-any-condition"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
    }
  }
}
`, rString, rString)
}
