---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_user"
sidebar_current: "docs-opsgenie-resource-user"
description: |-
  Manages a User within Opsgenie.
---

# opsgenie_user

Manages a User within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_user" "test" {
  username  = "user@domain.com"
  full_name = "Test User"
  role      = "User"
  locale    = "en_US"
  timezone  = "America/New_York"
  tags = ["sre", "opsgenie"]
  skype_username = "skypename"
  user_address {
      country = "Country"
      state = "State"
      city = "City"
      line = "Line"
      zipcode = "998877"
  }
  user_details = {
    key1 = "val1,val2"
    key2 = "val3,val4"
  }
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The email address associated with this user. Opsgenie defines that this must not be longer than 100 characters and must contain lowercase characters only.

* `full_name` - (Required) The Full Name of the User.

* `role` - (Required) The Role assigned to the User. Either a built-in such as 'Admin' or 'User' - or the name of a custom role.

* `locale` - (Optional) Location information for the user. Please look at [Supported Locale Ids](https://docs.opsgenie.com/docs/supported-locales) for available locales.

* `timezone` - (Optional) Timezone information of the user. Please look at [Supported Timezone Ids](https://docs.opsgenie.com/docs/supported-timezone-ids) for available timezones.

* `tags` - (Optional) A list of tags to be associated with the user.

* `skype_username` - (Optional) Skype username of the user.

* `user_details` - (Optional) Details about the user in form of key and list. of values.

* `user_address` - (Optional) Address of the user.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie User.

## Import

Users can be imported using the `user_id`, e.g.

`$ terraform import opsgenie_user.user user_id`
