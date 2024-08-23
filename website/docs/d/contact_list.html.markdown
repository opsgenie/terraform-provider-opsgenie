---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_contact_list"
sidebar_current: "docs-opsgenie-resource-contact-list"
description: |-
  Manages existing Contact within Opsgenie.
---

# opsgenie_user

Manages existing User's contacts within Opsgenie.

## Example Usage

```hcl
data "opsgenie_contact_list" "test" {
  username = "user@domain.com"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The email address associated with this user. Opsgenie defines that this must not be longer than 100 characters.

## Attributes Reference

The following attributes are exported:

* `contact_list` -  A list of Contact blocks as documented below.

Contact block contains:

* `id` - Id of user's contact.

* `method` - Contact method of user (should be one of email, sms or voice).

* `to` - Address of method.
