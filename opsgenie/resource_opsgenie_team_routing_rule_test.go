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
	rs := acctest.RandString(6)
	config := testAccOpsGenieTeamRoutingRule_basic(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieTeamRoutingRuleDestroy,
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

func testAccOpsGenieTeamRoutingRule_basic(rString string) string {
	return fmt.Sprintf(`
	resource "opsgenie_team_routing_rule" "test" {
	name        = "genieteam-%s"
	team_id = "a9b12343-c2c1-4d3f-8412-d715f7a2fc28"
	order = 0
	timezone = "America/Los_Angeles"
	criteria  {
	type = "match-any-condition"
	conditions {
	field = "message"
	operation = "contains"
	expected_value = "expected1"
    not = false

    }
	}
	time_restriction {
	type = "weekday-and-time-of-day"
    restrictions {
	start_day="monday"
	start_hour = 8
	start_min = 0
	end_day = "tuesday"
	end_hour = 18
	end_min = 30
	}
	}
 notify {
        name="uuu-test-team-ops_schedule"
        type="schedule"
    }

}

`, rString)
}
