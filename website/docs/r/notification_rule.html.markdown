---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_notification_rule"
sidebar_current: "docs-opsgenie-resource-notification-rule"
description: |-
  Manages a Notification Rule within Opsgenie.
---

# opsgenie\_notification\_rule

Manages a Notification Rule within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_user" "test" {
  username  = "Example user"
  full_name = "Name Lastname"
  role      = "User"
}

resource "opsgenie_notification_rule" "test" {
  name = "Example notification rule"
  username = opsgenie_user.test.username
  action_type = "schedule-end"
  notification_time = ["just-before", "15-minutes-ago"]
  steps {
    contact {
      method = "email"
      to = "example@user.com"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the notification policy

* `username` - (Required) Username of user to which this notification rule belongs to.

* `action_type` - (Required) Type of the action that notification rule will have. Allowed values: `create-alert`, `acknowledged-alert`, `closed-alert`, `assigned-alert`, `add-note`, `schedule-start`, `schedule-end`, `incoming-call-routing`

* `notification_time` - (Optional) List of Time Periods that notification for schedule start/end will be sent. Allowed values: `just-before`, `15-minutes-ago`, `1-hour-ago`, `1-day-ago`. If `action_type` is `schedule-start` or `schedule-end` then it is required.

* `steps` - (Optional) Notification rule steps to take (eg. SMS or email message). This is a block, structure is documented below.

* `enabled` - (Optional) If policy should be enabled. Default: `true`

The `steps` block supports:

* `enabled` - (Optional) Defined if this step is enabled. Default: `true`

* `send_after` - (Optional) Time period, in minutes, notification will be sent after.

* `contact` - (Required) Defines the contact that notification will be sent to. This is a block, structure is documented below.

The `contact` block supports:

* `method` - (Required) Contact method. Possible values: `email`, `sms`, `voice`, `mobile`

* `to` - (Required) Address of a given method (eg. email address for `email`, phone number for `sms`/`voice` or mobile application name for `mobile`)

The `criteria` block supports:

* `type` - (Required) Kind of matching filter. Possible values: `match-all`, `match-any-condition`, `match-all-conditions`

* `conditions` - (Optional) Defines the fields and values when the condition applies

The `conditions` block supports:

* `field` - (Required) Possible values: `message`, `alias`, `description`, `source`, `entity`, `tags`, `actions`, `details`, `extra-properties`, `recipients`, `teams`, `priority`

* `operation` - (Required) Possible values: `matches`, `contains`, `starts-with`, `ends-with`, `equals`, `contains-key`, `contains-value`, `greater-than`, `less-than`, `is-empty`, `equals-ignore-whitespace`

* `key` - (Optional) If 'field' is set as 'extra-properties', key could be used for key-value pair

* `not` - (Optional) Indicates behaviour of the given operation. Default: `false`

* `expected_value` - (Optional) User defined value that will be compared with alert field according to the operation. Default: empty string

* `order` - (Optional) Order of the condition in conditions list



## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Notification Rule.

## Import

Notification policies can be imported using the `user_id/notification_rule_id`, e.g.

`$ terraform import opsgenie_notification_rule.test user_id/notification_rule_id`
