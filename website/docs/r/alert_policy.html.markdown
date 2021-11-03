---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_alert_policy"
sidebar_current: "docs-opsgenie-resource-alert-policy"
description: |-
  Manages a Alert Policy within Opsgenie.
---

# opsgenie\_alert\_policy

Manages a Alert Policy within Opsgenie.

## Example Usage

```hcl

resource "opsgenie_team" "test" {
  name        = "example team"
  description = "This team deals with all the things"

}
resource "opsgenie_alert_policy" "test" {
  name               = "example policy"
  team_id            = opsgenie_team.test.id
  policy_description = "This is sample policy"
  message            = "{{message}}"

  filter {}
  time_restriction {
    type = "weekday-and-time-of-day"
    restrictions {
      end_day    = "monday"
      end_hour   = 7
      end_min    = 0
      start_day  = "sunday"
      start_hour = 21
      start_min  = 0
    }
    restrictions {
      end_day    = "tuesday"
      end_hour   = 7
      end_min    = 0
      start_day  = "monday"
      start_hour = 22
      start_min  = 0
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the alert policy

* `team_id` - (Optional) Id of team that this policy belongs to.

* `enabled` - (Optional) If policy should be enabled. Default: `true`

* `policy_description` - (Optional) Description of the policy. This can be max 512 characters.

* `filter` - (Required) A alert filter which will be applied. This filter can be empty: `filter {}` - this means `match-all`. This is a block, structure is documented below.

* `time_restriction` - (Optional) Time restrictions specified in this field must be met for this policy to work. This is a block, structure is documented below.

* `message` - (Required) Message of the alerts

* `continue_policy` - (Optional) It will trigger other modify policies if set to `true`. Default: `false`

* `alias` - (Optional) Alias of the alert. You can use `{{alias}}` to refer to the original alias. Default: `{{alias}}`

* `alert_description` - (Optional) Description of the alert. You can use `{{description}}` to refer to the original alert description. Default: `{{description}}`

* `entity` - (Optional) Entity field of the alert. You can use `{{entity}}` to refer to the original entity. Default: `{{entity}}`

* `source` - (Optional) Source field of the alert. You can use `{{source}}` to refer to the original source. Default: `{{source}}`

* `ignore_original_actions` - (Optional) If set to `true`, policy will ignore the original actions of the alert. Default: `false`

* `ignore_original_details` - (Optional) If set to `true`, policy will ignore the original details of the alert. Default: `false`

* `ignore_original_responders` - (Optional) If set to `true`, policy will ignore the original responders of the alert. Default: `false`

* `responders` - (Optional) Responders to add to the alerts original responders value as a list of teams, users or the reserved word none or all. If `ignore_original_responders` field is set to `true`, this will replace the original responders. The possible values for responders are: `user`, `team`. This is a block, structure is documented below.

* `ignore_original_tags` - (Optional) If set to `true`, policy will ignore the original tags of the alert. Default: `false`

* `actions` - (Optional) Actions to add to the alerts original actions value as a list of strings. If `ignore_original_actions` field is set to `true`, this will replace the original actions.

* `tags` - (Optional) Tags to add to the alerts original tags value as a list of strings. If `ignore_original_responders` field is set to `true`, this will replace the original responders.

* `priority` - (Optional) Priority of the alert. Should be one of `P1`, `P2`, `P3`, `P4`, or `P5`



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

The `responders` block supports:

* `type` - (Required) Type of responder. Acceptable values are: `user` or `team`

* `name` - (Optional) Name of the responder

* `id` - (Required) ID of the responder

* `username` - (Optional) Username of the responder

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Alert Policy.

## Import

Alert policies can be imported using the `team_id/policy_id`, e.g.

`$ terraform import opsgenie_notification_policy.test team_id/policy_id`

You can import global polices using only policy identifier

`$ terraform import opsgenie_alert_policy.test policy_id`
