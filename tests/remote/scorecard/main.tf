data "opslevel_scorecards" "all" {}

data "opslevel_scorecard" "first_scorecard_by_id" {
  identifier = data.opslevel_scorecards.all.scorecards[0].id
}

resource "opslevel_scorecard" "test" {
  affects_overall_service_levels = var.affects_overall_service_levels
  description                    = var.description
  filter_id                      = var.filter_id
  name                           = var.name
  owner_id                       = var.owner_id
}
