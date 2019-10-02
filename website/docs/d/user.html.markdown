---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_user"
sidebar_current: "docs-opsgenie-resource-user"
description: |-
  Manages existing User within Opsgenie.
---

# opsgenie_user

Manages existing User within Opsgenie.

## Example Usage

```hcl
data "opsgenie_user" "test" {
  username = "user@domain.com"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The email address associated with this user. Opsgenie defines that this must not be longer than 100 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie User.

* `full_name` - The Full Name of the User.

* `role` - The Role assigned to the User. Either a built-in such as 'Owner', 'Admin' or 'User' - or the name of a custom role.

* `locale` - Location information for the user. Please look at [Supported Locale Ids](https://docs.opsgenie.com/docs/supported-locales) for available locales.

* `timezone` - Timezone information of the user. Please look at [Supported Timezone Ids](https://docs.opsgenie.com/docs/supported-timezone-ids) for available timezones.
