resource "opslevel_alias" "service" {
  resource_type       = "Service"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "team" {
  resource_type       = "Team"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "domain" {
  resource_type       = "Domain"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "system" {
  resource_type       = "System"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "infra" {
  resource_type       = "Infrastructure_Resource"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}

resource "opslevel_alias" "scorecard" {
  resource_type       = "Scorecard"
  resource_identifier = "example_alias"

  aliases = ["example_alias_2", "example_alias_3"]
}