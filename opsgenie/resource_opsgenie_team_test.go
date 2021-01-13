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
	resource.AddTestSweepers("opsgenie_team", &resource.Sweeper{
		Name: "opsgenie_team",
		F:    testSweepTeam,
	})

}

func testSweepTeam(region string) error {
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

func TestAccOpsGenieTeam_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieTeam_basic(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamExists("opsgenie_team.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieTeam_basicNoMember(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomUser := acctest.RandString(6)
	config := testAccOpsGenieTeam_basicNoMember(randomUser, randomTeam)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamExists("opsgenie_team.test"),
					testEmptyMemberList("opsgenie_team.test"),
					addTeamMembers("opsgenie_team.test", "opsgenie_user.test"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testEmptyMemberList("opsgenie_team.test"),
					removeAllTeamMembers("opsgenie_team.test"),
				),
			},
		},
	})
}

func testEmptyMemberList(teamName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[teamName]
		if !ok {
			return fmt.Errorf("Not found: %s", teamName)
		}

		if len(teamResource.Primary.Attributes["member"]) > 0 {
			return fmt.Errorf("member list in state of resource %q is unexpectely not empty", teamName)
		}

		return nil
	}
}

func addTeamMembers(teamName string, userName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[teamName]
		if !ok {
			return fmt.Errorf("Not found: %s", teamName)
		}

		userResource, ok := s.RootModule().Resources[userName]
		if !ok {
			return fmt.Errorf("Not found: %s", userName)
		}

		updateRequest := team.UpdateTeamRequest{
			Id:          teamResource.Primary.Attributes["id"],
			Name:        teamResource.Primary.Attributes["name"],
			Description: teamResource.Primary.Attributes["description"],
			Members: []team.Member{
				{
					User: team.User{ID: userResource.Primary.ID},
				},
			},
		}

		return updateTeamMembers(s, updateRequest)
	}
}

func removeAllTeamMembers(teamName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[teamName]
		if !ok {
			return fmt.Errorf("Not found: %s", teamName)
		}

		updateRequest := team.UpdateTeamRequest{
			Id:          teamResource.Primary.Attributes["id"],
			Name:        teamResource.Primary.Attributes["name"],
			Description: teamResource.Primary.Attributes["description"],
			Members:     []team.Member{},
		}

		return updateTeamMembers(s, updateRequest)
	}
}

func updateTeamMembers(s *terraform.State, updateRequest team.UpdateTeamRequest) error {
	client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	_, err = client.Update(context.Background(), &updateRequest)
	if err != nil {
		return fmt.Errorf("error updating team %q with members %q", updateRequest.Name, updateRequest.Members)
	}

	return nil
}

func TestAccOpsGenieTeam_complete(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomUser := acctest.RandString(6)
	config := testAccOpsGenieTeam_complete(randomUser, randomTeam)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamExists("opsgenie_team.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieTeamDestroy(s *terraform.State) error {
	client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_team" {
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

func testCheckOpsGenieTeamExists(name string) resource.TestCheckFunc {
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

func testAccOpsGenieTeam_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
`, rString)
}

func testAccOpsGenieTeam_basicNoMember(randomUser, rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name           = "genieteam-%s"
  description    = "This team deals with all the things"
  ignore_members = true
  depends_on     = [opsgenie_user.test] # Just a hack for the test to destroy resources in the right order
}
resource "opsgenie_user" "test" {
  username   = "genietest-%s@opsgenie.com"
  full_name  = "Acceptance Test User"
  role       = "User"
}
`, randomUser, rString)
}

func testAccOpsGenieTeam_complete(randomUser, randomTeam string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"

  member {
    id = "${opsgenie_user.test.id}"
    role     = "admin"
  }
}
`, randomUser, randomTeam)
}
