package opsgenie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
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

func flattenTags(d *schema.ResourceData, fieldName string) []string {
	input := d.Get(fieldName).(*schema.Set)
	tags := make([]string, len(input.List()))

	if input == nil {
		return tags
	}

	for k, v := range input.List() {
		tags[k] = v.(string)
	}

	return tags
}

func expandOpsGenieTimeRestriction(d []interface{}) *og.TimeRestriction {
	timeRestriction := og.TimeRestriction{}

	for _, v := range d {
		config := v.(map[string]interface{})
		timeRestriction.Type = og.RestrictionType(config["type"].(string))

		if config["restrictions"].(*schema.Set).Len() > 0 {
			restrictionList := make([]og.Restriction, 0, config["restrictions"].(*schema.Set).Len())
			for _, v := range config["restrictions"].(*schema.Set).List() {
				config := v.(map[string]interface{})
				startHour := uint32(config["start_hour"].(int))
				startMin := uint32(config["start_min"].(int))
				endHour := uint32(config["end_hour"].(int))
				endMin := uint32(config["end_min"].(int))
				restriction := og.Restriction{
					StartDay:  og.Day(config["start_day"].(string)),
					StartHour: &startHour,
					StartMin:  &startMin,
					EndHour:   &endHour,
					EndDay:    og.Day(config["end_day"].(string)),
					EndMin:    &endMin,
				}
				restrictionList = append(restrictionList, restriction)
			}
			timeRestriction.RestrictionList = restrictionList
		} else {
			restriction := og.Restriction{}
			for _, v := range config["restriction"].(*schema.Set).List() {
				config := v.(map[string]interface{})
				startHour := uint32(config["start_hour"].(int))
				startMin := uint32(config["start_min"].(int))
				endHour := uint32(config["end_hour"].(int))
				endMin := uint32(config["end_min"].(int))
				restriction = og.Restriction{
					StartHour: &startHour,
					StartMin:  &startMin,
					EndHour:   &endHour,
					EndMin:    &endMin,
				}
			}

			timeRestriction.Restriction = restriction
		}
	}
	return &timeRestriction
}

func flattenOpsgenieTimeRestriction(input *og.TimeRestriction) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, 1)
	if input == nil || input.Type == "" {
		// If type is not set, time restriction should be empty.
		return output
	}

	element := make(map[string]interface{})

	if len(input.RestrictionList) > 0 {
		restrictions := make([]map[string]interface{}, 0, len(input.RestrictionList))
		for _, r := range input.RestrictionList {
			restrictionMap := make(map[string]interface{})
			restrictionMap["start_min"] = r.StartMin
			restrictionMap["start_hour"] = r.StartHour
			restrictionMap["start_day"] = r.StartDay
			restrictionMap["end_min"] = r.EndMin
			restrictionMap["end_hour"] = r.EndHour
			restrictionMap["end_day"] = r.EndDay
			restrictions = append(restrictions, restrictionMap)
		}
		element["restrictions"] = restrictions
	} else {
		restriction := make([]map[string]interface{}, 0, 1)
		restrictionMap := make(map[string]interface{})
		restrictionMap["start_min"] = input.Restriction.StartMin
		restrictionMap["start_hour"] = input.Restriction.StartHour
		restrictionMap["end_min"] = input.Restriction.EndMin
		restrictionMap["end_hour"] = input.Restriction.EndHour
		restriction = append(restriction, restrictionMap)
		element["restriction"] = restriction
	}

	element["type"] = input.Type
	output = append(output, element)
	return output
}
