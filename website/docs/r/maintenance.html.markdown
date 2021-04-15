---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_maintenance"
sidebar_current: "docs-opsgenie-resource-maintenance"
description: |-
  Manages a Maintenance within Opsgenie.
---

# opsgenie_maintenance

Manages a Maintenance within Opsgenie.

## Example Usage
```hcl
resource "opsgenie_maintenance" "test" {
  description = "geniemaintenance-%s"

  time {
    type       = "schedule"
    start_date = "2019-06-20T17:45:00Z"
    end_date   = "2019-06-20T17:50:00Z"
  }

  rules{}
}

resource "opsgenie_maintenance" "test" {
  description = "geniemaintenance-%s"

  time {
    type       = "schedule"
    start_date = "2019-06-20T17:45:00Z"
    end_date   = "2019-06-%dT17:50:00Z"
  }

  rules {
    state = "enabled"

    entity {
      id   = "${opsgenie_email_integration.test.id}"
      type = "integration"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `time` - (Required) Time configuration of maintenance. It takes a time object which has type, startDate and endDate fields

* `rules` - (Required) Rules of maintenance, which takes a list of rule objects and defines the maintenance rules over integrations and policies.

* `description` - (Optional) Description for the maintenance.


`times` supports the following:

* `type` - (Required) This parameter defines when the maintenance will be active. It can take one of for-5-minutes, for-30-minutes, for-1-hour, indefinitely or schedule.
* `start_date` - (Required) This parameter takes a date format as (yyyy-MM-dd'T'HH:mm:ssZ) (e.g. 2019-06-11T08:00:00+02:00).
* `end_date` - (Required) This parameter takes a date format as (yyyy-MM-dd'T'HH:mm:ssZ) (e.g. 2019-06-11T08:00:00+02:00).


`rules` supports the following:

* `entity` - (Required) This field represents the entity that maintenance will be applied. Entity field takes two mandatory fields as id and type.
  - `id` - (Required) The id of the entity that maintenance will be applied.
  - `type` - (Required) The type of the entity that maintenance will be applied. It can be either integration or policy.

* `state` - (Required) State of rule that will be defined in maintenance and can take either enabled or disabled for policy type rules. This field has to be disabled for integration type entity rules.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Maintenance Policy.

## Import

Maintenance policies can be imported using the `policy_id`, e.g.

`$ terraform import opsgenie_maintenance.test policy_id`
