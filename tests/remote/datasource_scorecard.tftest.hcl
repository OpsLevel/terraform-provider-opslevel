run "datasource_scorecards_all" {

  variables {
    datasource_type = "opslevel_scorecards"
  }

  assert {
    condition     = can(data.opslevel_scorecards.all.scorecards)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_scorecards.all.scorecards) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_scorecard_first" {

  variables {
    datasource_type = "opslevel_scorecard"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_scorecard.first_scorecard_by_id.affects_overall_service_levels),
      can(data.opslevel_scorecard.first_scorecard_by_id.aliases),
      can(data.opslevel_scorecard.first_scorecard_by_id.description),
      can(data.opslevel_scorecard.first_scorecard_by_id.filter_id),
      can(data.opslevel_scorecard.first_scorecard_by_id.id),
      can(data.opslevel_scorecard.first_scorecard_by_id.identifier),
      can(data.opslevel_scorecard.first_scorecard_by_id.name),
      can(data.opslevel_scorecard.first_scorecard_by_id.owner_id),
      can(data.opslevel_scorecard.first_scorecard_by_id.passing_checks),
      can(data.opslevel_scorecard.first_scorecard_by_id.service_count),
      can(data.opslevel_scorecard.first_scorecard_by_id.total_checks),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_scorecard.first_scorecard_by_id.id == data.opslevel_scorecards.all.scorecards[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

}
