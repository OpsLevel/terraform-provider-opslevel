run "datasource_scorecards_all" {

  variables {
    datasource_type = "opslevel_scorecards"
  }

  module {
    source = "./scorecard"
  }

  assert {
    condition     = can(data.opslevel_scorecards.all.scorecards)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_scorecards.all.scorecards) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

}

run "datasource_scorecard_first" {

  variables {
    datasource_type = "opslevel_scorecard"
  }

  module {
    source = "./scorecard"
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
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_scorecard.first_scorecard_by_id.id == data.opslevel_scorecards.all.scorecards[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

}

run "resource_scorecard_create_with_all_fields" {

  variables {
  }

  module {
    source = "./scorecard"
  }

}

run "resource_scorecard_update_unset_optional_fields" {

  variables {
  }

  module {
    source = "./scorecard"
  }

}

run "resource_scorecard_update_set_optional_fields" {

  variables {
  }

  module {
    source = "./scorecard"
  }

}
