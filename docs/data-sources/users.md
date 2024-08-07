---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_users Data Source - terraform-provider-opslevel"
subcategory: ""
description: |-
  List of all User data sources
---

# opslevel_users (Data Source)

List of all User data sources

## Example Usage

```terraform
data "opslevel_users" "all" {}

output "all" {
  value = data.opslevel_users.all.users
}

output "user_emails" {
  value = sort(data.opslevel_users.all.users[*].email)
}

output "user_names" {
  value = sort(data.opslevel_users.all.users[*].name)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `users` (Attributes List) List of user data sources (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `email` (String) The email of the user.
- `id` (String) The unique identifier for the user.
- `name` (String) The name of the user.
- `role` (String) The user's assigned role.


