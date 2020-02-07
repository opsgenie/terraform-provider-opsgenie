package alert

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateRequest_Validate(t *testing.T) {
	createRequestWithoutMessage := &CreateAlertRequest{
		Alias: "alias",
	}
	err := createRequestWithoutMessage.Validate()

	assert.Equal(t, err.Error(), errors.New("message can not be empty").Error())

	createRequest := &CreateAlertRequest{
		Message: "message",
	}

	err = createRequest.Validate()

	assert.Equal(t, err, nil)

}

func TestAcknowledgeAlertRequest_Validate(t *testing.T) {
	acknowledgeAlertRequestWithError := &AcknowledgeAlertRequest{}
	err := acknowledgeAlertRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	acknowledgeAlertRequest := &AcknowledgeAlertRequest{
		IdentifierType:  ALIAS,
		IdentifierValue: "alias1",
	}
	err = acknowledgeAlertRequest.Validate()

	assert.Equal(t, err, nil)

}

func TestAddNoteRequest_Validate(t *testing.T) {
	addNoteRequest := &AddNoteRequest{
		IdentifierType:  ALIAS,
		IdentifierValue: "alias2",
		Note:            "",
	}
	err := addNoteRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Note can not be empty").Error())

	addNoteRequestWithoutidentifier := &AddNoteRequest{
		IdentifierType: TINYID,
		Note:           "note test",
	}
	err = addNoteRequestWithoutidentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	noteRequest := &AddNoteRequest{
		IdentifierType:  ALIAS,
		IdentifierValue: "alias1",
		Note:            "note1",
	}
	err = noteRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestCloseAlertRequest_Validate(t *testing.T) {
	closeAlertRequestWithError := &CloseAlertRequest{}
	err := closeAlertRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	closeAlertRequest := &CloseAlertRequest{
		IdentifierType:  ALIAS,
		IdentifierValue: "alias1",
	}
	err = closeAlertRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestCreateSavedSearchRequest_Validate(t *testing.T) {
	createSavedSearchRequestWithoutName := &CreateSavedSearchRequest{
		Name:  "",
		Query: "status: open",
		Owner: User{
			ID:       "id",
			Username: "user",
		},
		Description: "Test",
		Teams:       nil,
	}
	err := createSavedSearchRequestWithoutName.Validate()

	assert.Equal(t, err.Error(), errors.New("Name can not be empty").Error())

	createSavedSearchRequestWithoutQuery := &CreateSavedSearchRequest{
		Name:  "name1",
		Query: "",
		Owner: User{
			ID:       "id",
			Username: "user",
		},
		Description: "Test",
		Teams:       nil,
	}
	err = createSavedSearchRequestWithoutQuery.Validate()

	assert.Equal(t, err.Error(), errors.New("Query can not be empty").Error())

	createSavedSearchRequestWithoutOwner := &CreateSavedSearchRequest{
		Name:        "name1",
		Query:       "status: open",
		Owner:       User{},
		Description: "Test",
		Teams:       nil,
	}
	err = createSavedSearchRequestWithoutOwner.Validate()

	assert.Equal(t, err.Error(), errors.New("Owner can not be empty").Error())

	savedSearchRequest := &CreateSavedSearchRequest{
		Name:  "name1",
		Query: "status: open",
		Owner: User{
			ID:       "id",
			Username: "user",
		},
		Description: "Test",
	}
	err = savedSearchRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestDeleteAlertRequest_Validate(t *testing.T) {
	deleteAlertRequestWithError := &DeleteAlertRequest{}
	err := deleteAlertRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	deleteAlertRequest := &DeleteAlertRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
	}
	err = deleteAlertRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestDeleteSavedSearchRequest_Validate(t *testing.T) {
	deleteSavedSearchRequestWithError := &DeleteSavedSearchRequest{}
	err := deleteSavedSearchRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	deleteSavedSearchRequest := &DeleteSavedSearchRequest{
		IdentifierType:  NAME,
		IdentifierValue: "name1",
	}
	err = deleteSavedSearchRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestEscalateToNextRequest_Validate(t *testing.T) {
	escalateToNextRequestWithoutEscalation := &EscalateToNextRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Escalation:      Escalation{},
	}
	err := escalateToNextRequestWithoutEscalation.Validate()

	assert.Equal(t, err.Error(), errors.New("Escalation ID or name must be defined").Error())

	escalateToNextRequestWithoutIdentifier := &EscalateToNextRequest{
		IdentifierType: ALERTID,
		Escalation: Escalation{
			ID:   "",
			Name: "escName",
		},
	}
	err = escalateToNextRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	escalateToNextRequest := &EscalateToNextRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "tiny1",
		Escalation: Escalation{
			ID:   "",
			Name: "escName",
		},
	}
	err = escalateToNextRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestExecuteCustomActionAlertRequest_Validate(t *testing.T) {
	executeCustomActionAlertRequestWithoutIdentifier := &ExecuteCustomActionAlertRequest{
		IdentifierType: ALIAS,
		Action:         "start",
	}
	err := executeCustomActionAlertRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	executeCustomActionAlertRequestWithoutAction := &ExecuteCustomActionAlertRequest{
		IdentifierType:  ALIAS,
		IdentifierValue: "Test",
	}
	err = executeCustomActionAlertRequestWithoutAction.Validate()

	assert.Equal(t, err.Error(), errors.New("Action can not be empty").Error())

	customActionAlertRequest := &ExecuteCustomActionAlertRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "tiny1",
		Action:          "start",
	}
	err = customActionAlertRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestGetAlertRequest_Validate(t *testing.T) {
	getAlertRequestWithError := &GetAlertRequest{}
	err := getAlertRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	getAlertRequest := &GetAlertRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "tiny1",
	}
	err = getAlertRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestGetAsyncRequestStatusRequest_Validate(t *testing.T) {
	getAsyncRequestStatusRequestWithError := &GetRequestStatusRequest{}
	err := getAsyncRequestStatusRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("RequestId can not be empty").Error())

	asyncRequestStatusRequest := &GetRequestStatusRequest{
		RequestId: "reqId",
	}
	err = asyncRequestStatusRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestSnoozeAlertRequest_Validate(t *testing.T) {
	snoozeAlertRequestWithoutIdentifier := &SnoozeAlertRequest{}
	err := snoozeAlertRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	snoozeAlertRequestWithInvalidEndtime := &SnoozeAlertRequest{
		IdentifierValue: "alias1",
		IdentifierType:  ALIAS,
	}
	err = snoozeAlertRequestWithInvalidEndtime.Validate()

	assert.Equal(t, err.Error(), errors.New("EndTime should at least be 2 seconds later.").Error())

	snoozeAlertRequestWithInvalidEndtime2 := &SnoozeAlertRequest{
		IdentifierValue: "alias1",
		IdentifierType:  ALIAS,
		EndTime:         time.Now(),
	}
	err = snoozeAlertRequestWithInvalidEndtime2.Validate()

	assert.Equal(t, err.Error(), errors.New("EndTime should at least be 2 seconds later.").Error())

	var timeNow = time.Now()
	snoozeAlertRequest := &SnoozeAlertRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "tiny1",
		EndTime:         time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour()+1, timeNow.Minute(), timeNow.Second(), 1, timeNow.Location()),
	}
	err = snoozeAlertRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestUnacknowledgeAlertRequest_Validate(t *testing.T) {
	unacknowledgeAlertRequestWithError := &UnacknowledgeAlertRequest{}
	err := unacknowledgeAlertRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	unacknowledgeAlertRequest := &UnacknowledgeAlertRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "tiny1",
	}
	err = unacknowledgeAlertRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestUpdateSavedSearchRequest_Validate(t *testing.T) {
	updateSavedSearchRequestWithoutIdentifier := &UpdateSavedSearchRequest{}
	err := updateSavedSearchRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	updateSavedSearchRequestWithoutName := &UpdateSavedSearchRequest{
		IdentifierValue: "name1",
		IdentifierType:  NAME,
		NewName:         "",
		Query:           "status: open",
		Owner: User{
			ID:       "id",
			Username: "user",
		},
		Description: "Test",
		Teams:       nil,
	}
	err = updateSavedSearchRequestWithoutName.Validate()

	assert.Equal(t, err.Error(), errors.New("Name can not be empty").Error())

	updateSavedSearchRequestWithoutQuery := &UpdateSavedSearchRequest{
		IdentifierValue: "name1",
		IdentifierType:  NAME,
		NewName:         "name1",
		Query:           "",
		Owner: User{
			ID:       "id",
			Username: "user",
		},
		Description: "Test",
		Teams:       nil,
	}
	err = updateSavedSearchRequestWithoutQuery.Validate()

	assert.Equal(t, err.Error(), errors.New("Query can not be empty").Error())

	updateSavedSearchRequestWithoutOwner := &UpdateSavedSearchRequest{
		IdentifierValue: "name1",
		IdentifierType:  NAME,
		NewName:         "name1",
		Query:           "status: open",
		Owner:           User{},
		Description:     "Test",
		Teams:           nil,
	}
	err = updateSavedSearchRequestWithoutOwner.Validate()

	assert.Equal(t, err.Error(), errors.New("Owner can not be empty").Error())

	updateSavedSearchRequest := &UpdateSavedSearchRequest{
		IdentifierValue: "name1",
		IdentifierType:  NAME,
		NewName:         "name1",
		Query:           "status: open",
		Owner: User{
			ID:       "id",
			Username: "user",
		},
	}
	err = updateSavedSearchRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestAssignRequest_Validate(t *testing.T) {
	assignRequestWithoutOwner := &AssignRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Owner:           User{},
	}
	err := assignRequestWithoutOwner.Validate()

	assert.Equal(t, err.Error(), errors.New("Owner ID or username must be defined").Error())

	assignRequestWithoutIdentifier := &AssignRequest{
		IdentifierType: ALERTID,
		Owner: User{
			ID: "user1",
		},
	}
	err = assignRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	assignRequest := &AssignRequest{
		IdentifierValue: "tiny1",
		IdentifierType:  TINYID,
		Owner: User{
			ID: "user1",
		},
	}
	err = assignRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestAddTeamRequest_Validate(t *testing.T) {
	addTeamRequestWithoutOwner := &AddTeamRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Team:            Team{},
	}
	err := addTeamRequestWithoutOwner.Validate()

	assert.Equal(t, err.Error(), errors.New("Team ID or name must be defined").Error())

	addTeamRequestWithoutIdentifier := &AddTeamRequest{
		IdentifierType: ALERTID,
		Team: Team{
			ID: "team1",
		},
	}
	err = addTeamRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	addTeamRequest := &AddTeamRequest{
		IdentifierValue: "tiny1",
		IdentifierType:  TINYID,
		Team: Team{
			ID: "team1",
		},
	}
	err = addTeamRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestAddTagsRequest_Validate(t *testing.T) {
	addTagsRequestWithoutTags := &AddTagsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Tags:            []string{},
	}
	err := addTagsRequestWithoutTags.Validate()

	assert.Equal(t, err.Error(), errors.New("Tags list can not be empty").Error())

	addTagsRequestWithoutIdentifier := &AddTagsRequest{
		IdentifierType: ALERTID,
		Tags: []string{
			"tags1",
		},
	}
	err = addTagsRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	addTagsRequest := &AddTagsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
		Tags: []string{
			"tags1",
		},
	}
	err = addTagsRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestRemoveTagsRequest_Validate(t *testing.T) {
	removeTagsRequestWithoutTags := &RemoveTagsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Tags:            "",
	}
	err := removeTagsRequestWithoutTags.Validate()

	assert.Equal(t, err.Error(), errors.New("Tags can not be empty").Error())

	removeTagsRequestWithoutIdentifier := &RemoveTagsRequest{
		IdentifierType: ALERTID,
		Tags:           "tags1",
	}
	err = removeTagsRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	removeTagsRequest := &RemoveTagsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
		Tags:            "tags1,tag2",
	}
	err = removeTagsRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestAddDetailsRequest_Validate(t *testing.T) {
	addDetailsRequestWithoutKey := &AddDetailsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
	}
	err := addDetailsRequestWithoutKey.Validate()

	assert.Equal(t, err.Error(), errors.New("Details can not be empty").Error())

	addDetailsRequestWithoutIdentifier := &AddDetailsRequest{
		IdentifierType: ALERTID,
		Details: map[string]string{
			"key":  "value1",
			"key2": "value2",
		},
	}
	err = addDetailsRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	addDetailsRequest := &AddDetailsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
		Details: map[string]string{
			"key":  "value1",
			"key2": "value2",
		},
	}
	err = addDetailsRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestRemoveDetailsRequest_Validate(t *testing.T) {
	removeDetailsRequestWithoutKey := &RemoveDetailsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Keys:            "",
	}
	err := removeDetailsRequestWithoutKey.Validate()

	assert.Equal(t, err.Error(), errors.New("Keys can not be empty").Error())

	removeDetailsRequestWithoutIdentifier := &RemoveDetailsRequest{
		IdentifierType: ALERTID,
		Keys:           "key1",
	}
	err = removeDetailsRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	removeDetailsRequest := &RemoveDetailsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
		Keys:            "key1,key2",
	}
	err = removeDetailsRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestUpdatePriorityRequest_Validate(t *testing.T) {
	updatePriorityRequestWithoutPriority := &UpdatePriorityRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Priority:        "",
	}
	err := updatePriorityRequestWithoutPriority.Validate()

	assert.Equal(t, err.Error(), errors.New("Priority can not be empty").Error())

	updatePriorityRequestWithInvalidPriority := &UpdatePriorityRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Priority:        "P8",
	}
	err = updatePriorityRequestWithInvalidPriority.Validate()

	assert.Equal(t, err.Error(), errors.New("Priority should be one of these: 'P1', 'P2', 'P3', 'P4' and 'P5'").Error())

	updatePriorityRequestWithoutIdentifier := &UpdatePriorityRequest{
		IdentifierType: ALERTID,
		Priority:       P1,
	}
	err = updatePriorityRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	updatePriorityRequest := &UpdatePriorityRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
		Priority:        P1,
	}
	err = updatePriorityRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestUpdateMessageRequest_Validate(t *testing.T) {
	updateMessageRequestWithoutMessage := &UpdateMessageRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Message:         "",
	}
	err := updateMessageRequestWithoutMessage.Validate()

	assert.Equal(t, err.Error(), errors.New("Message can not be empty").Error())

	updateMessageRequestWithoutIdentifier := &UpdateMessageRequest{
		IdentifierType: ALERTID,
		Message:        "mess1",
	}
	err = updateMessageRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	updateMessageRequest := &UpdateMessageRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
		Message:         "mess1",
	}
	err = updateMessageRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestUpdateDescriptionRequest_Validate(t *testing.T) {
	updateDescriptionRequestWithoutDescription := &UpdateDescriptionRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Description:     "",
	}
	err := updateDescriptionRequestWithoutDescription.Validate()

	assert.Equal(t, err.Error(), errors.New("Description can not be empty").Error())

	updateDescriptionRequestWithoutIdentifier := &UpdateDescriptionRequest{
		IdentifierType: ALERTID,
		Description:    "desc",
	}
	err = updateDescriptionRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	updateDescriptionRequest := &UpdateDescriptionRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
		Description:     "desc",
	}
	err = updateDescriptionRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestListAlertRecipientsRequest_Validate(t *testing.T) {
	listAlertRecipientsRequestWithError := &ListAlertRecipientRequest{}
	err := listAlertRecipientsRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	listAlertRecipientRequest := &ListAlertRecipientRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
	}
	err = listAlertRecipientRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestListAlertLogsRequest_Validate(t *testing.T) {
	listAlertLogsRequestWithError := &ListAlertLogsRequest{}
	err := listAlertLogsRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	alertLogsRequest := &ListAlertLogsRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
	}
	err = alertLogsRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestListAlertNotesRequest_Validate(t *testing.T) {
	listAlertNotesRequestWithError := &ListAlertNotesRequest{}
	err := listAlertNotesRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	listAlertNotesRequest := &ListAlertNotesRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "id1",
	}
	err = listAlertNotesRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestGetSavedSearch_Validate(t *testing.T) {
	getSavedSearchRequestWithError := &GetSavedSearchRequest{}
	err := getSavedSearchRequestWithError.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	getSavedSearchRequest := &GetSavedSearchRequest{
		IdentifierType:  NAME,
		IdentifierValue: "name1",
	}
	err = getSavedSearchRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestCreateAlertAttachmentRequest_Validate(t *testing.T) {
	createAlertAttachmentsRequestWithoutFileName := &CreateAlertAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
		FileName:        "",
		FilePath:        "/usr/",
	}
	err := createAlertAttachmentsRequestWithoutFileName.Validate()

	assert.Equal(t, err.Error(), errors.New("FileName can not be empty").Error())

	createAlertAttachmentsRequestWithoutFilePath := &CreateAlertAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
		FileName:        "test.txt",
		FilePath:        "",
	}
	err = createAlertAttachmentsRequestWithoutFilePath.Validate()

	assert.Equal(t, err.Error(), errors.New("FilePath can not be empty").Error())

	createAlertAttachmentsRequestWithoutIdentifier := &CreateAlertAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "",
		FileName:        "test.txt",
		FilePath:        "/usr/",
	}
	err = createAlertAttachmentsRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	createAlertAttachmentRequest := &CreateAlertAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
		FileName:        "test.txt",
		FilePath:        "/usr/",
	}
	err = createAlertAttachmentRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestGetAlertAttachmentRequest_Validate(t *testing.T) {
	getAttachmentRequestWithoutFileName := &GetAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
		AttachmentId:    "",
	}
	err := getAttachmentRequestWithoutFileName.Validate()

	assert.Equal(t, err.Error(), errors.New("AttachmentId can not be empty").Error())

	getAttachmentRequestWithoutIdentifier := &GetAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "",
		AttachmentId:    "123123",
	}
	err = getAttachmentRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	getAttachmentRequest := &GetAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
		AttachmentId:    "123123",
	}
	err = getAttachmentRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestListAttachmentsRequest_Validate(t *testing.T) {
	listAttachmentsRequestWithoutIdentifier := &ListAttachmentsRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "",
	}
	err := listAttachmentsRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	listAttachmentsRequest := &ListAttachmentsRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
	}
	err = listAttachmentsRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestDeleteAlertAttachmentRequest_Validate(t *testing.T) {
	deleteAttachmentRequestWithoutFileName := &DeleteAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
		AttachmentId:    "",
	}
	err := deleteAttachmentRequestWithoutFileName.Validate()

	assert.Equal(t, err.Error(), errors.New("AttachmentId can not be empty").Error())

	deleteAttachmentRequestWithoutIdentifier := &DeleteAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "",
		AttachmentId:    "123123",
	}
	err = deleteAttachmentRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	deleteAttachmentRequest := &DeleteAttachmentRequest{
		IdentifierType:  TINYID,
		IdentifierValue: "tiny1",
		AttachmentId:    "123123",
	}
	err = deleteAttachmentRequest.Validate()

	assert.Equal(t, err, nil)
}

