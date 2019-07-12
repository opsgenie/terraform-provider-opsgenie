package user

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"time"
)

type User struct {
	Id          string              `json:"id"`
	Username    string              `json:"username"`
	FullName    string              `json:"fullName"`
	Role        *UserRole           `json:"role"`
	Blocked     bool                `json:"blocked"`
	Verified    bool                `json:"verified"`
	UserAddress *UserAddress        `json:"userAddress"`
	Tags        []string            `json:"tags"`
	Details     map[string][]string `json:"details"`
	TimeZone    string              `json:"timeZone"`
	Locale      string              `json:"locale"`
	CreatedAt   time.Time           `json:"createdAt"`
}

type UserRole struct {
	RoleName string `json:"name"`
}

type UserAddress struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Line    string `json:"line"`
	ZipCode string `json:"zipCode"`
}

type UserContact struct {
	To            string `json:"to"`
	Id            string `json:"id"`
	ContactMethod string `json:"contactMethod"`
	Enabled       bool   `json:"enabled"`
}

type CreateResult struct {
	client.ResultMetadata
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
	Id      string            `json:"id"`
	Name    string            `json:"name"`
}

type GetResult struct {
	client.ResultMetadata
	Id            string              `json:"id"`
	Username      string              `json:"username"`
	FullName      string              `json:"fullName"`
	Role          *UserRole           `json:"role"`
	Blocked       bool                `json:"blocked"`
	Verified      bool                `json:"verified"`
	UserAddress   *UserAddress        `json:"userAddress"`
	SkypeUsername string              `json:"skypeUsername"`
	Tags          []string            `json:"tags"`
	Details       map[string][]string `json:"details"`
	TimeZone      string              `json:"timeZone"`
	Locale        string              `json:"locale"`
	CreatedAt     time.Time           `json:"createdAt"`
	UserContacts  []UserContact       `json:"userContacts"`
}

type UpdateResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type ListResult struct {
	client.ResultMetadata
	Users      []User `json:"data"`
	Paging     Paging `json:"paging"`
	TotalCount int    `json:"totalCount"`
}

type Paging struct {
	Next  string `json:"next"`
	Prev  string `json:"prev"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type ListUserEscalationsResult struct {
	client.ResultMetadata
	Escalations []UserEscalation `json:"data"`
}

type UserEscalation struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       []Rule `json:"rules"`
	OwnerTeam   Team   `json:"ownerTeam"`
}

type Rule struct {
	Condition  og.EscalationCondition `json:"condition"`
	NotifyType og.NotifyType          `json:"notifyType"`
	Recipient  og.Participant         `json:"recipient"`
	Delay      EscalationDelay        `json:"delay"`
}

type EscalationDelay struct {
	TimeUnit   og.TimeUnit `json:"timeUnit"`
	TimeAmount int         `json:"timeAmount"`
}

type ListUserTeamsResult struct {
	client.ResultMetadata
	Teams []Team `json:"data"`
}

type Team struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListUserForwardingRulesResult struct {
	client.ResultMetadata
	ForwardingRules []ForwardingRule `json:"data"`
}

type ForwardingRule struct {
	FromUser  ForwardedUser `json:"fromUser"`
	ToUser    ForwardedUser `json:"toUser"`
	StartDate time.Time     `json:"startDate"`
	EndDate   time.Time     `json:"endDate"`
	Alias     string        `json:"alias"`
	Id        string        `json:"id"`
}

type ForwardedUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type ListUserSchedulesResult struct {
	client.ResultMetadata
	Schedules []Schedule `json:"data"`
}

type Schedule struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type GetSavedSearchResult struct {
	client.ResultMetadata
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerId     string `json:"ownerId"`
	Query       string `json:"query"`
}

type ListSavedSearchesResult struct {
	client.ResultMetadata
	UsersSavedSearches []UsersSavedSearch `json:"data"`
}

type UsersSavedSearch struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DeleteSavedSearchResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}
