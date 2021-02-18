package opsgenie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

// handleNonExistentResource is a wrapper of resourceFunc that
// handles errors returned by a read function.
func handleNonExistentResource(f schema.ReadFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		if err := f(d, meta); err != nil {
			// if the error that we receive is an ApiError and
			// the status code is 404, it means we need to re-create
			// the specific resource
			apiErr, ok := err.(*client.ApiError)
			if !ok || apiErr.StatusCode != http.StatusNotFound {
				return err
			}
			d.SetId("")
			return nil
		}

		return nil
	}
}

func validateDateWithMinutes(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	layoutStr := "2006-01-02T15:04:05Z"
	selectedTime, err := time.Parse(layoutStr, value)
	if err != nil {
		errors = append(errors, fmt.Errorf("wrong date-time format:%s, Its should be like this : %s", value, layoutStr))
		return
	}
	min := selectedTime.Minute()
	if min%30 != 0 {
		errors = append(errors, fmt.Errorf("you can only select 30 or 00 for minutes. you enter : %d", selectedTime.Minute()))
	}

	return
}

func validateDate(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	layoutStr := "2006-01-02T15:04:05Z"
	_, err := time.Parse(layoutStr, value)
	if err != nil {
		errors = append(errors, fmt.Errorf("wrong date-time format:%s, Its should be like this : %s", value, layoutStr))
		return
	}

	return
}

func convertStringMapToInterfaceMap(old map[string]string) map[string]interface{} {
	new := map[string]interface{}{}
	for k, v := range old {
		new[k] = v
	}
	return new
}
func convertStringSliceToInterfaceSlice(old []string) []interface{} {
	new := make([]interface{}, len(old))
	for k, v := range old {
		new[k] = v
	}
	return new
}