func TestAddResponderRequest_Validate(t *testing.T) {
	addResponderRequestWithInvalidResponderType := &AddResponderRequest{
		IdentifierType:  ALERTID,
		IdentifierValue: "Id",
		Responder: Responder{
			Type: "escalation",
			Name: "Test",
		},
	}
	err := addResponderRequestWithInvalidResponderType.Validate()

	assert.Equal(t, err.Error(), errors.New("Responder type must be user or team").Error())

	addResponderRequestWithInvalidUser := &AddResponderRequest{
		IdentifierType: ALERTID,
		Responder: Responder{
			Type: "user",
		},
	}
	err = addResponderRequestWithInvalidUser.Validate()

	assert.Equal(t, err.Error(), errors.New("User ID or username must be defined").Error())

	addResponderRequestWithInvalidTeam := &AddResponderRequest{
		IdentifierType: ALERTID,
		Responder: Responder{
			Type: "team",
		},
	}
	err = addResponderRequestWithInvalidTeam.Validate()

	assert.Equal(t, err.Error(), errors.New("Team ID or name must be defined").Error())

	addResponderRequestWithoutIdentifier := &AddResponderRequest{
		IdentifierType: ALERTID,
		Responder: Responder{
			Type:     "user",
			Username: "usertest1",
		},
	}
	err = addResponderRequestWithoutIdentifier.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	addResponderRequest := &AddResponderRequest{
		IdentifierValue: "tiny1",
		IdentifierType:  TINYID,
		Responder: Responder{
			Type: "team",
			Id:   "id1",
		},
	}
	err = addResponderRequest.Validate()

	assert.Equal(t, err, nil)
}
