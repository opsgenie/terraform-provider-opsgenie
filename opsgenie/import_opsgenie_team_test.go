package opsgenie

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOpsGenieTeam_importBasic(t *testing.T) {
	resourceName := "opsgenie_team.test"

	ri := acctest.RandInt()
	config := testAccOpsGenieTeam_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOpsGenieTeam_importWithEmptyDescription(t *testing.T) {
	resourceName := "opsgenie_team.test"

	ri := acctest.RandInt()
	config := testAccOpsGenieTeam_withEmptyDescription(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOpsGenieTeam_importWithUser(t *testing.T) {
	resourceName := "opsgenie_team.test"

	ri := acctest.RandInt()
	config := testAccOpsGenieTeam_withUser(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOpsGenieTeam_importWithUserComplete(t *testing.T) {
	resourceName := "opsgenie_team.test"

	ri := acctest.RandInt()
	config := testAccOpsGenieTeam_withUserComplete(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
