---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_secret Resource - terraform-provider-opslevel"
subcategory: ""
description: |-
  Manages a secret
---

# opslevel_secret (Resource)

Manages a secret

## Example Usage

```terraform
data "opslevel_team" "devs" {
  alias = "devs"
}

resource "opslevel_secret" "my_secret" {
  alias = "secret-alias"
  owner = data.opslevel_team.devs.id
  value = "too_many_passwords"
}

resource "opslevel_secret" "my_secret_2" {
  alias = "secret-alias-2"
  owner = "devs"
  value = "0sd09wer0sdlkjwer90wer098sdfsewr"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `alias` (String) The alias for this secret.
- `owner` (String) The owner of this secret.
- `value` (String, Sensitive) A sensitive value.

### Optional

- `last_updated` (String)

### Read-Only

- `created_at` (String) Timestamp of time created at.
- `id` (String) The ID of this resource.
- `updated_at` (String) Timestamp of last update.

## Import

Import is supported using the following syntax:

```shell
terraform import opslevel_secret.my_secret Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS82MDI0
```