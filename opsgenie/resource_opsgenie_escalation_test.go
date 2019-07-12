package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/escalation"
	"log"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("opsgenie_escalation", &resource.Sweeper{
		Name: "opsgenie_escalation",
		F:    testSweepEscalation,
	})

}

func testSweepEscalation(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := escalation.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background())
	if err != nil {
		return err
	}

	for _, u := range resp.Escalations {
		if strings.HasPrefix(u.Description, "genieescalation-") {
			log.Printf("Destroying escalation %s", u.Description)

			deleteRequest := escalation.DeleteRequest{
				IdentifierType: escalation.Id,
				Identifier:     u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieEscalation_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieEscalation_basic(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieEscalationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieEscalationExists("opsgenie_escalation.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieEscalation_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieEscalation_complete(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieEscalationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieEscalationExists("opsgenie_escalation.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieEscalationDestroy(s *terraform.State) error {
	client, err := escalation.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_escalation" {
			continue
		}
		Id := rs.Primary.Attributes["id"]

		req := escalation.GetRequest{
			IdentifierType: escalation.Id,
			Identifier:     Id,
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Escalation still exists : %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieEscalationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := escalation.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		Id := rs.Primary.Attributes["id"]

		req := escalation.GetRequest{
			IdentifierType: escalation.Id,
			Identifier:     Id,
		}

		_, err = client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Escalation with id %q does not exist", Id)
		}
		return nil
	}
}

func testAccOpsGenieEscalation_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_escalation" "test" {
 name ="genieescalation-%s"
 rules {
  condition =   "if-not-acked"
    notify_type  =   "default"
    recipient {
      type  = "user"
      id  = "${opsgenie_user.test.id}"
		}
    delay = 1
	}
}
`, rString)
}

func testAccOpsGenieEscalation_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_team" "test" {
  name        = "genieteam"
  description = "This team deals with all the things"
}
resource "opsgenie_schedule" "test" {
  name = "genieschedule"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
}
resource "opsgenie_escalation" "test" {
 name ="genieescalation-%s"
 description="test"
 rules {
  condition =   "if-not-acked"
    notify_type  =   "default"
    recipient {
      type  = "user"
      id  = "${opsgenie_user.test.id}"
    }   
	recipient {
      type  = "team"
      id  = "${opsgenie_team.test.id}"
    }   
	recipient {
      type  = "schedule"
      id  = "${opsgenie_schedule.test.id}"
    }   
    delay = 1
 }
owner_team_id = "${opsgenie_team.test.id}"
repeat  {
  wait_interval = 10
  count = 1
  reset_recipient_states = true
  close_alert_after_all = false
  }
}
`, rString)
}
