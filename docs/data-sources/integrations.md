---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_integrations Data Source - terraform-provider-opslevel"
subcategory: ""
description: |-
  
---

# opslevel_integrations (Data Source)



## Example Usage

```terraform
data "opslevel_integrations" "all" {
}

output "found" {
  value = data.opslevel_integrations.all.id[0]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of String)
- `names` (List of String)


