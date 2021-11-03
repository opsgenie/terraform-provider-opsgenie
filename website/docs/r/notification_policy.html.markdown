---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_notification_policy"
sidebar_current: "docs-opsgenie-resource-notification-policy"
description: |-
  Manages a Notification Policy within Opsgenie.
---

# opsgenie\_notification\_policy

Manages a Notification Policy within Opsgenie.

## Example Usage

```hcl

resource "opsgenie_team" "test" {
  name        = "example team"
  description = "This team deals with all the things"

}
resource "opsgenie_notification_policy" "test" {
  name               = "example policy"
  team_id            = opsgenie_team.test.id
  policy_description = "This policy has a delay action"
  delay_action {
    delay_option = "next-time"
    until_minute = 1
    until_hour   = 9
  }
  filter {}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the notification policy

* `team_id` - (Required) Id of team that this policy belons to.

* `enabled` - (Optional) If policy should be enabled. Default: `true`

* `policy_description` - (Optional) Description of the policy. This can be max 512 characters.

* `filter` - (Required) A notification filter which will be applied. This filter can be empty: `filter {}` - this means `match-all`. This is a block, structure is documented below.

* `time_restriction` - (Optional) Time restrictions specified in this field must be met for this policy to work. This is a block, structure is documented below.

* `auto_close_action` - (Optional) Auto Restart Action of the policy. This is a block, structure is documented below.

* `auto_restart_action` - (Optional) Auto Restart Action of the policy. This is a block, structure is documented below.

* `de_duplication_action` - (Optional) Deduplication Action of the policy. This is a block, structure is documented below.

* `delay_action` - (Optional) Delay notifications. This is a block, structure is documented below.

* `suppress` - (Optional) Suppress value of the policy. Values are: `true`, `false`. Default: `false`


The `filter` block supports:

* `type` (Optional) - A filter type, supported types are: `match-all`, `match-any-condition`, `match-all-conditions`. Default: `match-all`

* `conditions` (Optional) Conditions applied to filter. This is a block, structure is documented below.

The `conditions` block supports:

* `field` - (Required) Specifies which alert field will be used in condition. Possible values are `message`, `alias`, `description`, `source`, `entity`, `tags`, `actions`, `details`, `extra-properties`, `recipients`, `teams`, `priority`

* `operation` - (Required) It is the operation that will be executed for the given field and key. Possible operations are `matches`, `contains`, `starts-with`, `ends-with`, `equals`, `contains-key`, `contains-value`, `greater-than`, `less-than`, `is-empty`, `equals-ignore-whitespace`.

* `key` - (Optional) If `field` is set as extra-properties, key could be used for key-value pair

* `not` - (Optional) Indicates behaviour of the given operation. Default: `false`

* `expected_value` - (Optional) User defined value that will be compared with alert field according to the operation. Default: empty string

* `order` - (Optional) Order of the condition in conditions list

The `time_restriction` block supports:

* `type` - (Required) Defines if restriction should apply daily on given hours or on certain days and hours. Possible values are: `time-of-day`, `weekday-and-time-of-day`

* `restrictions` - (Optional) List of days and hours definitions for field type = `weekday-and-time-of-day`. This is a block, structure is documented below.

* `restriction` - (Optional) A definition of hourly definition applied daily, this has to be used with combination: type = `time-of-day`. This is a block, structure is documented below.

The `restrictions` block supports:

* `start_day` - (Required) Starting day of restriction (eg. `monday`)

* `end_day` - (Required) Ending day of restriction (eg. `wednesday`)

* `start_hour` - (Required) Starting hour of restriction on defined `start_day`

* `end_hour` - (Required) Ending hour of restriction on defined `end_day`

* `start_min` - (Required) Staring minute of restriction on defined `start_hour`

* `end_min` - (Required) Ending minute of restriction on defined `end_hour`

The `restriction` block supports:

* `start_hour` - (Required) Starting hour of restriction.

* `end_hour` - (Required) Ending hour of restriction.

* `start_min` - (Required) Staring minute of restriction on defined `start_hour`

* `end_min` - (Required) Ending minute of restriction on defined `end_hour`

The `auto_close_action` block supports:

* `duration` - (Required) Duration of this action. This is a block, structure is documented below.

The `auto_restart_action` block supports:

* `duration` - (Required) Duration of this action. This is a block, structure is documented below.

* `max_repeat_count` - (Required) How many times to repeat. This is a integer attribute.

The `de_duplication_action` block supports:

* `de_duplication_action_type` - (Required) Deduplication type. Possible values are: "value-based", "frequency-based"

* `count` - (Required) - Count

* `duration` - (Optional) Duration of this action (only required for "frequency-based" de-duplication action). This is a block, structure is documented below.

The `delay_action` block supports:

* `delay_option` - (Required) Defines until what day to delay or for what duration. Possible values are: `for-duration`, `next-time`, `next-weekday`, `next-monday`, `next-tuesday`, `next-wednesday`, `next-thursday`, `next-friday`, `next-saturday`, `next-sunday`

* `duration` - (Optional) Duration of this action. If `delay_option` = `for-duration` this has to be set. This is a block, structure is documented below.

* `until_hour` - (Optional) Until what hour notifications will be delayed. If `delay_option` is set to antyhing else then `for-duration` this has to be set.

* `until_minute` - (Optional) Until what minute on `until_hour` notifications will be delayed. If `delay_option` is set to antyhing else then `for-duration` this has to be set.

The `duration` block supports:

* `time_unit` - (Optional) Valid time units are: `minutes`, `hours`, `days`. Default: `minutes`

* `time_amount` - (Required) A amount of time in `time_units`. This is a integer attribute.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Notification Policy.

## Import

Notification policies can be imported using the `team_id` and `notification_policy_id`, e.g.

`$ terraform import opsgenie_notification_policy.test team_id/notification_policy_id`
