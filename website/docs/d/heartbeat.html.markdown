---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_heartbeat"
sidebar_current: "docs-opsgenie-resource-heartbeat"
description: |-
  Manages existing Heartbeats within Opsgenie.
---

# opsgenie_heartbeat

Manages existing heartbeat within Opsgenie.

## Example Usage

```hcl
data "opsgenie_heartbeat" "test" {
  name = "genieheartbeat-existing"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the heartbeat


## Attributes Reference

The following attributes are exported:

* `description` -  An optional description of the heartbeat

* `interval_unit` - Interval specified as minutes, hours or days.

* `interval` - Specifies how often a heartbeat message should be expected.

* `enabled` -  Enable/disable heartbeat monitoring.

* `owner_team_id` - Owner team of the heartbeat.

* `alert_message` - Specifies the alert message for heartbeat expiration alert. If this is not provided, default alert message is "HeartbeatName is expired".

* `alert_priority` - Specifies the alert priority for heartbeat expiration alert. If this is not provided, default priority is P3.

* `alert_tags` - Specifies the alert tags for heartbeat expiration alert.

