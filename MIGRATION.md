# Migration Guide to v1.0.0

We at OpsLevel have been working on upgrading [our Terraform provider](https://github.com/OpsLevel/terraform-provider-opslevel) to version 1.0.0.

While the majority of these improvements are under the hood, there are a few Terraform configuration changes to be aware of.

For errors encountered during `terraform apply`, see the [State Migration](#state-migration) section below.

# BREAKING CHANGES to some resource fields

The following resources with fields that were “block type fields" will need an “=” added.

- opslevel_check_manual

```yaml
# OLD - no "=" sign
resource "opslevel_check_manual" "example" {
  update_frequency { ... }
}
# New has "=" sign
resource "opslevel_check_manual" "example" {
  update_frequency = { ... }
}
```
Just like `opslevel_check_manual.update_frequency`, the following `opslevel_resource.field`s will need an "=" added.
- `opslevel_check_alert_source_usage.alert_name_predicate`
- `opslevel_check_repository_file.file_contents_predicate`
- `opslevel_check_repository_grep.file_contents_predicate`
- `opslevel_check_repository_search.file_contents_predicate`
- `opslevel_check_service_ownership.tag_predicate`
- `opslevel_check_service_property.predicate`
- `opslevel_infrastructure.provider_data`
- `data.opslevel_repositories.filter` (datasource)
- `data.opslevel_services.filter` (datasource)


# BREAKING CHANGES with OpsLevel's updated "plural" datasources

Many of our [datasources](https://registry.terraform.io/providers/OpsLevel/opslevel/latest/docs) were updated to provide a better user experience and may require a few small updates.

## Terraform configuration changes during upgrade

When upgrading from versions older than 0.12.0-0, updating the outputs of "plural" datasources may be needed.

```terraform
# Requesting the datasource remains the same
data "opslevel_domains" "all" {}

# OLD - before v0.12.0-0
output "all" {
  value = data.opslevel_domains.all
}
output "domain_names" {
  value = data.opslevel_domains.all.names
}

# NEW - versions v0.12.0-0 and after
output "all" {
  value = data.opslevel_domains.all.domains
}
output "domain_names" {
  value = data.opslevel_domains.all.domains[*].name
}
```

This same pattern applies to `data.opslevel_scorecards`, `data.opslevel_services`, `data.opslevel_teams`, etc.

### Output Example - `data.opslevel_domains` with version >= 0.12.0

Given this configuration:

```terraform
data "opslevel_domains" "all" {}

output "all" {
  value = data.opslevel_domains.all.domains
}

output "domain_names" {
  value = sort(data.opslevel_domains.all.domains[*].name)
}
```

Will return:

```terraform
data.opslevel_domains.all: Reading...
data.opslevel_domains.all: Read complete after 1s

Changes to Outputs:
  + all          = [
      + {
          + aliases     = [
              + "alias_1",
              + "alias_2",
            ]
          + description = null
          + id          = "Z2lkOi8vb3BzbGV2ZWwvRW50aXR5T2JqZWN0LzEwOTk2NjU"
          + name        = "Big Domain"
          + owner       = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNDQwNg"
        },
      + {
          + aliases     = [
              + "example_1",
              + "example_2",
            ]
          + description = "Example description"
          + id          = "Z2lkOi8vb3BzbGV2ZWwvRW50aXR5T2JqZWN0LzE4ODg2NzA"
          + name        = "Example Domain"
          + owner       = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS85NzU5"
        },
    ]
  + domain_names = [
      + "Big Domain",
      + "Example Domain",
    ]
```

### Output Example - `data.opslevel_domains` before version 0.12.0-0

Given this configuration:

```terraform
data "opslevel_domains" "all" {}

output "all" {
  value = data.opslevel_domains.all
}

output "domain_names" {
  value = sort(data.opslevel_domains.all.names)
}
```

Will return:

```terraform
data.opslevel_domains.all: Reading...
data.opslevel_domains.all: Read complete after 1s

Changes to Outputs:
  + all          = [
      + aliases      = [
          + "alias_1",
          + "alias_2",
          + "example_1",
          + "example_2",
        ]
      + descriptions = [
          + "",
          + "Example description",
        ]
      + ids          = [
          + "Z2lkOi8vb3BzbGV2ZWwvRW50aXR5T2JqZWN0LzEwOTk2NjU",
          + "Z2lkOi8vb3BzbGV2ZWwvRW50aXR5T2JqZWN0LzE4ODg2NzA",
        ]
      + names        = [
          + "Big Domain",
          + "Example Domain",
        ]
      + owners       = [
          + "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNDQwNg",
          + "Z2lkOi8vb3BzbGV2ZWwvVGVhbS85NzU5",
        ]
    }
  + domain_names = [
      + "Big Domain",
      + "Example Domain",
    ]
```

# State Migration

The structure of some resource data has changed from v0.11.0 to v1.0.0.
If you encounter erros like `Error: Unable to Read Previously Saved State for UpgradeResourceState`
while upgrading to the OpsLevel v1.0.0 Terraform provider, the following will get you unstuck.

### Get ids of resources that raise errors

Assuming `opslevel_check_tag_defined.my_tag` raises an error like
`AttributeName("environment_predicate"): invalid JSON, expected "{", got "["`, run
`terraform state show opslevel_check_tag_defined.my_tag` and extract the value of `id`.
Extract and keep all ids for resources that raise errors.

If `terraform state show ...` raises an error, try `terraform state pull > my_tf_state.json`
and extract ids from `my_tf_state.json`.

### Remove problematic resources from state

This does NOT remove the resources from OpsLevel - only from Terraform's state.
Run `terraform state rm opslevel_check_tag_defined.my_tag` and repeat for all
resources that raise errors.
Optionally verify resources removed from state with `terraform state list`.

### Import resources from OpsLevel into Terraform state

Import OpsLevel resources into Terraform's state by running
`terraform import opslevel_check_tag_defined.my_tag <resource-id>` and repeat
for all resources that raise errors. Note <resource-id> matches the `id` of the OpsLevel
resource as returned by the [OpsLevel cli](https://github.com/opslevel/cli), if `id` is
not otherwise known.
Optionally verify resources imported to state with `terraform state list`.

