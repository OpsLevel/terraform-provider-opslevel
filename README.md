# terraform-provider-opslevel

Terraform Provider for OpsLevel.com

## Example

```hcl
terraform {
  required_providers {
    opslevel = {
      source = "zapier/opslevel"
    }
  }
}

provider "opslevel" {
  # token = "eyJhbGciOi..." // or environment variable OPSLEVEL_TOKEN
}

resource "opslevel_service" "foo" {
  name = "foo"

  description = "foo service"
  framework   = "rails"
  language    = "ruby"

  aliases = [
    "bar"
  ]

  tags = {
    foo = "bar"
  }
}

# Datasources to read and filter on Services

data "opslevel_service" "foo" {
  filter {
    field = "alias"
    value = opslevel_service.foo.aliases.0
  }
}

data "opslevel_service" "django" {
  filter {
    field = "framework"
    value = "django"
  }
}
```

