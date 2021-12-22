---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_incident_template"
sidebar_current: "docs-opsgenie-resource-incident-template"
description: |-
  Manages an Incident Template within Opsgenie.
---

# opsgenie\_incident_template

Manages an Incident Template within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_team" "test" {
  name        = "genietest-team"
  description = "This team deals with all the things"
}
resource "opsgenie_service" "test" {
  name  = "genietest-service"
  team_id = opsgenie_team.test.id
}
resource "opsgenie_incident_template" "test" {
  name = "genietest-incident-template"
  message = "Incident Message"
  priority = "P2"
  stakeholder_properties {
    enable = true
    message = "Stakeholder Message"
    description = "Stakeholder Description"
  }
  tags = ["tag1", "tag2"]
  description = "Incident Description"
  details = {
    key1 = "value1"
    key2 = "value2"
  }
  impacted_services = [
    opsgenie_service.test.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) Name of the incident template.

* `message` (Required) Message of the related incident template. This field must not be longer than 130 characters.

* `description` (Optional) Description field of the incident template. This field must not be longer than 10000 characters.

* `tags` (Optional) Tags of the incident template.

* `details` (Optional) Map of key-value pairs to use as custom properties of the incident template. This field must not be longer than 8000 characters.

* `priority` (Required) Priority level of the incident. Possible values are `P1`, `P2`, `P3`, `P4` and `P5`.

* `impactedServices` (Optional) Impacted services of incident template. Maximum 20 services.

* `stakeholderProperties` (Required)

   * `enable` (Optional) Option to enable stakeholder notifications.Default value is true.

   * `message` (Required) Message that is to be passed to audience that is generally used to provide a content information about the alert.

   * `description` (Optional) Description that is generally used to provide a detailed information about the alert. This field must not be longer than 15000 characters.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Incident Template.

## Import

Service can be imported using the `template_id`, e.g.

`$ terraform import opsgenie_incident_template.test template_id`
