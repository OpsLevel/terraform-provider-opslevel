---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_team Data Source - terraform-provider-opslevel"
subcategory: ""
description: |-
  
---

# opslevel_team (Data Source)



## Example Usage

```terraform
data "opslevel_team" "devs" {
  alias = "developers"
}

data "opslevel_team" "devs" {
  id = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS83NzQ0"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `alias` (String) An alias of the team to find by.
- `id` (String) The id of the team to find.

### Read-Only

- `members` (List of Object) List of members in the team with email address and role. (see [below for nested schema](#nestedatt--members))
- `name` (String) The name of the team.
- `parent_alias` (String) The alias of the parent team.
- `parent_id` (String) The id of the parent team.

<a id="nestedatt--members"></a>
### Nested Schema for `members`

Read-Only:

- `email` (String)
- `role` (String)


