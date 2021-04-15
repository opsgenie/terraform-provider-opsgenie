---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_schedule_rotation"
sidebar_current: "docs-opsgenie-resource-schedule-rotation"
description: |-
  Manages a Schedule Rotation within Opsgenie.
---

# opsgenie_schedule_rotation

Manages a Schedule Rotation within Opsgenie.

## Example Usage
```hcl
resource "opsgenie_schedule_rotation" "test" {
  schedule_id = "${opsgenie_schedule.test.id}"
  name        = "test"
  start_date  = "2019-06-18T17:00:00Z"
  end_date    = "2019-06-20T17:30:00Z"
  type        = "hourly"
  length      = 6

  participant {
    type = "user"
    id   = "${opsgenie_user.test.id}"
  }

  time_restriction {
    type = "time-of-day"

    restriction {
      start_hour = 1
      start_min  = 1
      end_hour   = 10
      end_min    = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `schedule_id` - (Required) Identifier of the schedule.

* `name` - (Optional) Name of rotation.

* `start_date` - (Required) This parameter takes a date format as (yyyy-MM-dd'T'HH:mm:ssZ) (e.g. 2019-06-11T08:00:00+02:00). Minutes may take 0 or 30 as value. Otherwise they will be converted to nearest 0 or 30 automatically

* `end_date` - (Optional)  This parameter takes a date format as (yyyy-MM-dd'T'HH:mm:ssZ) (e.g. 2019-06-11T08:00:00+02:00). Minutes may take 0 or 30 as value. Otherwise they will be converted to nearest 0 or 30 automatically

* `type` - (Required) Type of rotation. May be one of daily, weekly and hourly.

* `length` - (Optional) Length of the rotation with default value 1.

* `participant` - (Required) List of escalations, teams, users or the reserved word none which will be used in schedule. Each of them can be used multiple times and will be rotated in the order they given. "user,escalation,team,none"

* `time_restriction` - (Optional)

`participant` supports the following:

* `type` - (Required) The responder type.
* `id` - (Required) The id of the responder.

`time_restriction` supports the following:

* `type` - (Required) This parameter should be set to `time-of-day` or `weekday-and-time-of-day`.

* `restriction` - (Optional) It is a restriction object which is described below. In this case startDay/endDay fields are not supported. This can be used only if time restriction type is `time-of-day`.

    `restriction` supports the following:

     * `start_hour` - (Required) Value of the hour that frame will start.
     * `start_min` - (Required) Value of the minute that frame will start. Minutes may take 0 or 30 as value. Otherwise they will be converted to nearest 0 or 30 automatically.
     * `end_hour` - (Required) Value of the hour that frame will end.
     * `end_min` - (Required) Value of the minute that frame will end. Minutes may take 0 or 30 as value. Otherwise they will be converted to nearest 0 or 30 automatically.

* `restrictions` - (Optional) It is a restriction object which is described below. This can be used only if time restriction type is `weekday-and-time-of-day`.

    `restrictions` supports the following:

     * `start_day` - (Required) Value of the day that frame will start.
     * `start_hour` - (Required) Value of the hour that frame will start
     * `start_min` - (Required) Value of the minute that frame will start. Minutes may take 0 or 30 as value. Otherwise they will be converted to nearest 0 or 30 automatically.
     * `end_day` - (Required) Value of the day that frame will end.
     * `end_hour` - (Required) Value of the hour that frame will end.
     * `end_min` - (Required) Value of the minute that frame will end. Minutes may take 0 or 30 as value. Otherwise they will be converted to nearest 0 or 30 automatically.

     Both `start_day` and `end_day` can assume only `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, or `sunday` values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Schedule Rotation

## Import

Schedule Rotations can be imported using the `schedule_id/rotation_id`, e.g.

* `terraform import opsgenie_schedule_rotation.test schedule_id/rotation_id`
