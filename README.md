# terraform-provider-opslevel

Terraform Provider for OpsLevel.com

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.15.x
-	[Go](https://golang.org/doc/install) >= 1.16


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
```

### Create a service

```hcl
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
```

### Read and filter on Services

```terraform
# read a single service by `alias` or `id`
data "opslevel_services" "foo" {
  filter {
    field = "alias"
    value = opslevel_service.foo.aliases.0
  }
}

# retrieve several services matching framework
data "opslevel_services" "django" {
  filter {
    field = "framework"
    value = "django"
  }
}

# retrieve several services matching tag
data "opslevel_services" "production" {
  filter {
    field = "tag"
    value = "production:true"
  }
}
```

### Read Teams from OpsLevel

```terraform
data "opslevel_teams" "all" {}

output "all_teams" {
  value = data.opslevel_teams.all.teams
}
```

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:
```sh
$ go install
```