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
  team_id = "$opsgenie_team.payment.id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the service. This field must not be longer than 100 characters.

* `team_id` - (Required)  Team id of the service. This field must not be longer than 512 characters.

* `description` - (Optional) Description field of the service that is generally used to provide a detailed information about the service.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Service.

## Import

Teams can be imported using the `service_id`, e.g.

`$ terraform import opsgenie_service.this service_id`
