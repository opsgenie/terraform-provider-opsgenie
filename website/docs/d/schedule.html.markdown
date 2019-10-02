---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_schedule"
sidebar_current: "docs-opsgenie-resource-schedule"
description: |-
  Manages a Schedule within Opsgenie.
---

# opsgenie_schedule

Manages a Schedule within Opsgenie.

## Example Usage
```hcl
resource "opsgenie_schedule" "test" {
  name        = "genieschedule-%s"
  description = "schedule test"
  timezone    = "Europe/Rome"
  enabled     = false
}

resource "opsgenie_schedule" "test" {
  name          = "genieschedule-%s"
  description   = "schedule test"
  timezone      = "Europe/Rome"
  enabled       = false
  owner_team_id = "${opsgenie_team.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the schedule.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Schedule.

* `rules` - A Member block as documented below.

* `description` - Timezone of schedule. Please look at [Supported Timezone Ids](https://docs.opsgenie.com/docs/supported-timezone-ids) for available timezones - Defaults to "America/New_York".

* `timezone` - The description of schedule.

* `enabled` - Enable/disable state of schedule

* `owner_team_id` - Owner team id of the schedule.
