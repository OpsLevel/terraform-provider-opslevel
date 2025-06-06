---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_property_assignment Resource - terraform-provider-opslevel"
subcategory: ""
description: |-
  Property Assignment Resource
---

# opslevel_property_assignment (Resource)

Property Assignment Resource

## Example Usage

```terraform
resource "opslevel_property_assignment" "example" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8xODA"
  owner      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value      = jsonencode("green")
}

resource "opslevel_property_assignment" "example_2" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi85OA"
  owner      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value      = jsonencode(true)
}

resource "opslevel_property_assignment" "example_3" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8xODI"
  owner      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value      = jsonencode({ "container_id" : "1c6098d6-952a-4062-9293-1dc06e991118", "container_name" : "gcr.io/containername" })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `definition` (String) The custom property definition's ID or alias.
- `owner` (String) The ID or alias of the entity that the property has been assigned to.

### Optional

- `value` (String) The value of the custom property (must be a valid JSON value or null).

### Read-Only

- `id` (String) The ID of this resource.
- `locked` (Boolean) If locked = true, the property has been set in opslevel.yml and cannot be modified in Terraform!

## Import

Import is supported using the following syntax:

```shell
terraform import opslevel_property_assignment.example Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ:Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8xODA
```
