# TODO: Scorecard tests works on orange. Need to add to PAT acct.
#run "datasource_scorecards_all" {
#
#  assert {
#    condition     = length(data.opslevel_scorecards.all.scorecards) > 0
#    error_message = "zero scorecards found in data.opslevel_scorecards"
#  }
#
#  assert {
#    condition = alltrue([
#      can(data.opslevel_scorecards.all.scorecards[0].id),
#    ])
#    error_message = "cannot set all expected scorecard datasource fields"
#  }
#
#}
#
#run "datasource_scorecard_first" {
#
#  assert {
#    condition     = data.opslevel_scorecard.first_scorecard_by_id.id == data.opslevel_scorecards.all.scorecards[0].id
#    error_message = "wrong ID on opslevel_scorecard"
#  }
#
#}
