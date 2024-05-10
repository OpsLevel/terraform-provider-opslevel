data "opslevel_scorecards" "all" {}

data "opslevel_scorecard" "first_scorecard_by_id" {
  identifier = data.opslevel_scorecards.all.scorecards[0].id
}
