package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func init() {
	resource.AddTestSweepers("opsgenie_team_role", &resource.Sweeper{
		Name: "opsgenie_team_role",
		F:    testSweepTeamRole,
	})
}

func testSweepTeamRole(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.ListRole(context.Background(), &team.ListTeamRoleRequest{})
	if err != nil {
		return err
	}

	for _, teamRole := range resp.TeamRoles {
		if strings.HasPrefix(teamRole.Name, "genietest-") {
			log.Printf("Destroying team role %s", teamRole.Name)

			deleteRequest := team.DeleteTeamRoleRequest{
				TeamID:   strconv.Itoa(int(team.Id)),
				RoleID:   teamRole.Id,
				RoleName: teamRole.Name,
			}

			if _, err := client.DeleteRole(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieTeamRole_basic(t *testing.T) {
	teamName := acctest.RandString(6)
	roleName := acctest.RandString(6)
	config := testAccOpsGenieTeamRole_basic(teamName, roleName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieTeamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamRoleExists("opsgenie_team_role.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieTeamRole_complete(t *testing.T) {
	teamName := acctest.RandString(6)
	roleName := acctest.RandString(6)
	config := testAccOpsGenieTeamRole_complete(teamName, roleName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieTeamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieTeamRoleExists("opsgenie_team_role.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieTeamRoleDestroy(s *terraform.State) error {
	client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_team_role" {
			continue
		}
		req := team.GetTeamRoleRequest{
			TeamID:   strconv.Itoa(int(team.Id)),
			RoleID:   rs.Primary.Attributes["role_id"],
			RoleName: rs.Primary.Attributes["role_name"],
		}
		_, err := client.GetRole(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Team role still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieTeamRoleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]
		teamRoleName := rs.Primary.Attributes["role_name"]

		client, err := team.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := team.GetTeamRoleRequest{
			TeamID:   strconv.Itoa(int(team.Id)),
			RoleID:   rs.Primary.Attributes["role_id"],
			RoleName: rs.Primary.Attributes["role_name"],
		}

		result, err := client.GetRole(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: teamrole %q (teamRoleName: %q) does not exist", id, teamRoleName)
		} else {
			log.Printf("Team role found :%s ", result.Name)
		}

		return nil
	}
}

func TestAccOpsGenieTeamRole_rightsValidationError(t *testing.T) {
	teamName := acctest.RandString(6)
	roleName := acctest.RandString(6)
	config := testAccOpsGenieTeamRole_rightsValidationError(teamName, roleName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				ExpectError: regexp.MustCompile(`config is invalid: expected rights.0 to be one of \[manage-members
	edit-team-roles delete-team-roles access-member-profiles edit-member-profiles edit-routing-rules delete-routing-rules edit-escalations delete-escalations edit-schedules delete-schedules edit-integrations delete-integrations edit-heartbeats delete-heartbeats access-reports edit-services delete-services edit-rooms delete-rooms send-service-status-update",\], got invalid-right`),
			},
		},
	})
}

func testAccOpsGenieTeamRole_basic(teamName string, roleName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}

resource "opsgenie_team_role" "test" {
  team_id  = "${opsgenie_team.test.id}"
  role_name  = "opsgenie-%s"
  rights = [
    {
      "right": "access-member-profiles",
      "granted": true
    }
  ]
}
`, teamName, roleName)
}

func testAccOpsGenieTeamRole_complete(teamName string, roleName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}

resource "opsgenie_team_role" "test" {
  team_id  = "${opsgenie_team.test.id}"
  role_name  = "opsgenie-%s"
  rights = [
    {
      "right": "access-member-profiles",
      "granted": true
    }
  ]
}
`, teamName, roleName)
}

func testAccOpsGenieTeamRole_rightsValidationError(teamName string, roleName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}

resource "opsgenie_team_role" "test" {
  team_id  = "${opsgenie_team.test.id}"
  role_name  = "opsgenie-%s"
  rights = [
    {
      "right": "access-member",
      "granted": true
    }
  ]
}
`, teamName, roleName)
}
