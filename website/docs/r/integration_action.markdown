---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_integration_action"
sidebar_current: "docs-opsgenie-resource-integration-action"
description: |-
  Manages advanced actions for integrations within Opsgenie
---

# opsgenie_integration_action

Manages an API Integration within Opsgenie. The action types that are supported are:
* create
* close
* acknowledge
* add_note

## Example Usage

```hcl
resource "opsgenie_integration_action" "test_action" {
  integration_id = opsgenie_api_integration.test.id


  create {
    name = "create action"
    user = "{{user}}"
    note = "{{note}}"
    alias = "{{alias}}"
    source = "{{source]}"
    message = "{{message}}"
    description = "{{description}}"
    entity = "{{entity}}"

    responders {
      type = "escalation"
      id = "${opsgenie_escalation.test.id}"
    }

    filter {
      type = "match-all"
      conditions {
        field = "message"
        operation = "contains"
        expected_value = "ERROR"
      }
    }
  }

  close {
    name = "close action"
    filter {
      type = "match-any-condition"
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P4"
      }
      conditions {
        field = "priority"
        operation = "equals"
        expected_value = "P5"
      }
    }
  }

}
```

## Argument Reference

The following arguments are common and supported for all actions:

* `integration_id` - (Required) ID of the parent integration resource to bind to.

* `name` - (Required) Name of the integration action.

* `alias` - (Optional) An identifier that is used for alert deduplication. Defaults to `{{alias}}`.

* `order` - (Optional) Integer value that defines in which order the action will be performed. Defaults to `1`.

* `user` - (Optional) Owner of the execution for integration action.

* `note` - (Optional) Integer value that defines in which order the action will be performed.

* `filter` - (Optional) Used to specify rules for matching alerts and the filter type.

### Additional Arguments for Create Action

* `description` - (Optional)  Detailed description of the alert, anything that may not have fit in the `message` field.

* `entity` - (Optional) The entity that the alert is related to.

* `extra_properties` - (Optional) Set of user defined properties specified as a map.

* `message` - (Optional) Alert text limited to 130 characters.

* `responders` - (Optional) User, schedule, teams or escalation names to calculate which users will receive notifications of the alert.

* `source` - (Optional) User defined field to specify source of action.

* `tags` - (Optional) Comma separated list of labels to be attached to the alert.

* `ignore_responders_from_payload` - (Optional) If enabled, the integration will ignore responders sent in request payloads.

* `ignore_teams_from_payload` - (Optional) If enabled, the integration will ignore teams sent in request payloads.

`responders` is supported only in create action and supports the following:

* `type` - (Required) The responder type.
* `id` - (Required) The id of the responder.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie API Integration.

## Import

API Integrations can be imported using the `id`, e.g.

`$ terraform import opsgenie_integration_action.defaultintegration 812be1a1-32c8-4666-a7fb-03ecc385106c`
