---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_team"
sidebar_current: "docs-opsgenie-resource-team"
description: |-
  Manages a Team within Opsgenie.
---

# opsgenie\_team

Manages a Team within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_user" "first" {
  username  = "user@domain.com"
  full_name = "name "
  role      = "User"
}

resource "opsgenie_user" "second" {
  username  = "test@domain.com"
  full_name = "name "
  role      = "User"
}

resource "opsgenie_team" "test" {
  name        = "example"
  description = "This team deals with all the things"

  member {
    id   = "${opsgenie_user.first.id}"
    role = "admin"
  }

  member {
    id   = "${opsgenie_user.second.id}"
    role = "user"
  }
}

resource "opsgenie_team" "self-service" {
  name           = "Self Service"
  description    = "Membership in this team is managed via OpsGenie web UI only"
  ignore_members = true
  delete_default_resources = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name associated with this team. Opsgenie defines that this must not be longer than 100 characters.

* `description` - (Optional) A description for this team.

* `ignore_members` - (Optional) Set to true to ignore any configured member blocks and any team member added/updated/removed via OpsGenie web UI. Use this option e.g. to maintain membership via web UI only and use it only for new teams. Changing the value for existing teams might lead to strange behaviour. Default: `false`.

* `delete_default_resources` - (Optional) Set to true to remove default escalation and schedule for newly created team. **Be careful its also changes that team routing rule to None. That means you have to define routing rule as well**


* `member` - (Optional) A Member block as documented below.

`member` supports the following:

* `id` - (Required) The UUID for the member to add to this Team.
* `role` - (Optional) The role for the user within the Team - can be either `admin` or `user`. Default: `user`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Team.

## Import

Teams can be imported using the `team_id`, e.g.

`$ terraform import opsgenie_team.team1 team_id`
