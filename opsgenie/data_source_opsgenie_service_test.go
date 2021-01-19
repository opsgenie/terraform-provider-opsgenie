package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceOpsGenieService_Basic(t *testing.T) {
	randomTeamName := acctest.RandString(6)
	randomServiceName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieServiceConfig(randomTeamName, randomServiceName),
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
			return fmt.Errorf("Expected to get a Service ID from OpsGenie")
		}

		testAtts := []string{"name", "team_id", "description"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the Service %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieServiceConfig(randomTeamName, randomServiceName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test" {
  name = "genieservice-%s"
  description = "This is our main service"
  team_id = "${opsgenie_team.test.id}"
}
data "opsgenie_service" "existingservice" {
  name = "${opsgenie_service.test.name}"
  depends_on = [opsgenie_service.test]

}
`, randomTeamName, randomServiceName)
}
