data "opslevel_team" "devs" {
  alias = "developers"
}

data "opslevel_filter" "tier_1" {
  filter {
    field = "name"
    value = "Tier 1 Services"
  }
}

resource "opslevel_scorecard" "my_scorecard" {
    name = "My Scorecard"
    description = "This is my example scorecard"
    ownerId = data.opslevel_team.devs.id
    filterId = data.opslevel_filter.tier_1.id
}

// Example of how to assign a check to a scorecard
resource "opslevel_check_manual" "my_check" {
  name      = "My Check"
  category  = resource.opslevel_scorecard.my_scorecard.id
  ...
}
