---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_service"
sidebar_current: "docs-opsgenie-datasource-service"
description: |-
  Manages existing Service within Opsgenie.
---

# opsgenie\_service

Manages existing Service within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_team" "payment" {
  name        = "example"
  description = "This team deals with all the things"
}

resource "opsgenie_service" "this" {
  name  = "Payment"
  team_id = "$opsgenie_team.payment.id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the service. This field must not be longer than 100 characters.

* `team_id` - (Required)  Team id of the service. This field must not be longer than 512 characters.

The following attributes are exported:

* `id` - The ID of the Opsgenie Service.

* `team_id` - Team id of the service.

* `description` - Description field of the service that is generally used to provide a detailed information about the service.
