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
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func init() {
	resource.AddTestSweepers("opsgenie_team_routing_rule", &resource.Sweeper{
		Name: "opsgenie_team_routing_roule",
		F:    testSweepTeamRoutingRule,
	})

}

func testSweepTeamRoutingRule(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &team.ListTeamRequest{})
	if err != nil {
		return err
	}

	for _, u := range resp.Teams {
		if strings.HasPrefix(u.Name, "genieteam") {
			log.Printf("Destroying team routing rule %s", u.Name)
			resp2, err := client.ListRoutingRules(context.Background(), &team.ListRoutingRulesRequest{
				TeamIdentifierType:  team.Id,
				TeamIdentifierValue: u.Id,
			})
			if err != nil {
				return err
			}
			for _, k := range resp2.RoutingRules {

				deleteRequest := team.DeleteRoutingRuleRequest{
					TeamIdentifierType:  team.Id,
					TeamIdentifierValue: u.Id,
					RoutingRuleId:       k.Id,
				}

				if _, err := client.DeleteRoutingRule(context.Background(), &deleteRequest); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func TestAccOpsGenieTeamRoutingRule_basic(t *testing.T) {
	teamName := acctest.RandString(6)
	scheduleName := acctest.RandString(6)
	routeRuleName := acctest.RandString(6)

	config := testAccOpsGenieTeamRoutingRule_basic(scheduleName, teamName, routeRuleName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieTeamRoutingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamRoutingRuleExists("opsgenie_team_routing_rule.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieTeamRoutingRuleDestroy(s *terraform.State) error {
	client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_team_routing_rule" {
			continue
		}

		req := team.GetRoutingRuleRequest{
			TeamIdentifierType:  team.Id,
			TeamIdentifierValue: rs.Primary.Attributes["team_id"],
			RoutingRuleId:       rs.Primary.Attributes["id"],
		}
		_, err := client.GetRoutingRule(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Team routing rule still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieTeamRoutingRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]

		client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := team.GetRoutingRuleRequest{
			TeamIdentifierType:  team.Id,
			TeamIdentifierValue: rs.Primary.Attributes["team_id"],
			RoutingRuleId:       id,
		}

		_, err = client.GetRoutingRule(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Team routing rule %q  does not exist", id)
		}
		return nil
	}
}

func testAccOpsGenieTeamRoutingRule_basic(scheduleName, teamName, routingRuleName string) string {
	return fmt.Sprintf(`
resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
  owner_team_id = "${opsgenie_team.test.id}"
}

resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}

resource "opsgenie_team_routing_rule" "test" {
  name     = "genieteam-%s"
  team_id  = "${opsgenie_team.test.id}"
  order    = 0
  timezone = "America/Los_Angeles"
  criteria {
    type = "match-all-conditions"
    conditions {
      field          = "message"
      operation      = "contains"
      expected_value = "unexpected1"
      not            = true
    }
    conditions {
      field          = "message"
      operation      = "contains"
      expected_value = "expected1"
      not            = false
    }
  }
  time_restriction {
    type = "weekday-and-time-of-day"
    restrictions {
      start_day  = "monday"
      start_hour = 8
      start_min  = 0
      end_day    = "tuesday"
      end_hour   = 18
      end_min    = 30
    }
  }
  notify {
    name = "${opsgenie_schedule.test.name}"
    type = "schedule"
  }
}
`, scheduleName, teamName, routingRuleName)
}
