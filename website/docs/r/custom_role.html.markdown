---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_custom_role"
sidebar_current: "docs-opsgenie-custom-role"
description: |-
  Manages custom user roles within Opsgenie.
---

# opsgenie_custom_role

Manages custom user roles within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_custom_role" "test" {
    role_name           = "genierole"
    extended_role       = "user"
    granted_rights      = ["alert-delete"]
    disallowed_rights   = ["profile-edit", "contacts-edit"]
}
```

## Argument Reference

The following arguments are supported:

* `role_name` - (Required) Name of the custom role.

* `extended_role` - (Required) The role from which this role has been derived. Allowed Values: "user", "observer", "stakeholder".

* `granted_rights` - (Optional) The rights granted to this role. For allowed values please refer [User Right Prerequisites](https://docs.opsgenie.com/docs/custom-user-role-api#section-user-right-prerequisites)

* `disallowed_rights` - (Optional) The rights this role cannot have. For allowed values please refer [User Right Prerequisites](https://docs.opsgenie.com/docs/custom-user-role-api#section-user-right-prerequisites)

## Attributes Reference

Only the arguments listed above are exposed as attributes.
