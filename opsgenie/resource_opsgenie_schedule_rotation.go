package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpsgenieScheduleRotation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieScheduleRotationCreate,
		Read:   handleNonExistentResource(resourceOpsgenieScheduleRotationRead),
		Update: resourceOpsgenieScheduleRotationUpdate,
		Delete: resourceOpsgenieScheduleRotationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected schedule_id/schedule_rotation_id", d.Id())
				}
				scheduleId := idParts[0]
				scheduleRotationId := idParts[1]
				d.Set("schedule_id", scheduleId)
				d.SetId(scheduleRotationId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"schedule_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_date": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDateWithMinutes,
			},
			"end_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDateWithMinutes,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsgenieScheduleRotationType,
			},
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"participant": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateScheduleRotationParticipantType,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"time_restriction": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateRestrictionType,
						},
						"restrictions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_day": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateDay,
									},
									"end_day": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateDay,
									},
									"start_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"start_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
									"end_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"end_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
								},
							},
						},
						"restriction": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"start_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
									"end_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"end_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceOpsgenieScheduleRotationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	scheduleIdentiferValue := d.Get("schedule_id").(string)

	name := d.Get("name").(string)
	start_date := d.Get("start_date").(string)
	end_date := d.Get("end_date").(string)
	rotationType := d.Get("type").(string)
	length := d.Get("length").(int)
	timeRestriction := d.Get("time_restriction").([]interface{})
	participants := d.Get("participant").([]interface{})
	layoutStr := "2006-01-02T15:04:05Z"
	startDate, err := time.Parse(layoutStr, start_date)
	if err != nil {
		return fmt.Errorf("Cannot parse date-time")
	}

	createRequest := &schedule.CreateRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		Rotation: &og.Rotation{
			StartDate:    &startDate,
			Length:       uint32(length),
			Type:         og.RotationType(rotationType),
			Participants: expandOpsgenieScheduleParticipants(participants),
		},
	}

	if name != "" {
		createRequest.Rotation.Name = name
	}
	if end_date != "" {
		endDate, err := time.Parse(layoutStr, end_date)
		if err != nil {
			log.Fatal(err)
		}
		createRequest.Rotation.EndDate = &endDate
	}
	if length != 0 {
		createRequest.Rotation.Length = uint32(length)
	}
	if len(timeRestriction) > 0 {
		createRequest.Rotation.TimeRestriction = expandTimeRestrictions(timeRestriction)
	}

	log.Printf("[INFO] Creating OpsGenie rotation '%s'", name)

	result, err := client.CreateRotation(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsgenieScheduleRotationRead(d, meta)
}

func resourceOpsgenieScheduleRotationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	scheduleIdentiferValue := d.Get("schedule_id").(string)

	getRequest := &schedule.GetRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		RotationId:              d.Id(),
	}
	getResponse, err := client.GetRotation(context.Background(), getRequest)
	if err != nil {
		return err
	}
	startDate := getResponse.StartDate.Format("2006-01-02T15:04:05Z")
	d.SetId(getResponse.Rotation.Id)
	d.Set("name", getResponse.Rotation.Name)
	d.Set("length", getResponse.Length)
	d.Set("type", getResponse.Type)
	d.Set("participant", flattenOpsgenieScheduleRotationParticipant(getResponse.Participants))
	if getResponse.TimeRestriction != nil {
		d.Set("time_restriction", flattenOpsgenieScheduleRotationTimeRestriction(getResponse.TimeRestriction))
	}
	d.Set("start_date", startDate)
	if getResponse.EndDate != nil {
		endDate := getResponse.EndDate.Format("2006-01-02T15:04:05Z")
		d.Set("end_date", endDate)
	}

	return nil
}

func flattenOpsgenieScheduleRotationParticipant(input []og.Participant) []map[string]interface{} {
	participants := make([]map[string]interface{}, 0, len(input))
	for _, part := range input {
		outputMember := make(map[string]interface{})
		outputMember["id"] = part.Id
		outputMember["type"] = part.Type
		participants = append(participants, outputMember)
	}

	return participants
}

func flattenOpsgenieScheduleRotationRestriction(input og.Restriction) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, 1)
	restriction := make(map[string]interface{})

	restriction["start_hour"] = input.StartHour
	restriction["end_hour"] = input.EndHour
	restriction["start_min"] = input.StartMin
	restriction["end_min"] = input.EndMin
	output = append(output, restriction)

	return output
}

func flattenOpsgenieScheduleRotationRestrictions(input []og.Restriction) []map[string]interface{} {
	restrictions := make([]map[string]interface{}, 0, len(input))
	for _, restriction := range input {
		outputRestriction := make(map[string]interface{})
		outputRestriction["start_hour"] = restriction.StartHour
		outputRestriction["end_hour"] = restriction.EndHour
		outputRestriction["start_min"] = restriction.StartMin
		outputRestriction["end_min"] = restriction.EndMin
		outputRestriction["start_day"] = restriction.StartDay
		outputRestriction["end_day"] = restriction.EndDay
		restrictions = append(restrictions, outputRestriction)
	}

	return restrictions
}

