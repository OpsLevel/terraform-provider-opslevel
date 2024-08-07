---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_service_tag Resource - terraform-provider-opslevel"
subcategory: ""
description: |-
  Service Tag Resource
---

# opslevel_service_tag (Resource)

Service Tag Resource

## Example Usage

```terraform
resource "opslevel_service_tag" "service_tag_1" {
  key     = "hello"
  value   = "world"
  service = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MzMyOQ"
}

resource "opslevel_service_tag" "service_tag_2" {
  key           = "hello_with_alias"
  value         = "world_with_alias"
  service_alias = "cart"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) The tag's key.
- `value` (String) The tag's value.

### Optional

- `service` (String) The id of the service that this will be added to.
- `service_alias` (String) The alias of the service that this will be added to.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import opslevel_service_tag.example Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ:Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8xODA
```
