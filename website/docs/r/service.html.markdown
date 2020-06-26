---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_service"
sidebar_current: "docs-opsgenie-resource-service"
description: |-
  Manages a Service within Opsgenie.
---

# opsgenie\_service

Manages a Service within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_team" "payment" {
  name        = "example"
  description = "This team deals with all the things"
}

resource "opsgenie_service" "this" {
  name  = "Payment"
  team_id = "$opsgenie_team.this.id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the service. This field must not be longer than 100 characters.

* `team_id` - (Required)  Team id of the service. This field must not be longer than 512 characters.

* `description` - (Optional) Description field of the service that is generally used to provide a detailed information about the service.

* `visibility` - (Optional) Teams and users that the service will become visible to. "TEAM\_MEMBERS" would mean that service is visible to only given team members; whereas "OPSGENIE\_USERS" mean that service is visible to all users of the given customers. Defaults to: "TEAM\_MEMBERS".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Service.

## Import

Teams can be imported using the `id`, e.g.

`$ terraform import opsgenie_service.this 812be1a1-32c8-4666-a7fb-03ecc385106c`
