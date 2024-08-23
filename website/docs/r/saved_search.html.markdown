---
layout: "opsgenie"
page_title: "Opsgenie: opsgenie_saved_search"
sidebar_current: "docs-opsgenie-saved-search"
description: |-
  Manages saved searches within Opsgenie.
---

# opsgenie\_saved\_search

Manages a Saved Search within Opsgenie.

## Example Usage

```hcl
resource "opsgenie_saved_search" "example" {
    name = "Example Saved Search"
    owner {
        id = data.opsgenie_user.test.id
    }
    query = "NOT priority: (\"P5\" OR \"P4\") AND status: open "
} 

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Unique name of the saved search. Maximum length is 100 characters.

* `query` - (Required) Search query to be used while filtering the alerts. Maximum length is 1000 characters.

* `owner` - (Required) User that will be assigned as owner of the saved search. Saved searches are always accessible to their owners.

* `description` - (Optional) Informational description of the saved search. Maximum length is 15000 characters.

* `teams` - (Optional) Teams that saved search is assigned to. If a saved-search is assigned to at least one team, saved-search will only be accessible to the owner and members of the assigned teams. A saved-search can be assigned to at most 20 teams.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Opsgenie Saved Search.

## Import

Teams can be imported using the `saved_search_id`, e.g.

`$ terraform import opsgenie_service.this saved_search_id`
