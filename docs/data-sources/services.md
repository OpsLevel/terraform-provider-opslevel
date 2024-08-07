---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_services Data Source - terraform-provider-opslevel"
subcategory: ""
description: |-
  Services data source
---

# opslevel_services (Data Source)

Services data source

## Example Usage

```terraform
data "opslevel_services" "all" {}

data "opslevel_tier" "tier1" {
  filter {
    field = "alias"
    value = "tier_1"
  }
}

data "opslevel_services" "tier1" {
  filter = {
    field = "tier"
    value = data.opslevel_tier.tier1.alias
  }
}

data "opslevel_services" "frontend" {
  filter = {
    field = "owner"
    value = "frontend"
  }
}

output "all_services" {
  value = data.opslevel_services.all.services
}

output "all_service_names" {
  value = sort(data.opslevel_services.all.services[*].name)
}

output "tier1_services" {
  value = data.opslevel_services.tier1.services
}

output "frontend_services" {
  value = data.opslevel_services.frontend.services
}


output "frontend_services_urls" {
  value = sort(data.opslevel_services.frontend.services[*].url)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Attributes) Used to filter services by one of 'filter`, `framework`, `language`, `lifecycle`, `owner`, `product`, `tag`, `tier' (see [below for nested schema](#nestedatt--filter))

### Read-Only

- `services` (Attributes List) List of Service data sources (see [below for nested schema](#nestedatt--services))

<a id="nestedatt--filter"></a>
### Nested Schema for `filter`

Required:

- `field` (String) The field of the target resource to filter upon. One of `filter`, `framework`, `language`, `lifecycle`, `owner`, `product`, `tag`, `tier`
- `value` (String) The field value of the target resource to match.


<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- `id` (String) The id of the service
- `name` (String) The display name of the service.
- `url` (String) A link to the HTML page for the resource


