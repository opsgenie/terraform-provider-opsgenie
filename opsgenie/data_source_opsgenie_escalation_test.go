package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceOpsGenieEscalation_Basic(t *testing.T) {
	randomUserName := acctest.RandString(6)
	randomEscalation := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieEscalationConfig(randomUserName, randomEscalation),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieEscalation("opsgenie_escalation.test", "data.opsgenie_escalation.existingescalation"),
				),
			},
		},
	})
}

func testAccDataSourceOpsGenieEscalation(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["name"] == "" {
			return fmt.Errorf("expected to get an escalation name from OpsGenie")
		}

		testAtts := []string{"name", "rules", "description", "repeat"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("expected the escalation %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieEscalationConfig(randomUser, randomEscalation string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
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
data "opsgenie_escalation" "existingescalation" {
  name = "${opsgenie_escalation.test.name}"
  depends_on = [opsgenie_escalation.test]
}
`, randomUser, randomEscalation)
}
