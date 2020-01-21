---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_user_contact"
sidebar_current: "docs-opsgenie-resource-user-contact"
description: |-
  Manages a User Contact.
---

# opsgenie_user_contact

Manages a User Contact.

## Example Usage

```hcl
resource "opsgenie_user_contact" "sms" {
  username = "${opsgenie_user.exampleuser.username}"
  to       = "39-123"
  method   = "sms"
}

resource "opsgenie_user_contact" "email" {
  username = "${opsgenie_user.exampleuser.username}"
  to       = "fahri@opsgenie.com"
  method   = "email"
}

resource "opsgenie_user_contact" "voice" {
  username = "${opsgenie_user.exampleuser.username}"
  to       = "39-123"
  method   = "voice"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The username for contact.(reference)

* `to` - (Required) to field is the address of given method.
                    
* `method` - (Required) This parameter is the contact method of user and should be one of email, sms or voice. Please note that adding mobile is not supported from API.

* `enabled` - (Optional) Enable contact of the user in OpsGenie. Default value is true.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Contact.

## Import

Users can be imported using the `id`, e.g.

`$ terraform import opsgenie_user_contact.testcontact username/contactId`

For this example:
- Username = `genie@awesometeam.com` 
- Contact Id = `2d1a78d0-c13e-47d3-af0a-8b6d0cc2b7b1`

`$ terraform import opsgenie_user_contact.testcontact genie@awesometeam.com/2d1a78d0-c13e-47d3-af0a-8b6d0cc2b7b1`
