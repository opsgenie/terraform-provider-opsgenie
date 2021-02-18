package opsgenie

import (
	"context"
	"fmt"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

func init() {
	resource.AddTestSweepers("opsgenie_heartbeat", &resource.Sweeper{
		Name: "opsgenie_heartbeat",
		F:    testSweepHeartbeat,
	})

}

func testSweepHeartbeat(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := heartbeat.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background())
	if err != nil {
		return err
	}

	for _, u := range resp.Heartbeats {
		if strings.HasPrefix(u.Name, "genieheartbeat-") {
			log.Printf("Destroying heartbeat %s", u.Name)

			if _, err := client.Delete(context.Background(), u.Name); err != nil {
				return err
			}

		}
	}

	return nil
}

func TestAccOpsgenieHeartbeat_basic(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomHeartbeat := acctest.RandString(6)
	config := testAccOpsGenieHeartbeat_basic(randomTeam, randomHeartbeat)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieHeartbeatDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieHeartbeatExists("opsgenie_heartbeat.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieHeartbeatDestroy(s *terraform.State) error {
	client, err := heartbeat.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_api_integration" {
			continue
		}
		name := rs.Primary.Attributes["id"]

		_, err := client.Get(context.Background(), name)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Heartbeat still exists: %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieHeartbeatExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := heartbeat.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		name := rs.Primary.Attributes["id"]

		_, err = client.Get(context.Background(), name)
		if err != nil {
			return fmt.Errorf("Bad: Heartbeat with name %q does not exist", name)
		}
		return nil
	}
}

func testAccOpsGenieHeartbeat_basic(randomTeam, randomHeartbeat string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteamw-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_heartbeat" "test" {
	name = "genieheartbeat-%s"
	description = "test opsgenie heartbeat terraform"
	interval_unit = "minutes"
	interval = 10
	enabled = false
	alert_message = "Test"
	alert_priority = "P3"
	alert_tags = ["test","fahri"]
	owner_team_id = "${opsgenie_team.test.id}"

}
`, randomTeam, randomHeartbeat)

}
