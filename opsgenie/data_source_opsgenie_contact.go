package opsgenie

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/contact"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

func dataSourceOpsGenieContact() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieContactRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieUserUsername,
			},
			"contact_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceOpsGenieContactGetID(username string, meta interface{}) (string, error) {
	client, err := user.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return "", err
	}

	log.Printf("[INFO] Reading OpsGenie user '%s'", username)

	usr, err := client.Get(context.Background(), &user.GetRequest{
		Identifier: username,
	})
	if err != nil {
		return "", err
	}
	return usr.Id, err
}

func dataSourceOpsGenieContactRead(d *schema.ResourceData, meta interface{}) error {
	client, err := contact.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	usr := d.Get("username").(string)
	log.Printf("[INFO] Reading OpsGenie user contact '%s'", usr)

	getRequest := &contact.ListRequest{
		UserIdentifier: usr,
	}

	getResponse, err := client.List(context.Background(), getRequest)
	if err != nil {
		return err
	}

	id, err := dataSourceOpsGenieContactGetID(usr, meta)
	if err != nil {
		return err
	}
	id = id + "_contact"

	log.Printf("[TEST]")
	log.Print(flattenOpsgenieContactList(getResponse.Contact))
	d.SetId(id)
	d.Set("contact_list", flattenOpsgenieContactList(getResponse.Contact))
	// d.Set("contact_list", getResponse.Contact)

	return nil
}

func flattenOpsgenieContactList(input []contact.Contact) []map[string]interface{} {

	res := make([]map[string]interface{}, 0, len(input))
	for _, r := range input {
		out := make(map[string]interface{})
		out["id"] = r.Id
		out["method"] = r.MethodOfContact
		out["to"] = r.To
		res = append(res, out)
	}
	return res
}
