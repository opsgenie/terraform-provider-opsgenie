package opsgenie

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func timeRestrictionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"time-of-day", "weekday-and-time-of-day"}, false),
				},
				"restrictions": {
					Type:          schema.TypeSet,
					Optional:      true,
					ConflictsWith: []string{"time_restriction.0.restriction"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"start_day": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}, false),
							},
							"end_day": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}, false),
							},
							"start_hour": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 23),
							},
							"start_min": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 59),
							},
							"end_hour": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 23),
							},
							"end_min": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 59),
							},
						},
					},
				},
				"restriction": {
					Type:          schema.TypeSet,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"time_restriction.0.restrictions"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"start_hour": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 23),
							},
							"start_min": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 59),
							},
							"end_hour": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 23),
							},
							"end_min": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 59),
							},
						},
					},
				},
			},
		},
	}
}
