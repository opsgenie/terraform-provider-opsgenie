---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_escalation"
sidebar_current: "docs-opsgenie-resource-escalation"
description: |-
  Manages an Escalation within Opsgenie.
---

# opsgenie_escalation

Manages an Escalation within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_escalation" "test" {
  name = "genieescalation-%s"

  rules {
    condition   = "if-not-acked"
    notify_type = "default"
    delay       = 1

    recipient {
      type = "user"
      id   = "${opsgenie_user.test.id}"
		}
	}
}

resource "opsgenie_escalation" "test" {
  name          = "genieescalation-%s"
  description   = "test"
  owner_team_id = "${opsgenie_team.test.id}"

  rules {
    condition   = "if-not-acked"
    notify_type = "default"
    delay       = 1

    recipient {
      type = "user"
      id   = "${opsgenie_user.test.id}"
    }

    recipient {
      type = "team"
      id   = "${opsgenie_team.test.id}"
    }

    recipient {
      type = "schedule"
      id   = "${opsgenie_schedule.test.id}"
    }
  }

  repeat  {
    wait_interval          = 10
    count                  = 1
    reset_recipient_states = true
    close_alert_after_all  = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the escalation.

* `rules` - (Required) List of the escalation rules.

* `description` - (Optional) Description of the escalation.

* `owner_team_id` - (Optional) Owner team id of the escalation.

* `repeat` - (Optional) Repeat preferences of the escalation including repeat interval, count, reverting acknowledge and seen states back and closing an alert automatically as soon as repeats are completed


`rules` supports the following:

* `condition` - (Required) The condition for notifying the recipient of escalation rule that is based on the alert state. Possible values are: `if-not-acked` and `if-not-closed`. Default: `if-not-acked`
* `notify_type` - (Required) Recipient calculation logic for schedules. Possible values are:

  - `default`: on call users
  - `next`: next users in rotation
  - `previous`: previous users on rotation
  - `users`: users of the team
  - `admins`: admins of the team
  - `all`: all members of the team

* `recipient` - (Required) Object of schedule, team, or users which will be notified in escalation. The possible values for participants are: `user`, `schedule`, `team`.
* `delay` - (Required) Time delay of the escalation rule, in minutes.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Escalation.

## Import

Escalations can be imported using the `escalation_id`, e.g.

`$ terraform import opsgenie_escalation.test escalation_id`

