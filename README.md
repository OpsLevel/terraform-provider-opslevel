<p align="center">
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/terraform-provider-opslevel.svg" /></a>
    <a href="http://golang.org" alt="Made With Go">
        <img src="https://img.shields.io/github/go-mod/go-version/OpsLevel/terraform-provider-opslevel" /></a>
    <a href="https://GitHub.com/OpsLevel/terraform-provider-opslevel/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/terraform-provider-opslevel?include_prereleases" /></a>  
    <a href="https://masterminds.github.io/stability/active.html" alt="Stability: Active">
        <img src="https://masterminds.github.io/stability/active.svg" /></a>   
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/terraform-provider-opslevel" /></a>
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/terraform-provider-opslevel" /></a>
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/releases" alt="Downloads">
        <img src="https://img.shields.io/github/downloads/OpsLevel/terraform-provider-opslevel/total" /></a>
</p>

[![Overall](https://img.shields.io/endpoint?style=flat&url=https%3A%2F%2Fapp.opslevel.com%2Fapi%2Fservice_level%2FOYbJw2HuOqY7Np42eBzMn_RCwWebqaywVSJAQczStEY)](https://app.opslevel.com/services/opslevel_terraform_provider/maturity-report)

Terraform Provider for [OpsLevel](https://opslevel.com)
===============================

[Provider Documentation](https://registry.terraform.io/providers/OpsLevel/opslevel/latest/docs)
[Quickstart](https://www.opslevel.com/docs/terraform)
[Importing All Existing Account Data](https://www.opslevel.com/docs/terraform/#Importing)

## Example

```hcl
provider "opslevel" {
  api_token = "XXX" // or environment variable OPSLEVEL_API_TOKEN
}

resource "opslevel_team" "foo" {
  name = "foo"
  responsibilities = "Responsible for foo frontend and backend"

  member {
    email = "foo@example.com"
    role = "manager"
  }
  member {
    email = "bar@example.com"
    role = "contributor"
  }
}

resource "opslevel_service" "foo-frontend" {
  name = "foo-frontend"

  description = "The foo frontend service"
  framework   = "rails"
  language    = "ruby"

  lifecycle_alias = "beta"
  tier_alias = "tier_3"
  owner = opslevel_team.foo.alias

  tags = [
    "environment:production",
  ]
}

data "opslevel_rubric_category" "security" {
  filter {
    field = "name"
    value = "Security"
  }
}

data "opslevel_rubric_level" "bronze" {
  filter {
    field = "name"
    value = "Bronze"
  }
}

resource "opslevel_filter" "filter" {
  name = "foo"
  predicate {
    key = "tier_index"
    type = "equals"
    value = "tier_3"
  }
  connective = "and"
}

resource "opslevel_check_repository_integrated" "foo" {
  name = "foo"
  enabled = true
  category = data.opslevel_rubric_category.security.id
  level = data.opslevel_rubric_level.bronze.id
  owner = opslevel_team.foo.id
  filter = opslevel_filter.filter.id
  notes = "Optional additional info on why this check is run or how to fix it"
}
```

# Useful Terraform Snippets

Get a service's tag keys

```hcl
output "service_tag_keys" {
    value = values({for entry in opslevel_service.example.tags : entry => split(":", entry)[0]})
}
```

Get a service's tag values

```hcl
output "service_tag_values" {
    value = values({for entry in opslevel_service.example.tags : entry => split(":", entry)[1]})
}
```
