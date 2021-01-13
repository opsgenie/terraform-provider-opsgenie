package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
	"log"
)

func dataSourceOpsGenieUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieUserRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieUserUsername,
			},
			"full_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOpsGenieUserFullName,
			},
			"role": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOpsGenieUserRole,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en_US",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "America/New_York",
			},
		},
	}
}

func dataSourceOpsGenieUserRead(d *schema.ResourceData, meta interface{}) error {
	client, err := user.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	username := d.Get("username").(string)

	log.Printf("[INFO] Reading OpsGenie user '%s'", username)

	usr, err := client.Get(context.Background(), &user.GetRequest{
		Identifier: username,
	})
	if err != nil {
		return err
	}

	d.SetId(usr.Id)
	d.Set("username", usr.Username)
	d.Set("full_name", usr.FullName)
	d.Set("role", usr.Role.RoleName)
	d.Set("locale", usr.Locale)
	d.Set("timezone", usr.TimeZone)

	return nil
}
