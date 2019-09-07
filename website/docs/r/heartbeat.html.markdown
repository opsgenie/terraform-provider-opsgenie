---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_heartbeat"
sidebar_current: "docs-opsgenie-resource-heartbeat"
description: |-
  Manages Heartbeat within Opsgenie.
---

# opsgenie_heartbeat

Manages an heartbeat within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_heartbeat" "test" {
	name = "genieheartbeat-%s"
	description = "test opsgenie heartbeat terraform"
	interval_unit = "minutes"
	interval = 10
	enabled = false
	alert_message = "Test"
	alert_priority = "P3"
	alert_tags = ["test","fahri"]
	owner_team_id = "${opsgenie_team.test.id}"

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the integration. Name must be unique for each integration. 

* `description` - (Optional) Type of the integration. (API,Marid,Prometheus ...)

* `interval_unit` - (Optional) This parameter is for configuring the write access of integration. If write access is restricted, the integration will not be authorized to write within any domain. Defaults to true.

* `interval` - (Optional) This parameter is for specifying whether the integration will be enabled or not. Defaults to true

* `enabled` - (Optional) If enabled, the integration will ignore recipients sent in request payloads. Defaults to false.

* `owner_team_id` - (Optional) Owner team id of the integration.

* `alert_message` - (Optional) If enabled, notifications that come from alerts will be suppressed. Defaults to false.

* `alert_priority` - (Optional)  User, schedule, teams or escalation names to calculate which users will receive the notifications of the alert.

* `alert_tags` - (Optional)  User, schedule, teams or escalation names to calculate which users will receive the notifications of the alert.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Heartbeat.

## Import

Heartbeat Integrations can be imported using the `id`, e.g.

```
$ terraform import opsgenie_heartbeat.test 812be1a1-32c8-4666-a7fb-03ecc385106c
```
