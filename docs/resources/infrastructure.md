---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_infrastructure Resource - terraform-provider-opslevel"
subcategory: ""
description: |-
  Manages an infrastructure resource
---

# opslevel_infrastructure (Resource)

Manages an infrastructure resource

## Example Usage

```terraform
data "opslevel_group" "foo" {
    identifier = "foo"
}

// Minimum example
resource "opslevel_infrastructure" "example_1" {
    schema = "Database"
    provider_data {
        account = "dev"
    }
    data = jsonencode({
        name = "my-database"
    })
}

// Detailed example
resource "opslevel_infrastructure" "example_2" {
    schema = "Database"
    owner = data.opslevel_group.devs.id
    provider_data {
        account = "dev"
        name = "google cloud"
        type = "BigQuery"
        url = "https://console.cloud.google.com/..."
    }
    data = jsonencode({
        name = "big-query"
        zone = "us-east-1"
        engine = "bigquery"
        engine_version = "1.28.0"
        endpoint = "https://console.cloud.google.com/..."
        replica = false
        publicly_accessible = false
        storage_size = {
            unit = "GB"
            value = 700
        }
        storage_type = "gp3"
        storage_iops = {
            unit = "per second"
            value = 12000
        }
    })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `schema` (String) The schema of the infrastructure resource that determines its data specification.

### Optional

- `data` (String) The data of the infrastructure resource in JSON format.
- `last_updated` (String)
- `owner` (String) The id of the owner for the infrastructure resource.  Can be a team or group. Does not support aliases!
- `provider_data` (Block List, Max: 1) The provider specific data for the infrastructure resource. (see [below for nested schema](#nestedblock--provider_data))

### Read-Only

- `aliases` (List of String) The aliases of the infrastructure resource.
- `id` (String) The ID of this resource.

<a id="nestedblock--provider_data"></a>
### Nested Schema for `provider_data`

Required:

- `account` (String) The canonical account name for the provider of the infrastructure resource.

Optional:

- `name` (String) The name of the provider of the infrastructure resource. (eg. AWS, GCP, Azure)
- `type` (String) The type of the infrastructure resource as defined by its provider.
- `url` (String) The url for the provider of the infrastructure resource.

