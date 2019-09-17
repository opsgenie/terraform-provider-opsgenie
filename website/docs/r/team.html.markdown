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
  username  = "eggman@dr-robotnik.com"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name associated with this team. Opsgenie defines that this must not be longer than 100 characters.

* `description` - (Optional) A description for this team.

* `member` - (Optional) A Member block as documented below.

`member` supports the following:

* `id` - (Required) The UUID for the member to add to this Team.
* `role` - (Optional) The role for the user within the Team - can be either 'admin' or 'user', defaults to 'user' if not set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Team.

## Import

Teams can be imported using the `id`, e.g.

```
$ terraform import opsgenie_team.team1 812be1a1-32c8-4666-a7fb-03ecc385106c
```
