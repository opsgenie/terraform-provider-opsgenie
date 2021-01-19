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
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieIntegrationActionDestroy,
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

func TestAccOpsGenieIntegrationAction_custompriority(t *testing.T) {
	customPriority := "{{condition_name.extract(/^\\[(\\S+)\\].*$/, 1)}"
	customPriorityEscaped := "{{condition_name.extract(/^\\\\[(\\\\S+)\\\\].*$/, 1)}"
	rs := acctest.RandString(6)
	config := testAccOpsGenieIntegrationAction_custompriority(rs, customPriorityEscaped)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieIntegrationActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieIntegrationActionCustomPriorityExists("opsgenie_integration_action.testcustom", customPriority),
				),
			},
		},
	})
}

func TestAccOpsGenieIntegrationAction_complete(t *testing.T) {
	rString := acctest.RandString(6)

	config := testAccOpsGenieIntegrationAction_complete(rString)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieIntegrationActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
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

func testCheckOpsGenieIntegrationActionCustomPriorityExists(name, customPriority string) resource.TestCheckFunc {
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

		iAction, err := client.GetActions(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: ApiIntegration with id %q does not exist", Id)
		}
		if iAction.Create[0].CustomPriority != customPriority {
			return fmt.Errorf("Bad: CustomPriority %s for ApiIntegration with id %q does not exist", customPriority, Id)
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

func testAccOpsGenieIntegrationAction_custompriority(rString, crString string) string {
	return fmt.Sprintf(`
resource "opsgenie_api_integration" "test" {
  type = "API"
  name = "genieintegration-custompri-%s"
}
resource "opsgenie_integration_action" "testcustom" {
  integration_id = opsgenie_api_integration.test.id
  create {
    name = "Create high priority alerts"
    tags = ["CRITICAL", "SEV-0"]
    user = "Acceptance test user"
    note = "{{note}}"
    alias = "{{alias}}"
    source = "{{source}}"
    message = "{{message}}"
    description = "{{description}}"
    entity = "{{entity}}"
    custom_priority = "%s"
    alert_actions = ["Check error rate"]
    extra_properties = map(
      "Environment", "test-env",
      "Region", "us-west-2"
    )
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P1"
      }
    }
  }
}
`, rString, crString)
}

func testAccOpsGenieIntegrationAction_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "genietest-%s"
  role      = "User"
}
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_api_integration" "test" {
  name = "genieintegration-api-%s"
  type = "API"
  allow_write_access = false
  ignore_responders_from_payload = true
  suppress_notifications = true
}
resource "opsgenie_email_integration" "test" {
  name = "geniemailintegration-%s"
  email_username = "geniemailintegration-%s"
}
resource "opsgenie_integration_action" "test_api" {
  integration_id = "${opsgenie_api_integration.test.id}"
  create {
    name = "Create high priority alerts"
    tags = ["CRITICAL", "SEV-0"]
    user = "Acceptance test user"
    note = "{{note}}"
	alias = "{{alias}}"
	source = "{{source}}"
	message = "{{message}}"
	description = "{{description}}"
	entity = "{{entity}}"
    priority = "P2"
	alert_actions = ["Check error rate"]
    extra_properties = map(
      "Environment", "test-env",
      "Region", "us-west-2"
    )
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P1"
      }
    }
    responders {
      id = "${opsgenie_team.test.id}"
      type = "team"
    }
  }
  create {
    name = "Create medium priority alerts"
	order = 2
    tags = ["SEVERE", "SEV-1"]
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P2"
      }
    }
    responders {
      id = "${opsgenie_user.test.id}"
      type = "user"
    }
  }
  ignore {
    name = "Ignore alerts with priority P5"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
    }
  }
  close {
    name = "Low priority alerts"
    filter {
      type = "match-any-condition"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
    }
  }
  acknowledge {
    name = "Auto-ack test alerts"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "message"
        operation = "contains"
        expected_value = "TEST"
      }
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
    }
  }
  add_note {
    name = "Add note to all alerts"
    note = "Created from test integration"
    filter {
      type = "match-all"
    }
  }
}
resource "opsgenie_integration_action" "test_email" {
  integration_id = "${opsgenie_email_integration.test.id}"
  create {
    name = "Accept from source"
    filter {
      type = "match-any-condition"
      conditions {
        field = "from_address"
        operation = "equals"
        expected_value = "alerts@opsgenie.com"
      }
      conditions {
        field = "from_name"
        operation = "equals"
        expected_value = "admin"
      }
      conditions {
        field = "subject"
        operation = "contains"
        expected_value = "S1:"
      }
    }
  }
}
`, rString, rString, rString, rString, rString, rString)
}
