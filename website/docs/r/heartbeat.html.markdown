---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_heartbeat"
sidebar_current: "docs-opsgenie-resource-heartbeat"
description: |-
  Manages Heartbeat within Opsgenie.
---

# opsgenie_heartbeat

Manages heartbeat within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_heartbeat" "test" {
	name           = "genieheartbeat-%s"
	description    = "test opsgenie heartbeat terraform"
	interval_unit  = "minutes"
	interval       = 10
	enabled        = false
	alert_message  = "Test"
	alert_priority = "P3"
	alert_tags     = ["test","fahri"]
	owner_team_id  = "${opsgenie_team.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the heartbeat

* `description` - (Optional) An optional description of the heartbeat

* `interval_unit` - (Required) Interval specified as minutes, hours or days.

* `interval` - (Required) Specifies how often a heartbeat message should be expected.

* `enabled` - (True) Enable/disable heartbeat monitoring.

* `owner_team_id` - (Optional) Owner team of the heartbeat.

* `alert_message` - (Optional) Specifies the alert message for heartbeat expiration alert. If this is not provided, default alert message is "HeartbeatName is expired".

* `alert_priority` - (Optional) Specifies the alert priority for heartbeat expiration alert. If this is not provided, default priority is P3.

* `alert_tags` - (Optional)  Specifies the alert tags for heartbeat expiration alert.


## Attributes Reference

Only the arguments listed above are exposed as attributes.


## Import

Heartbeat Integrations can be imported using the `name`, e.g.

`$ terraform import opsgenie_heartbeat.test name`
