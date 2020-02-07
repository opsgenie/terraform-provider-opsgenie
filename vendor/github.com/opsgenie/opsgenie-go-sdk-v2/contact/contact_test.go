package contact

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequest_Validate(t *testing.T) {
	var err error
	createRequest := &CreateRequest{}
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	createRequest.UserIdentifier = "123"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("to cannot be empty.").Error())

	createRequest.To = "ab@email.com"
	err = createRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Method cannot be empty.").Error())

	createRequest.MethodOfContact = Email
	err = createRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, createRequest.ResourcePath(), "/v2/users/123/contacts")
	assert.Equal(t, createRequest.Method(), http.MethodPost)

}

func TestGetRequest_Validate(t *testing.T) {
	var err error
	getRequest := GetRequest{}
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	getRequest.UserIdentifier = "123"
	err = getRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Contact identifier cannot be empty.").Error())

	getRequest.ContactIdentifier = "1234"
	err = getRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, getRequest.ResourcePath(), "/v2/users/123/contacts/1234")
	assert.Equal(t, getRequest.Method(), http.MethodGet)

}

func TestUpdateRequest_Validate(t *testing.T) {
	var err error
	updateRequest := &UpdateRequest{}
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	updateRequest.UserIdentifier = "123"
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Contact identifier cannot be empty.").Error())

	updateRequest.ContactIdentifier = "1234"
	err = updateRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("to cannot be empty.").Error())

	updateRequest.To = "ab@email.com"
	err = updateRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, updateRequest.ResourcePath(), "/v2/users/123/contacts/1234")
	assert.Equal(t, updateRequest.Method(), http.MethodPatch)

}
func TestDeleteRequest_Validate(t *testing.T) {
	var err error
	deleteRequest := &DeleteRequest{}
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	deleteRequest.UserIdentifier = "123"
	err = deleteRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Contact identifier cannot be empty.").Error())

	deleteRequest.ContactIdentifier = "1234"
	err = deleteRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, deleteRequest.ResourcePath(), "/v2/users/123/contacts/1234")
	assert.Equal(t, deleteRequest.Method(), http.MethodDelete)

}

func TestListRequest_Validate(t *testing.T) {
	var err error
	listRequest := &ListRequest{}
	err = listRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	listRequest.UserIdentifier = "123"
	err = listRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, listRequest.ResourcePath(), "/v2/users/123/contacts")
	assert.Equal(t, listRequest.Method(), http.MethodGet)

}

func TestEnableRequest_Validate(t *testing.T) {
	var err error
	enableRequest := &EnableRequest{}
	err = enableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	enableRequest.UserIdentifier = "123"
	err = enableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Contact identifier cannot be empty.").Error())

	enableRequest.ContactIdentifier = "1234"
	err = enableRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, enableRequest.ResourcePath(), "/v2/users/123/contacts/1234/enable")
	assert.Equal(t, enableRequest.Method(), http.MethodPost)

}

func TestDisableRequest_Validate(t *testing.T) {
	var err error
	disableRequest := &DisableRequest{}
	err = disableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User identifier cannot be empty.").Error())

	disableRequest.UserIdentifier = "123"
	err = disableRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("Contact identifier cannot be empty.").Error())

	disableRequest.ContactIdentifier = "1234"
	err = disableRequest.Validate()
	assert.Nil(t, err)

	assert.Equal(t, disableRequest.ResourcePath(), "/v2/users/123/contacts/1234/disable")
	assert.Equal(t, disableRequest.Method(), http.MethodPost)

}
