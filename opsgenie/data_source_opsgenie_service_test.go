package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceOpsGenieService_Basic(t *testing.T) {
	randomUserName := acctest.RandString(6)
	randomTeamName := acctest.RandString(6)
	randomServiceName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieServiceConfig(randomUserName, randomTeamName, randomServiceName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieService("opsgenie_service.test", "data.opsgenie_service.existingservice"),
				),
			},
		},
	})
}

func testAccDataSourceOpsGenieService(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a team ID from OpsGenie")
		}

		testAtts := []string{"name", "team_id", "description"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the team %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieServiceConfig(randomUserName, randomTeamName, randomServiceName string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_user" "test2" {
  username  = "genietest-2%s@opsgenie.com"
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
  member {
    id = "${opsgenie_user.test2.id}"
    role     = "admin"
  }
}
resource "opsgenie_service" "test" {
  name = "genieservice-%s"
  description = "This is our main service"
  team_id = "${opsgenie_team.test.id}"
}
data "opsgenie_service" "existingservice" {
  name = "${opsgenie_service.test.name}"
}
`, randomUserName, randomUserName, randomTeamName, randomServiceName)
}
