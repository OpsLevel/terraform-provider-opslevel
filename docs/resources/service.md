---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_service Resource - terraform-provider-opslevel"
subcategory: ""
description: |-
  Service Resource
---

# opslevel_service (Resource)

Service Resource

## Example Usage

```terraform
data "opslevel_lifecycle" "beta" {
  filter {
    field = "alias"
    value = "beta"
  }
}

data "opslevel_system" "parent" {
  identifier = "example"
}

data "opslevel_tier" "tier3" {
  filter {
    field = "index"
    value = "3"
  }
}

resource "opslevel_team" "foo" {
  name             = "foo"
  responsibilities = "Responsible for foo frontend and backend"
  aliases          = ["foo", "bar", "baz"] # NOTE: if set, slugified value of "name" must be included

  member {
    email = "john.doe@example.com"
    role  = "manager"
  }
}

resource "opslevel_service" "foo" {
  name = "foo"

  description = "foo service"
  framework   = "rails"
  language    = "ruby"

  lifecycle_alias = data.opslevel_lifecycle.beta.alias
  tier_alias      = data.opslevel_tier.tier3.alias
  owner           = opslevel_team.foo.id
  parent          = data.opslevel_system.parent.id

  api_document_path             = "/swagger.json"
  preferred_api_document_source = "PULL" //or "PUSH"

  aliases = ["foo", "bar", "baz"] # NOTE: if set, value of "name" must be included
  tags    = ["foo:bar"]
}

output "foo_aliases" {
  value = opslevel_service.foo.aliases
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The display name of the service.

### Optional

- `aliases` (Set of String) A list of human-friendly, unique identifiers for the service.
- `api_document_path` (String) The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path.
- `description` (String) A brief description of the service.
- `framework` (String) The primary software development framework that the service uses.
- `language` (String) The primary programming language that the service is written in.
- `lifecycle_alias` (String) The lifecycle stage of the service.
- `note` (String) Additional information about the service.
- `owner` (String) The team that owns the service. ID or Alias may be used.
- `parent` (String) The id or alias of the parent system of this service
- `preferred_api_document_source` (String) The API document source (PULL or PUSH) used to determine the displayed document. If null, defaults to PUSH.
- `product` (String) A product is an application that your end user interacts with. Multiple services can work together to power a single product.
- `tags` (Set of String) A list of tags applied to the service.
- `tier_alias` (String) The software tier that the service belongs to.
- `type` (String) The component type of the service.

### Read-Only

- `id` (String) The id of the service to find

## Import

Import is supported using the following syntax:

```shell
terraform import opslevel_service.example Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS82MDI0
```
