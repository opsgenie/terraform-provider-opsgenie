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
	"github.com/opsgenie/opsgenie-go-sdk-v2/service"
)

func init() {
	resource.AddTestSweepers("opsgenie_service", &resource.Sweeper{
		Name: "opsgenie_service",
		F:    testSweepService,
	})

}

func testSweepService(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := service.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &service.ListRequest{})
	if err != nil {
		return err
	}

	for _, svc := range resp.Services {
		if strings.HasPrefix(svc.Name, "genietest-") {
			log.Printf("Destroying service %s", svc.Name)

			deleteRequest := service.DeleteRequest{
				Id: svc.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieService_basic(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomService := acctest.RandString(6)

	config := testAccOpsGenieService_basic(randomTeam, randomService)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieServiceExists("opsgenie_service.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieService_complete(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomService := acctest.RandString(6)
	randomDescription := acctest.RandString(20)

	config := testAccOpsGenieService_complete(randomTeam, randomService, randomDescription)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieServiceExists("opsgenie_service.test"),
					resource.TestCheckResourceAttr("opsgenie_service.test", "description", randomDescription),
				),
			},
		},
	})
}

func testCheckOpsGenieServiceDestroy(s *terraform.State) error {
	client, err := service.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_service" {
			continue
		}
		req := service.GetRequest{
			Id: rs.Primary.Attributes["id"],
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Service still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]
		serviceName := rs.Primary.Attributes["name"]

		client, err := service.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := service.GetRequest{
			Id: rs.Primary.Attributes["id"],
		}

		result, err := client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Service %q (service: %q) does not exist", id, serviceName)
		} else {
			log.Printf("Service found :%s ", result.Service.Name)
		}

		return nil
	}
}

func testAccOpsGenieService_basic(randomTeam, randomService string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test" {
  name  = "genietest-%s"
  team_id = "${opsgenie_team.test.id}"
}
`, randomTeam, randomService)
}

func testAccOpsGenieService_complete(randomTeam, randomService, randomDescription string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test" {
  name  = "genietest-%s"
  team_id = "${opsgenie_team.test.id}"
  description = "%s"
}
`, randomTeam, randomService, randomDescription)
}
