package user

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	userRequest := &CreateRequest{
		Username: "",
		FullName: "Name Surname",
		Role: &UserRoleRequest{
			RoleName: "admin",
		},
	}
	err := userRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Username can not be empty").Error())

	userRequest = &CreateRequest{
		Username: "Name@gmail.com",
		FullName: "",
		Role: &UserRoleRequest{
			RoleName: "admin",
		},
	}
	err = userRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("FullName can not be empty").Error())

	userRequest = &CreateRequest{
		Username: "Name@gmail.com",
		FullName: "Name Surname",
		Role:     &UserRoleRequest{},
	}
	err = userRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User Role can not be empty").Error())

	userRequest = &CreateRequest{
		Username: "Name@gmail.com",
		FullName: "Name Surname",
		Role: &UserRoleRequest{
			RoleName: "admin",
		},
	}
	err = userRequest.Validate()
	assert.Equal(t, err, nil)

}

func TestGetUserRequest_Validate(t *testing.T) {
	userRequest := &GetRequest{
		Identifier: "",
	}
	err := userRequest.Validate()
	reqParam := userRequest.RequestParams()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())
	assert.Equal(t, len(reqParam), 0)

	userRequest = &GetRequest{
		Identifier: "id1",
		Expand:     "company",
	}
	err = userRequest.Validate()
	reqParam = userRequest.RequestParams()

	assert.Equal(t, err, nil)
	assert.Equal(t, reqParam["expand"], "company")
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	userRequest := &UpdateRequest{
		Identifier: "",
	}
	err := userRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	userRequest = &UpdateRequest{
		Identifier: "id1",
		Role: &UserRoleRequest{
			RoleName: "",
		},
	}

	err = userRequest.Validate()
	assert.Equal(t, err.Error(), errors.New("User role name can not be empty").Error())

	userRequest = &UpdateRequest{
		Identifier: "id1",
	}

	err = userRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestDeleteUserRequest_Validate(t *testing.T) {
	userRequest := &DeleteRequest{
		Identifier: "",
	}
	err := userRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	userRequest = &DeleteRequest{
		Identifier: "id1",
	}

	err = userRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestListUserRequest_Validate(t *testing.T) {
	userRequest := &ListRequest{}
	reqParam := userRequest.RequestParams()

	assert.Equal(t, len(reqParam), 0)

	userRequest = &ListRequest{
		Limit:  1,
		Offset: 2,
		Sort:   Username,
		Order:  Asc,
		Query:  "query:1",
	}
	reqParam = userRequest.RequestParams()

	assert.Equal(t, reqParam["limit"], "1")
	assert.Equal(t, reqParam["offset"], "2")
	assert.Equal(t, reqParam["sort"], "username")
	assert.Equal(t, reqParam["order"], "asc")
	assert.Equal(t, reqParam["query"], "query:1")
}

func TestListUserEscalationsRequest_Validate(t *testing.T) {
	userRequest := &ListUserEscalationsRequest{
		Identifier: "",
	}
	err := userRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	userRequest = &ListUserEscalationsRequest{
		Identifier: "id1",
	}

	err = userRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestListUserTeamsRequest_Validate(t *testing.T) {
	userRequest := &ListUserTeamsRequest{
		Identifier: "",
	}
	err := userRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	userRequest = &ListUserTeamsRequest{
		Identifier: "id1",
	}

	err = userRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestListUserForwardingRulesRequest_Validate(t *testing.T) {
	userRequest := &ListUserForwardingRulesRequest{
		Identifier: "",
	}
	err := userRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	userRequest = &ListUserForwardingRulesRequest{
		Identifier: "id1",
	}

	err = userRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestListUserSchedulesRequest_Validate(t *testing.T) {
	userRequest := &ListUserSchedulesRequest{
		Identifier: "",
	}
	err := userRequest.Validate()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())

	userRequest = &ListUserSchedulesRequest{
		Identifier: "id1",
	}

	err = userRequest.Validate()
	assert.Equal(t, err, nil)
}

func TestGetUsersSavedSearchRequest_Validate(t *testing.T) {
	userRequest := &GetSavedSearchRequest{
		Identifier: "",
	}
	err := userRequest.Validate()
	reqParam := userRequest.RequestParams()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())
	assert.Equal(t, reqParam["identifierType"], "id")

	userRequest = &GetSavedSearchRequest{
		Identifier:     "id1",
		IdentifierType: Name,
	}
	err = userRequest.Validate()
	reqParam = userRequest.RequestParams()

	assert.Equal(t, err, nil)
	assert.Equal(t, reqParam["identifierType"], "name")
}

func TestDeleteUsersSavedSearchRequest_Validate(t *testing.T) {
	userRequest := &DeleteSavedSearchRequest{
		Identifier: "",
	}
	err := userRequest.Validate()
	reqParam := userRequest.RequestParams()

	assert.Equal(t, err.Error(), errors.New("Identifier can not be empty").Error())
	assert.Equal(t, reqParam["identifierType"], "id")

	userRequest = &DeleteSavedSearchRequest{
		Identifier:     "id1",
		IdentifierType: Name,
	}
	err = userRequest.Validate()
	reqParam = userRequest.RequestParams()

	assert.Equal(t, err, nil)
	assert.Equal(t, reqParam["identifierType"], "name")
}
