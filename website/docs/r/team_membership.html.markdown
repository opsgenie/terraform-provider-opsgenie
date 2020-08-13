---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_team_membership"
sidebar_current: "docs-opsgenie-resource-team-membership"
description: |-
  Manages team memberships for users.
---

# opsgenie\_team

Manages team memberships for users.

## Example Usage

```hcl
resource "opsgenie_user" "first" {
  username  = "user@test.example.com"
  full_name = "name "
  role      = "User"
}

resource "opsgenie_user" "second" {
  username  = "test@test.example.com"
  full_name = "name "
  role      = "User"
}

resource "opsgenie_team" "test" {
  name           = "example"
  description    = "This team deals with all the things"
  ignore_members = true # we're using opsgenie_team_membership for it
}

resource "opsgenie_team_membership" "first" {
  user_id = opsgenie_user.first.id
  role    = "user"
  team_id = opsgenie_team.test.id
}

resource "opsgenie_team_membership" "second" {
  user_id = opsgenie_user.second.id
  role    = "admin"
  team_id = opsgenie_team.test.id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required) The ID of an user.

* `role` - (Optional) The role for the user within the Team - can be either 'admin' or 'user', defaults to 'user' if not set.

* `teamr_id` - (Required) The ID of a team the user should be member of.

## Attributes Reference

The following attributes are exported:

* `id` - The virtual ID of the Opsgenie Team Membership.

## Import

Import is not supported for team memberships.