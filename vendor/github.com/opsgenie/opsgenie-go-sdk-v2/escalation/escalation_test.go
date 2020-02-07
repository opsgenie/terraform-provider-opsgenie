package escalation

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRequest_Validate(t *testing.T) {
	createRequest := &CreateRequest{

		Name: "",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Default,
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
	}
	err := createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Name cannot be empty.").Error())

	createRequest = &CreateRequest{
		Name: "name1",
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rules list cannot be empty.").Error())

	createRequest = &CreateRequest{

		Name: "name1",
		Rules: []RuleRequest{
			{
				Condition:  "",
				NotifyType: og.Default,
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule Condition should be one of these: 'if-not-acked', 'if-not-closed'.").Error())

	createRequest = &CreateRequest{
		Name: "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: "",
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Notify Type should be one of these: 'next', 'previous', 'default', 'users', 'admins', 'all'.").Error())

	createRequest = &CreateRequest{
		Name: "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Name: "test",
				},
			},
		},
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Recipient type cannot be empty.").Error())

	createRequest = &CreateRequest{
		Name: "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Type: og.Escalation,
					Name: "test",
				},
			},
		},
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Recipient type should be one of these: 'User', 'Team', 'Schedule'").Error())

	createRequest = &CreateRequest{
		Name: "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Type: og.User,
					Name: "test",
				},
			},
		},
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("For recipient type user either username or id must be provided.").Error())

	createRequest = &CreateRequest{
		Name: "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Type:     og.Schedule,
					Username: "test",
				},
			},
		},
	}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("For recipient type team and schedule either name or id must be provided.").Error())

	createRequest = &CreateRequest{
		Name: "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
	}
	err = createRequest.Validate()
	assert.Equal(t, err, nil)

}

func TestUpdateRequest_Validate(t *testing.T) {
	updateRequest := &UpdateRequest{

		Name: "",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Default,
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
		IdentifierType: Name,
		Identifier:     "",
	}
	err := updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Identifier cannot be empty.").Error())

	updateRequest = &UpdateRequest{

		Name: "",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Default,
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
		IdentifierType: "anotherIdentifier",
		Identifier:     "id1",
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Identifier Type should be one of this : 'id', 'name' or empty.").Error())

	updateRequest = &UpdateRequest{
		IdentifierType: Name,
		Identifier:     "test",
		Name:           "name1",
		Rules: []RuleRequest{
			{
				Condition:  "",
				NotifyType: og.Default,
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Rule Condition should be one of these: 'if-not-acked', 'if-not-closed'.").Error())

	updateRequest = &UpdateRequest{
		IdentifierType: Name,
		Identifier:     "test",
		Name:           "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: "",
				Recipient: og.Participant{
					Type: og.Schedule,
					Name: "test",
				},
			},
		},
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Notify Type should be one of these: 'next', 'previous', 'default', 'users', 'admins', 'all'.").Error())

	updateRequest = &UpdateRequest{
		IdentifierType: Name,
		Identifier:     "test",
		Name:           "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Name: "test",
				},
			},
		},
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Recipient type cannot be empty.").Error())

	updateRequest = &UpdateRequest{
		IdentifierType: Name,
		Identifier:     "test",
		Name:           "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Type: og.Escalation,
					Name: "test",
				},
			},
		},
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Recipient type should be one of these: 'User', 'Team', 'Schedule'").Error())

	updateRequest = &UpdateRequest{
		IdentifierType: Name,
		Identifier:     "test",
		Name:           "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Type: og.User,
					Name: "test",
				},
			},
		},
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("For recipient type user either username or id must be provided.").Error())

	updateRequest = &UpdateRequest{
		IdentifierType: Name,
		Identifier:     "test",
		Name:           "name1",
		Rules: []RuleRequest{
			{
				Condition:  og.IfNotAcked,
				NotifyType: og.Admins,
				Recipient: og.Participant{
					Type:     og.Schedule,
					Username: "test",
				},
			},
		},
	}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("For recipient type team and schedule either name or id must be provided.").Error())

	updateRequest = &UpdateRequest{
		IdentifierType: Name,
		Identifier:     "test",
	}
	err = updateRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestDeleteRequest_Validate(t *testing.T) {
	deleteRequest := &DeleteRequest{
		IdentifierType: Name,
		Identifier:     "",
	}
	err := deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Identifier cannot be empty.").Error())

	deleteRequest = &DeleteRequest{
		IdentifierType: "anotherIdentifier",
		Identifier:     "id1",
	}
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Identifier Type should be one of this : 'id', 'name' or empty.").Error())

	deleteRequest = &DeleteRequest{
		IdentifierType: Name,
		Identifier:     "test",
	}
	err = deleteRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestGetRequest_Validate(t *testing.T) {
	getRequest := &GetRequest{
		IdentifierType: Name,
		Identifier:     "",
	}
	err := getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Identifier cannot be empty.").Error())
	getRequest = &GetRequest{
		IdentifierType: "anotherIdentifier",
		Identifier:     "id1",
	}
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Identifier Type should be one of this : 'id', 'name' or empty.").Error())

	getRequest = &GetRequest{
		IdentifierType: Name,
		Identifier:     "test",
	}
	err = getRequest.Validate()
	assert.Equal(t, err, nil)
}
