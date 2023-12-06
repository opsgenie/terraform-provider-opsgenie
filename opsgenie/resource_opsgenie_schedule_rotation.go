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
			"time_restriction": timeRestrictionSchema(),
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
		createRequest.Rotation.TimeRestriction = expandOpsGenieTimeRestriction(timeRestriction)
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
		d.Set("time_restriction", flattenOpsgenieTimeRestriction(getResponse.TimeRestriction))
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
		updateRequest.Rotation.TimeRestriction = expandOpsGenieTimeRestriction(timeRestriction)
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
