resource "opslevel_alias" "service" {
  resource_type       = "service"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "team" {
  resource_type       = "team"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "domain" {
  resource_type       = "domain"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "system" {
  resource_type       = "system"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "infra" {
  resource_type       = "infrastructure_resource"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "scorecard" {
  resource_type       = "scorecard"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}