func flattenOpsgenieScheduleRotationTimeRestriction(input *og.TimeRestriction) []map[string]interface{} {
	output := make([]map[string]interface{}, 0, 1)
	timeRestriction := make(map[string]interface{})
	timeRestriction["type"] = input.Type
	if timeRestriction["type"] == og.TimeOfDay {
		timeRestriction["restriction"] = flattenOpsgenieScheduleRotationRestriction(input.Restriction)
	} else if timeRestriction["type"] == og.WeekdayAndTimeOfDay {
		timeRestriction["restrictions"] = flattenOpsgenieScheduleRotationRestrictions(input.RestrictionList)
	}
	output = append(output, timeRestriction)

	return output
}

func resourceOpsgenieScheduleRotationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	scheduleIdentiferValue := d.Get("schedule_id").(string)

	name := d.Get("name").(string)
	start_date := d.Get("start_date").(string)
	end_date := d.Get("end_date").(string)
	rotationType := d.Get("type").(string)
	length := d.Get("length").(int)
	timeRestriction := d.Get("time_restriction").([]interface{})
	participants := d.Get("participant").([]interface{})
	layoutStr := "2006-01-02T15:04:05Z"
	startDate, err := time.Parse(layoutStr, start_date)
	if err != nil {
		log.Fatal(err)
	}

	updateRequest := &schedule.UpdateRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		RotationId:              d.Id(),
		Rotation: &og.Rotation{
			StartDate:    &startDate,
			Length:       uint32(length),
			Type:         og.RotationType(rotationType),
			Participants: expandOpsgenieScheduleParticipants(participants),
		},
	}

	if name != "" {
		updateRequest.Rotation.Name = name
	}
	if end_date != "" {
		endDate, err := time.Parse(layoutStr, end_date)
		if err != nil {
			log.Fatal(err)
		}
		updateRequest.Rotation.EndDate = &endDate
	}
	if len(timeRestriction) > 0 {
		updateRequest.Rotation.TimeRestriction = expandTimeRestrictions(timeRestriction)
	}
	log.Printf("[INFO] Updating OpsGenie schedule rotation '%s'", name)

	_, err = client.UpdateRotation(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieScheduleRotationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie schedule rotation '%s'", d.Get("name").(string))
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	scheduleIdentiferValue := d.Get("schedule_id").(string)

	deleteRequest := &schedule.DeleteRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		RotationId:              d.Id(),
	}

	_, err = client.DeleteRotation(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func expandOpsgenieScheduleParticipants(input []interface{}) []og.Participant {
	participants := make([]og.Participant, 0, len(input))

	if input == nil {
		return participants
	}

	for _, v := range input {
		config := v.(map[string]interface{})

		participantType := config["type"].(string)
		Id := config["id"].(string)

		participant := og.Participant{
			Type: og.ParticipantType(participantType),
			Id:   Id,
		}

		participants = append(participants, participant)
	}

	return participants
}
func expandTimeRestrictions(d []interface{}) *og.TimeRestriction {
	timeRestriction := og.TimeRestriction{}

	for _, v := range d {
		config := v.(map[string]interface{})

		timeRestrictionType := config["type"].(string)
		timeRestriction.Type = og.RestrictionType(timeRestrictionType)

		if len(config["restrictions"].([]interface{})) > 0 {
			timeRestriction.RestrictionList = expandOpsgenieRestrictions(config["restrictions"].([]interface{}))
		} else {
			timeRestriction.Restriction = expandOpsgenieRestriction(config["restriction"].([]interface{}))
		}
	}

	return &timeRestriction
}

func expandOpsgenieRestrictions(input []interface{}) []og.Restriction {
	restrictionList := make([]og.Restriction, 0, len(input))

	if input == nil {
		return restrictionList
	}

	for _, v := range input {
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

	return restrictionList
}

func expandOpsgenieRestriction(input []interface{}) og.Restriction {

	restriction := og.Restriction{}
	for _, v := range input {
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

	return restriction
}

func validateOpsgenieScheduleRotationType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"daily":  true,
		"weekly": true,
		"hourly": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Opsgenie Schedule Rotation Type  can only be 'Daily' ,'Weekly' or 'Hourly'"))
	}
	return
}

func validateScheduleRotationParticipantType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"user":       true,
		"team":       true,
		"escalation": true,
		"none":       true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("it can only be one of these 'user', 'schedule', 'team', 'escalation'"))
	}
	return
}

func validateDay(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"monday":    true,
		"tuesday":   true,
		"wednesday": true,
		"thursday":  true,
		"friday":    true,
		"saturday":  true,
		"sunday":    true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("it can only be day of week (monday,tuesday...)"))
	}
	return
}

func validateHourParams(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value < 0 || value > 24 {
		errors = append(errors, fmt.Errorf("hour must between 0-24"))
	}
	return
}
func validateMinParams(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value < 0 || value > 59 {
		errors = append(errors, fmt.Errorf("minute must in between of 0-59"))

	}
	return
}

func validateRestrictionType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	families := map[string]bool{
		"time-of-day":             true,
		"weekday-and-time-of-day": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("time restriction type must be time-of-day or weekday-and-time-of-day"))
	}
	return
}
