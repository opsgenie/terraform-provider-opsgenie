package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("opsgenie_saved_search", &resource.Sweeper{
		Name: "opsgenie_saved_search",
		F:    testSweepSavedSearch,
	})

}

func testSweepSavedSearch(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}
	fmt.Println("Starting testSweepSavedSearch")

	client, err := alert.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	resp, err := client.ListSavedSearches(context.Background(), &alert.ListSavedSearchRequest{})
	if err != nil {
		return err
	}

	if strings.HasPrefix(resp.Name, "genietest-") {
		deleteRequest := &alert.DeleteSavedSearchRequest{
			IdentifierType:  "name",
			IdentifierValue: resp.Name,
		}
		fmt.Printf("Deleting saved search %s", resp.Name)
		if _, err := client.DeleteSavedSearch(context.Background(), deleteRequest); err != nil {
			return err
		}
	}

	return nil
}

func TestAccOpsGenieSavedSearch_basic(t *testing.T) {
	randomUser := acctest.RandString(6)
	config := testAccOpsGenieSavedSearch_basic(randomUser)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieSavedSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieSavedSearchExists("opsgenie_saved_search.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieSavedSearch_complete(t *testing.T) {
	randomUser := acctest.RandString(6)
	randomTeam := acctest.RandString(6)
	config := testAccOpsGenieSavedSearch_complete(randomUser, randomTeam)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieSavedSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieSavedSearchExists("opsgenie_saved_search.test"),
				),
			},
		},
	})
}
func testCheckOpsGenieSavedSearchDestroy(s *terraform.State) error {
	client, err := alert.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_saved_search" {
			continue
		}
		Id := rs.Primary.Attributes["id"]

		req := alert.GetSavedSearchRequest{
			IdentifierValue: Id,
		}
		_, err := client.GetSavedSearch(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Saved search still exists : %s", x.Error()))
			}
		}
	}
	return nil
}

func testCheckOpsGenieSavedSearchExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		Id := rs.Primary.Attributes["id"]

		client, err := alert.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := alert.GetSavedSearchRequest{
			IdentifierValue: Id,
		}

		_, err = client.GetSavedSearch(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: SavedSearch with id %q does not exist", Id)
		}
		return nil
	}
}

func testAccOpsGenieSavedSearch_basic(randomUser string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}


resource "opsgenie_saved_search" "test" { 
    name = "test-search"
    owner ="${opsgenie_user.test.id}"
    query = "priority: P3"
`, randomUser)
}

func testAccOpsGenieSavedSearch_complete(randomUser, randomTeam string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}

resource "opsgenie_team" "test" {
	name        = "genieteam-%s"
	description = "This team deals with all the things"
  }


resource "opsgenie_saved_search" "test" { 
    name = "test-search"
    owner ="${opsgenie_user.test.id}"
    query = "priority: P3"
	description = "Testing saved search"
	teams {
        id = opsgenie_team.test.id
    }

`, randomUser, randomTeam)
}
