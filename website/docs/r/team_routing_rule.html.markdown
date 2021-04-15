---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_team_routing_rule"
sidebar_current: "docs-opsgenie-resource-team-routing-rule"
description: |-
  Manages a Team Routing Rule within Opsgenie.
---

# opsgenie\_team\_routing\_rule

Manages a Team Routing Rule within Opsgenie.

## Example Usage

```hcl

resource "opsgenie_schedule" "test" {
  name = "genieschedule"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
}

resource "opsgenie_team" "test" {
  name        = "example team"
  description = "This team deals with all the things"

}

resource "opsgenie_team_routing_rule" "test" {
  name     = "routing rule example"
  team_id  = "${opsgenie_team.test.id}"
  order    = 0
  timezone = "America/Los_Angeles"
  criteria {
    type = "match-any-condition"
    conditions {
      field          = "message"
      operation      = "contains"
      expected_value = "expected1"
      not            = false

    }
  }
  time_restriction {
    type = "weekday-and-time-of-day"
    restrictions {
      start_day  = "monday"
      start_hour = 8
      start_min  = 0
      end_day    = "tuesday"
      end_hour   = 18
      end_min    = 30
    }
  }
  notify {
    name = "${opsgenie_schedule.test.name}"
    type = "schedule"
  }

}

```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the team routing rule

* `team_id` - (Required) Id of the team owning the routing rule

* `order` - (Optional) The order of the team routing rule within the rules. order value is actually the index of the team routing rule whose minimum value is 0 and whose maximum value is n-1 (number of team routing rules is n)

* `timezone` - (Optional) Timezone of team routing rule. If timezone field is not given, account timezone is used as default.You can refer to Supported Locale IDs for available timezones

* `criteria` - (Optional) You can refer Criteria for detailed information about criteria and its fields

* `timeRestriction` - (Optional) You can refer Time Restriction for detailed information about time restriction and its fields.

* `notify` - (Required) Target entity of schedule, escalation, or the reserved word none which will be notified in routing rule. The possible values are: `schedule`, `escalation`, `none`

`notify` supports the following:

* `type` - (Required)

* `name` - (Optional)

* `id` - (Optional)


`criteria` supports the following:

* `type` - (Required) Type of the operation will be applied on conditions. Should be one of `match-all`, `match-any-condition` or `match-all-conditions`.

* `conditions` - (Optional) List of conditions will be checked before applying team routing rule. This field declaration should be omitted if the criteria type is set to match-all.


`conditions` supports the following:

* `field` - (Required) Specifies which alert field will be used in condition. Possible values are `message`, `alias`, `description`, `source`, `entity`, `tags`, `actions`, `extra-properties`, `recipients`, `teams` or `priority`.

* `key` - (Optional) If field is set as extra-properties, key could be used for key-value pair.

* `not` - (Optional) Indicates behaviour of the given operation. Default value is false.

* `operation` - (true) It is the operation that will be executed for the given field and key. Possible operations are `matches`, `contains`, `starts-with`, `ends-with`, `equals`, `contains-key`, `contains-value`, `greater-than`, `less-than`, `is-empty` and `equals-ignore-whitespace`.

* `expectedValue` - (Optional) User defined value that will be compared with alert field according to the operation. Default: empty string.

* `order` - (Optional) Order of the condition in conditions list.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Team Routing Rule.

## Import

Team Routing Rules can be imported using the `team_id/routing_rule_id`, e.g.

`$ terraform import opsgenie_team_routing_rule.ruletest team_id/routing_rule_id`
