---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_service_incident_rule"
sidebar_current: "docs-opsgenie-resource-service-incident-rule"
description: |-
  Manages a Service Incident Rule within Opsgenie.
---

# opsgenie\_service\_incident\_rule

Manages a Service Incident Rule within Opsgenie.

## Example Usage

```hcl

resource "opsgenie_team" "test" {
  name        = "example-team"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test" {
  name  = "example-service"
  team_id = opsgenie_team.test.id
}
resource "opsgenie_service_incident_rule" "test" {
  service_id = opsgenie_service.test.id
  incident_rule {
	condition_match_type = "match-any-condition"
	conditions {
		field = "message"
		not =  false
		operation = "contains"
		expected_value = "expected1"
	}
	conditions {
		field = "message"
		not =  false
		operation = "contains"
		expected_value = "expected2"
	}
	incident_properties {
		message = "This is a test message"
		priority = "P3"
		stakeholder_properties {
			message = "Message for stakeholders"
			enable = "true"
		}
	}
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) ID of the service associated

* `incident_rule` - (Required) This is the rule configuration for this incident rule. This is a block, structure is documented below.

The `incident_rule` block supports:

* `condition_match_type` - (Optional) A Condition type, supported types are: `match-all`, `match-any-condition`, `match-all-conditions`. Default: `match-all`

* `conditions` - (Optional) Conditions applied to incident. This is a block, structure is documented below.

* `incident_properties`- (Required) Properties for incident rule. This is a block, structure is documented below.


The `conditions` block supports:

* `field` - (Required) Specifies which alert field will be used in condition. Possible values are `message`, `alias`, `description`, `source`, `entity`, `tags`, `actions`, `details`, `extra-properties`, `recipients`, `teams`, `priority`

* `operation` - (Required) It is the operation that will be executed for the given field and key. Possible operations are `matches`, `contains`, `starts-with`, `ends-with`, `equals`, `contains-key`, `contains-value`, `greater-than`, `less-than`, `is-empty`, `equals-ignore-whitespace`.

* `not` - (Optional) Indicates behaviour of the given operation. Default: false

* `expected_value` - (Optional) User defined value that will be compared with alert field according to the operation. Default: empty string


The `incident_properties` block supports:

* `message` - (Required) Message of the related incident rule.

* `tags` - (Optional) Tags of the alert.

* `details` - (Optional) Map of key-value pairs to use as custom properties of the alert.

* `description` - (Optional) Description field of the incident rule.

* `priority` - (Required) Priority level of the alert. Possible values are `P1`, `P2`, `P3`, `P4` and `P5`

* `stakeholder_properties` - (Required) DEtails about stakeholders for this rule. This is a block, structure is documented below.


The `stakeholder_properties` block supports:

* `enable` - (Optional) Option to enable stakeholder notifications.Default value is true.

* `message` - (Required) Message that is to be passed to audience that is generally used to provide a content information about the alert.

* `description` - (Optional) Description that is generally used to provide a detailed information about the alert.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Service Incident Policy.

## Import

Service Incident Rule can be imported using the `service_id/service_incident_rule_id`, e.g.

`$ terraform import opsgenie_service_incident_rule.this service_id/service_incident_rule_id`
