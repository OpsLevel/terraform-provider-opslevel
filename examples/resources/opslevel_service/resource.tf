data "opslevel_lifecycle" "beta" {
  filter {
    field = "alias"
    value = "beta"
  }
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
  aliases          = ["bar", "baz"]

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
  owner           = opslevel_team.foo.alias

  api_document_path             = "/swagger.json"
  preferred_api_document_source = "PULL" //or "PUSH"

  aliases = ["bar", "baz"]
  tags    = ["foo:bar"]
}

output "foo_aliases" {
  value = opslevel_service.foo.aliases
}
