---
layout: "opsgenie"
page_title: "Provider: Opsgenie"
sidebar_current: "docs-opsgenie-index"
description: |-
  The Opsgenie provider is used to interact with the many resources supported by Opsgenie. The provider needs to be configured with the proper credentials before it can be used.
---

# Opsgenie Provider

The Opsgenie provider is used to interact with the
many resources supported by Opsgenie. The provider needs to be configured
with the proper credentials before it can be used.

**Breaking Change - v0.6.0**

With 0.6.0 version provider adopted Terraform Plugin SDK v2 therefore some resources reads has changed. 
If you encounter any problems you can contact us via Github

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Opsgenie Provider
provider "opsgenie" {
  api_key = "key"
  api_url = "api.eu.opsgenie.com" #default is api.opsgenie.com
}

# Create a user
resource "opsgenie_user" "test" {
  # ...
}
```

## Configuration Reference

The following arguments are supported:

* `api_key` - (Required) The API Key for the Opsgenie Integration. If omitted, the
  `OPSGENIE_API_KEY` environment variable is used.

* `api_url` - (Optional) The API url for the Opsgenie.

You can generate an API Key within Opsgenie by creating a new API Integration with Read/Write permissions.

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OPSGENIE_API_KEY` - The API Key used for the Opsgenie Integration.
