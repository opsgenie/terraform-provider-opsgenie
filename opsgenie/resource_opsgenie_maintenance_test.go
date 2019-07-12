package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/maintenance"
	"log"
	"strings"
	"testing"
	"time"
)

func init() {
	resource.AddTestSweepers("opsgenie_maintenance", &resource.Sweeper{
		Name: "opsgenie_maintenance",
		F:    testSweepMaintenance,
	})

}

func testSweepMaintenance(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := maintenance.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &maintenance.ListRequest{})
	if err != nil {
		return err
	}

	for _, u := range resp.Maintenances {
		if strings.HasPrefix(u.Description, "geniemaintenance-") {
			log.Printf("Destroying maintenance %s", u.Description)

			deleteRequest := maintenance.DeleteRequest{
				Id: u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieMaintenance_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieMaintenance_complete(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieMaintenanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieMaintenanceExists("opsgenie_maintenance.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieMaintenanceDestroy(s *terraform.State) error {
	client, err := maintenance.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_maintenance" {
			continue
		}
		Id := rs.Primary.Attributes["id"]

		req := maintenance.GetRequest{
			Id: Id,
		}
		_,err:=client.Get(context.Background(),&req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Maintenance still exists : %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieMaintenanceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := maintenance.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		Id := rs.Primary.Attributes["id"]

		req := maintenance.GetRequest{
			Id: Id,
		}

		_, err = client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Maintenance with id %q does not exist", Id)
		}
		return nil
	}
}

func testAccOpsGenieMaintenance_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_maintenance" "test" {
  description = "geniemaintenance-%s"
  time {
    type = "schedule"
    start_date = "2019-06-20T17:45:00Z"
    end_date  = "2019-06-20T17:50:00Z"
  }
rules{
}
}
`, rString)
}

func testAccOpsGenieMaintenance_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_email_integration" "test" {
  name = "testemailapi-maintenance"
  email_username ="fahri"
}
resource "opsgenie_maintenance" "test" {
  description = "geniemaintenance-%s"
  time {
    type = "schedule"
    start_date = "2019-06-20T17:45:00Z"
    end_date  = "2019-06-%dT17:50:00Z"
  }
  rules {
    state = "enabled"
    entity {
      id = "${opsgenie_email_integration.test.id}"
      type = "integration"
    }
  }
}
`, rString, time.Now().Day())
}
