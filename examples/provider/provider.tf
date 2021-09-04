provider "opslevel" {
  api_token = "XXX" // or environment variable OPSLEVEL_API_TOKEN
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