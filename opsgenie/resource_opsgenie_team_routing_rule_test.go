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
		F:    testSweepTeam,
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
			log.Printf("Destroying team %s", u.Name)

			deleteRequest := team.DeleteTeamRequest{
				IdentifierType:  team.Id,
				IdentifierValue: u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
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
					testCheckOpsGenieTeamExists("opsgenie_team_routing_rule.test"),
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

		req := team.GetTeamRequest{
			IdentifierType:  team.Id,
			IdentifierValue: rs.Primary.Attributes["id"],
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Team still exists : %s", x.Error()))
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
		teamname := rs.Primary.Attributes["name"]

		client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := team.GetTeamRequest{
			IdentifierType:  team.Id,
			IdentifierValue: rs.Primary.Attributes["id"],
		}

		_, err = client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Team %q (teamname: %q) does not exist", id, teamname)
		}
		return nil
	}
}

func testAccOpsGenieTeamRoutingRule_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_team_routing_rule" "test" {
  name        = "genieteam-%s"
  team_id = "a"
  order = 1
criteria = {
id="id"
role="fahri"
}	
}
`, rString)
}
