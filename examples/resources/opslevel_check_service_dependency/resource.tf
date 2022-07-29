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

data "opslevel_team" "devs" {
  alias = "developers"
}

data "opslevel_filter" "tier1" {
  filter {
    field = "name"
    value = "Tier 1"
  }
}

resource "opslevel_check_service_dependency" "example" {
  name    = "foo"
  enabled = true
  # To set a future enable date remove field 'enabled' and use 'enable_on'
  # enable_on = "2022-05-23T14:14:18.782000Z"
  category = data.opslevel_rubric_category.security.id
  level    = data.opslevel_rubric_level.bronze.id
  owner    = data.opslevel_team.devs.id
  filter   = data.opslevel_filter.tier1.id
  notes    = "Optional additional info on why this check is run or how to fix it"
}
