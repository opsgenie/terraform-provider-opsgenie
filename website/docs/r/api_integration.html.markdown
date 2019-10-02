---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_api_integration"
sidebar_current: "docs-opsgenie-resource-api-integration"
description: |-
  Manages an API Integration within Opsgenie.
---

# opsgenie_api_integration

Manages an API Integration within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_api_integration" "example-api-integration" {
  name = "api-based-int"
  type = "API"

  responders {
    type = "user"
    id   = "${opsgenie_user.user.id}"
  }

  responders {
    type = "user"
    id   = "${opsgenie_user.fahri.id}"
  }
}

resource "opsgenie_api_integration" "example-api-integration" {
  name = "api-based-int-2"
  type = "Prometheus"

  responders {
    type = "user"
    id   = "${opsgenie_user.user.id}"
  }

  enabled                        = false
  allow_write_access             = false
  ignore_responders_from_payload = true
  suppress_notifications         = true
  owner_team_id                  = "${opsgenie_team_genies.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the integration. Name must be unique for each integration.

* `type` - (Optional) Type of the integration. (API,Marid,Prometheus ...)

* `allow_write_access` - (Optional) This parameter is for configuring the write access of integration. If write access is restricted, the integration will not be authorized to write within any domain. Defaults to true.

* `enabled` - (Optional) This parameter is for specifying whether the integration will be enabled or not. Defaults to true

* `ignore_responders_from_payload` - (Optional) If enabled, the integration will ignore recipients sent in request payloads. Defaults to false.

* `suppress_notifications` - (Optional) If enabled, notifications that come from alerts will be suppressed. Defaults to false.

* `owner_team_id` - (Optional) Owner team id of the integration.

* `responder` - (Optional)  User, schedule, teams or escalation names to calculate which users will receive the notifications of the alert.

`responder` supports the following:

* `type` - (Required) The responder type.
* `id` - (Required) The id of the responder.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie API Integration.

* `api_key` - (Computed) API key of the created integration

## Import

API Integrations can be imported using the `id`, e.g.

`$ terraform import opsgenie_team.team1 812be1a1-32c8-4666-a7fb-03ecc385106c`
