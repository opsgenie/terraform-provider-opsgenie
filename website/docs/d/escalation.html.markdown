---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_escalation"
sidebar_current: "docs-opsgenie-resource-escalation"
description: |-
  You can use this data source to map existing Escalation within Opsgenie.
---

# opsgenie_escalation

Manages an Escalation within Opsgenie.

## Example Usage
```hcl
data "opsgenie_escalation" "test" {
  name = "existing-escalation"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the escalation.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Escalation.

* `rules` - Escalation rules

* `description` - Escalation Description

* `repeat`  - Escalation repeat preferences

* `owner_team_id` - If owner team exist the id of the team is exported


