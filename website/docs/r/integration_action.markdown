---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_integration_action"
sidebar_current: "docs-opsgenie-resource-integration-action"
description: |-
  Manages advanced actions for integrations within Opsgenie
---

# opsgenie_integration_action

Manages advanced actions for Integrations within Opsgenie. This applies for the following resources:
* [`opsgenie_api_integration`](api_integration.html)
* [`opsgenie_email_integration`](email_integration.html)

The actions that are supported are:
* `create`
* `close`
* `acknowledge`
* `add_note`
* `ignore`

## Example Usage

```hcl
resource "opsgenie_integration_action" "test_action" {
  integration_id = opsgenie_api_integration.test.id


  create {
    name = "create action"
    tags = ["CRITICAL", "SEV-0"]
    user = "Example-service"
    note = "{{note}}"
	alias = "{{alias}}"
	source = "{{source}}"
	message = "{{message}}"
	description = "{{description}}"
	entity = "{{entity}}"
	alert_actions = ["Runbook ID#342"]
    
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P1"
      }
    }
    responders {
      id = "${opsgenie_team.test.id}"
      type = "team"
    }
  }

  create {
    name = "Create medium priority alerts"
    tags = ["SEVERE", "SEV-1"]
    priority = "P3"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P2"
      }
    }
  }
  
  create {
    name = "Create alert with priority from message"
    custom_priority = "{{message.substringAfter(\"[custom]\")}}"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "tags"
        operation = "contains"
        expected_value = "P5"
      }
      conditions {
        field = "message"
        operation = "Starts With"
        expected_value = "[custom]"
      }
    }
  }

  close {
    name = "Low priority alerts"
    filter {
      type = "match-any-condition"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
      conditions {
        field = "message"
        operation = "contains"
        expected_value = "DEBUG"
      }
    }
  }

  acknowledge {
    name = "Auto-ack test alerts"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "message"
        operation = "contains"
        expected_value = "TEST"
      }
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
    }
  }

  add_note {
    name = "Add note to all alerts"
    note = "Created from test integration"
    filter {
      type = "match-all"
    }
  }
  
  ignore {
    name = "Ignore alerts with ignore tag"
    filter {
      type = "match-all-conditions"
      conditions {
        field = "tags"
        operation = "contains"
        expected_value = "ignore"
      }
    }
  }
}
```

## Argument Reference

The following arguments are common and supported for all actions:

* `integration_id` - (Required) ID of the parent integration resource to bind to.

* `name` - (Required) Name of the integration action.

* `alias` - (Optional) An identifier that is used for alert deduplication. Default: `{{alias}}`.

* `order` - (Optional) Integer value that defines in which order the action will be performed. Default: `1`.

* `user` - (Optional) Owner of the execution for integration action.

* `note` - (Optional) Additional alert action note.

* `filter` - (Optional) Used to specify rules for matching alerts and the filter type. Please note that depending on the integration type the field names in the filter conditions are:
  * For SNS integration: `actions`, `alias`, `entity`, `Message`, `recipients`, `responders`, `Subject`, `tags`, `teams`, `eventType`, `Timestamp`, `TopicArn`.
  * For API integration: `message`, `alias`, `description`, `source`, `entity`, `tags`, `actions`, `details`, `extra-properties`, `recipients`, `teams`, `priority`, `eventType`.
  * For Email integration: `from_address`, `from_name`, `conversationSubject`, `subject`

### Additional Arguments for Create Action

* `description` - (Optional)  Detailed description of the alert, anything that may not have fit in the `message` field.

* `entity` - (Optional) The entity the alert is related to.

* `priority` - (Optional) Alert priority.

* `custom_priority` - (Optional) Custom alert priority. e.g. ``{{message.substring(0,2)}}``

* `extra_properties` - (Optional) Set of user defined properties specified as a map.

* `message` - (Optional) Alert text limited to 130 characters.

* `responders` - (Optional) User, schedule, teams or escalation names to calculate which users will receive notifications of the alert.

* `source` - (Optional) User defined field to specify source of action.

* `tags` - (Optional) Comma separated list of labels to be attached to the alert.

* `ignore_responders_from_payload` - (Optional) If enabled, the integration will ignore responders sent in request payloads.

* `ignore_teams_from_payload` - (Optional) If enabled, the integration will ignore teams sent in request payloads.

`responders` is supported only in create action and supports the following:

* `id` - (Required) The id of the responder.
* `type` - (Required) The responder type - can be `escalation`, `team` or `user`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie API Integration.
