---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_team"
sidebar_current: "docs-opsgenie-datasource-team"
description: |-
  Manages existing Team within Opsgenie.
---

# opsgenie\_team

Manages existing Team within Opsgenie.

## Example Usage

```hcl
data "opsgenie_team" "sre-team" {
  name = "sre-team"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name associated with this team. Opsgenie defines that this must not be longer than 100 characters.

The following attributes are exported:

* `id` - The ID of the Opsgenie Team.

* `member` - A Member block as documented below.

* `description` - A description for this team.
