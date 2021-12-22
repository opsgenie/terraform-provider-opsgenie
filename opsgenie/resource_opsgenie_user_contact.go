package opsgenie

import (
	"context"
	"fmt"
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk-v2/contact"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpsGenieUserContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieUserContactCreate,
		Read:   handleNonExistentResource(resourceOpsGenieUserContactRead),
		Update: resourceOpsGenieUserContactUpdate,
		Delete: resourceOpsGenieUserContactDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected username/contact_id", d.Id())
				}
				username := idParts[0]
				contactId := idParts[1]
				d.Set("username", username)
				d.SetId(contactId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"to": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieContactTo,
			},
			"method": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieContactMethod,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceOpsGenieUserContactCreate(d *schema.ResourceData, meta interface{}) error {

	client, err := contact.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	userId := d.Get("username").(string)
	method := d.Get("method").(string)
	enabled := d.Get("enabled").(bool)
	to := d.Get("to").(string)

	createRequest := &contact.CreateRequest{
		UserIdentifier:  userId,
		To:              to,
		MethodOfContact: contact.MethodType(method),
	}

	result, err := client.Create(context.Background(), createRequest)
	if err != nil {
		return err
	}
	d.SetId(result.Id)

	if enabled {
		_, err = client.Enable(context.Background(), &contact.EnableRequest{
			UserIdentifier:    userId,
			ContactIdentifier: result.Id,
		})
		if err != nil {
			return err
		}
	} else {
		_, err = client.Disable(context.Background(), &contact.DisableRequest{
			UserIdentifier:    userId,
			ContactIdentifier: result.Id,
		})
		if err != nil {
			return err
		}
	}

	return resourceOpsGenieUserContactRead(d, meta)
}

func resourceOpsGenieUserContactRead(d *schema.ResourceData, meta interface{}) error {
	client, err := contact.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	userId := d.Get("username").(string)

	contactsResult, err := client.Get(context.Background(), &contact.GetRequest{
		UserIdentifier:    userId,
		ContactIdentifier: d.Id(),
	})
	if err != nil {
		return err
	}

	d.Set("method", contactsResult.MethodOfContact)
	d.Set("to", contactsResult.To)
	d.Set("enabled", contactsResult.Status.Enabled)

	return nil
}

func resourceOpsGenieUserContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := contact.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	userId := d.Get("username").(string)
	enabled := d.Get("enabled").(bool)
	to := d.Get("to").(string)

	updateRequest := &contact.UpdateRequest{
		UserIdentifier:    userId,
		ContactIdentifier: d.Id(),
		To:                to,
	}
	_, err = client.Update(context.Background(), updateRequest)
	if err != nil {
		return err
	}
	if enabled {
		_, err = client.Enable(context.Background(), &contact.EnableRequest{
			UserIdentifier:    userId,
			ContactIdentifier: d.Id(),
		})
		if err != nil {
			return err
		}
	} else {
		_, err = client.Disable(context.Background(), &contact.DisableRequest{
			UserIdentifier:    userId,
			ContactIdentifier: d.Id(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceOpsGenieUserContactDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := contact.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	userId := d.Get("username").(string)

	deleteRequest := &contact.DeleteRequest{
		UserIdentifier:    userId,
		ContactIdentifier: d.Id(),
	}

	dr, err := client.Delete(context.Background(), deleteRequest)
	if err != nil {
		return err
	}
	_ = dr

	return nil
}

func validateOpsGenieContactMethod(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"sms":   true,
		"voice": true,
		"email": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("you write wrong opsgenie contact method. valid methods are: sms,voice and email"))
	}
	return
}

func validateOpsGenieContactTo(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 512 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 512 characters: %q %d", k, value, len(value)))
	}
	return
}
