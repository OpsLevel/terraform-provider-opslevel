<p align="center">
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/terraform-provider-opslevel.svg" /></a>
    <a href="http://golang.org" alt="Made With Go">
        <img src="https://img.shields.io/github/go-mod/go-version/OpsLevel/terraform-provider-opslevel" /></a>
    <a href="https://GitHub.com/OpsLevel/terraform-provider-opslevel/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/terraform-provider-opslevel?include_prereleases" /></a>  
    <a href="https://GitHub.com/OpsLevel/terraform-provider-opslevel/issues/" alt="Issues">
        <img src="https://img.shields.io/github/issues/OpsLevel/terraform-provider-opslevel.svg" /></a>  
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/terraform-provider-opslevel" /></a>
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/terraform-provider-opslevel" /></a>
    <a href="https://github.com/OpsLevel/terraform-provider-opslevel/releases" alt="Downloads">
        <img src="https://img.shields.io/github/downloads/OpsLevel/terraform-provider-opslevel/total" /></a>
</p>

Terraform Provider for [OpsLevel](https://opslevel.com)
===============================

[Documentation](https://registry.terraform.io/providers/OpsLevel/opslevel/latest/docs)

## Using the provider

```terraform
provider "opslevel" {
  apitoken = "XXX" // or environment variable OPSLEVEL_APITOKEN
}

resource "opslevel_team" "foo" {
  name = "foo"
  manager_email = "foo@example.com"
  responsibilities = "Responsible for foo frontend and backend"
}

resource "opslevel_service" "foo-frontend" {
  name = "foo-frontend"

  description = "The foo frontend service"
  framework   = "rails"
  language    = "ruby"

  lifecycle_alias = "beta"
  tier_alias = "tier_3"
  owner_alias = opslevel_team.foo.alias

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

# Useful Snippets

Get Tag Values

```hcl
values({for entry in data.opslevel_service.example.tags : entry => split(":", entry)[1]})
```
