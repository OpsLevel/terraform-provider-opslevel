data "opslevel_team" "devs" {
  alias = "developers"
}

data "opslevel_rubric_level" "bronze" {
  filter {
    field = "name"
    value = "Bronze"
  }
}

data "opslevel_filter" "tier_1" {
  filter {
    field = "name"
    value = "Tier 1 Services"
  }
}

resource "opslevel_scorecard" "my_scorecard" {
  name                           = "My Scorecard"
  affects_overall_service_levels = true
  description                    = "This is my example scorecard"
  owner_id                       = data.opslevel_team.devs.id
  filter_id                      = data.opslevel_filter.tier_1.id
}

// Example of how to assign a check to a scorecard
resource "opslevel_check_manual" "my_check" {
  name                    = "My check that uses a scorecard"
  category                = opslevel_scorecard.my_scorecard.categories[0]
  level                   = data.opslevel_rubric_level.bronze.id
  update_requires_comment = true
}
