---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_email_integration"
sidebar_current: "docs-opsgenie-resource-email-integration"
description: |-
  Manages an Email Integration within Opsgenie.
---

# opsgenie_email_integration

Manages an Email Integration within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_email_integration" "test" {
  name           = "genieintegration-name"
  email_username = "fahri"
}

resource "opsgenie_email_integration" "test" {
  name = "genieintegration-%s"

  responders {
    type = "user"
    id   = "${opsgenie_user.test.id}"
  }

  responders {
    type = "schedule"
    id   = "${opsgenie_schedule.test.id}"
  }

  responders {
    type = "escalation"
    id   = "${opsgenie_escalation.test.id}"
  }

  responders {
    type = "team"
    id   = "${opsgenie_team.test2.id}"
  }

  email_username                 = "test"
  enabled                        = true
  ignore_responders_from_payload = true
  suppress_notifications         = true
}


resource "opsgenie_email_integration" "test" {
  name = "genieintegration-%s"

  responders {
    type = "user"
    id   = "${opsgenie_user.test.id}"
  }

  email_username                 = "test"
  enabled                        = true
  ignore_responders_from_payload = true
  suppress_notifications         = true
  owner_team_id                  = "${opsgenie_team_genies.id}"

}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the integration. Name must be unique for each integration.

* `email_username` - (Required) The username part of the email address. It must be unique for each integration.

* `enabled` - (Optional) A Member block as documented below.

* `ignore_responders_from_payload` - (Optional) If enabled, the integration will ignore recipients sent in request payloads. Default: `false`.

* `suppress_notifications` - (Optional) If enabled, notifications that come from alerts will be suppressed. Default: `false`.

* `owner_team_id` - (Optional) Owner team id of the integration.

* `responder` - (Optional) User, schedule, teams or escalation names to calculate which users will receive the notifications of the alert.

`responder` supports the following:

* `type` - (Required) The responder type.
* `id` - (Required) The id of the responder.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Email based Integration.

## Import

Email Integrations can be imported using the `id`, e.g.

`$ terraform import opsgenie_email_integration.test id`
