package opsgenie

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOpsGenieTeamMembership_basic(t *testing.T) {
	rString := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config:  testAccOpsGenieTeamMembership_basic(rString),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamExists("opsgenie_team.monkeys"),
					testCheckOpsGenieUserExists("opsgenie_user.kong"),
					testCheckOpsGenieTeamMembershipExists("opsgenie_team_membership.chaos_kong", "opsgenie_team.monkeys", "opsgenie_user.kong"),
				),
			},
			{
				Config: testAccOpsGenieTeamMembership_basicUpdated(rString),
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamExists("opsgenie_team.monkeys"),
					testCheckOpsGenieUserExists("opsgenie_user.kong"),
					testCheckOpsGenieTeamMembershipExists("opsgenie_team_membership.chaos_kong", "opsgenie_team.monkeys", "opsgenie_user.kong"),
				),
			},
			{
				Config: testAccOpsGenieTeamMembership_basicWithoutMembership(rString),
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamExists("opsgenie_team.monkeys"),
					testCheckOpsGenieUserExists("opsgenie_user.kong"),
					testCheckOpsGenieTeamMembershipRemoved("opsgenie_team_membership.chaos_kong", "opsgenie_team.monkeys", "opsgenie_user.kong"),
				),
			},
		},
	})
}

func testCheckOpsGenieTeamMembershipExists(membershipResource string, teamResource string, userResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rsMembership, ok := s.RootModule().Resources[membershipResource]
		if !ok {
			return fmt.Errorf("not found: %s", membershipResource)
		}

		rsTeam, ok := s.RootModule().Resources[teamResource]
		if !ok {
			return fmt.Errorf("not found: %s", teamResource)
		}
		teamName := rsTeam.Primary.Attributes["name"]

		rsUser, ok := s.RootModule().Resources[userResource]
		if !ok {
			return fmt.Errorf("not found: %s", userResource)
		}
		userName := rsUser.Primary.Attributes["username"]

		client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := team.GetTeamRequest{
			IdentifierType:  team.Name,
			IdentifierValue: teamName,
		}
		getResponse, err := client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("failed to detect team membership for user %q in team %q: %s", userName, teamName, err)
		}

		// compare what we've actually done
		if len(getResponse.Members) != 1 {
			return fmt.Errorf("there's no team membership at all. something went wrong :(")
		}

		if getResponse.Members[0].User.Username != userName {
			return fmt.Errorf("expected userName in team membership (%q) doesn't match actual username (%q)", userName, getResponse.Members[0].User.Username)
		}

		if getResponse.Members[0].User.ID != rsUser.Primary.ID {
			return fmt.Errorf("expected user ID in team membership (%q) doesn't match actual username (%q)", rsUser.Primary.ID, getResponse.Members[0].User.ID)
		}

		if getResponse.Members[0].Role != rsMembership.Primary.Attributes["role"] {
			return fmt.Errorf("expected user role in team membership (%q) doesn't match actual user role (%q)", rsMembership.Primary.Attributes["role"], getResponse.Members[0].Role)
		}

		return nil
	}
}

func testCheckOpsGenieTeamMembershipRemoved(membershipResource string, teamResource string, userResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		_, ok := s.RootModule().Resources[membershipResource]
		if ok {
			return fmt.Errorf("resource %s still in state. this is bad", membershipResource)
		}

		rsTeam, ok := s.RootModule().Resources[teamResource]
		if !ok {
			return fmt.Errorf("not found: %s", teamResource)
		}
		teamName := rsTeam.Primary.Attributes["name"]

		client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := team.GetTeamRequest{
			IdentifierType:  team.Name,
			IdentifierValue: teamName,
		}
		getResponse, err := client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("failed to verify team memberships of team %q: %s", teamName, err)
		}

		// compare what we've actually done
		if len(getResponse.Members) != 0 {
			return fmt.Errorf("there is still an unexpected number of team membership(s) (%#v). something went wrong :(", getResponse.Members)
		}

		return nil
	}
}

func testAccOpsGenieTeamMembership_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "monkeys" {
  name           = "monkeys-%s"
  description    = "They exist."
  ignore_members = true
}
resource "opsgenie_user" "kong" {
  username   = "kong-%s@inovex.de"
  full_name  = "Chaos Kong"
  role       = "User"
}
resource "opsgenie_team_membership" "chaos_kong" {
  username   = opsgenie_user.kong.username
  role       = "user"
  team       = opsgenie_team.monkeys.name
}
`, rString, rString)
}

func testAccOpsGenieTeamMembership_basicUpdated(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "monkeys" {
  name           = "monkeys-%s"
  description    = "They exist."
  ignore_members = true
}
resource "opsgenie_user" "kong" {
  username   = "kong-%s@test.example.com"
  full_name  = "Chaos Kong"
  role       = "User"
}
resource "opsgenie_team_membership" "chaos_kong" {
  username   = opsgenie_user.kong.username
  role       = "admin"
  team       = opsgenie_team.monkeys.name
}
`, rString, rString)
}

func testAccOpsGenieTeamMembership_basicWithoutMembership(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "monkeys" {
  name           = "monkeys-%s"
  description    = "They exist."
  ignore_members = true
  depends_on     = [opsgenie_user.kong] # Just a hack for the test to destroy resources in the right order
}
resource "opsgenie_user" "kong" {
  username   = "kong-%s@test.example.com"
  full_name  = "Chaos Kong"
  role       = "User"
}
`, rString, rString)
}
