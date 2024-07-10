variables {
  scorecard_one  = "opslevel_scorecard"
  scorecards_all = "opslevel_scorecards"

  # required fields
  affects_overall_service_levels = true
  name                           = "TF Test Scorecard"
  owner_id                       = null

  # optional fields
  description = "TF Scorecard description"
  filter_id   = null
}

run "from_filter_get_filter_id" {
  command = plan

  variables {
    connective     = null
    predicate_list = null
  }

  module {
    source = "./filter"
  }
}

run "from_team_get_owner_id" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }
}

run "resource_scorecard_create_with_all_fields" {

  variables {
    affects_overall_service_levels = var.affects_overall_service_levels
    description                    = var.description
    filter_id                      = run.from_filter_get_filter_id.first_filter.id
    name                           = var.name
    owner_id                       = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./scorecard"
  }

  assert {
    condition = alltrue([
      can(opslevel_scorecard.test.affects_overall_service_levels),
      can(opslevel_scorecard.test.aliases),
      can(opslevel_scorecard.test.categories),
      can(opslevel_scorecard.test.description),
      can(opslevel_scorecard.test.filter_id),
      can(opslevel_scorecard.test.id),
      can(opslevel_scorecard.test.name),
      can(opslevel_scorecard.test.owner_id),
      can(opslevel_scorecard.test.passing_checks),
      can(opslevel_scorecard.test.service_count),
      can(opslevel_scorecard.test.total_checks),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.scorecard_one)
  }

  assert {
    condition     = opslevel_scorecard.test.affects_overall_service_levels == var.affects_overall_service_levels
    error_message = "wrong affects_overall_service_levels for opslevel_scorecard resource"
  }

  assert {
    condition     = opslevel_scorecard.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.scorecard_one)
  }

  assert {
    condition     = opslevel_scorecard.test.filter_id == var.filter_id
    error_message = "wrong filter_id for opslevel_scorecard resource"
  }

  assert {
    condition     = startswith(opslevel_scorecard.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.scorecard_one)
  }

  assert {
    condition     = opslevel_scorecard.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.scorecard_one)
  }

  assert {
    condition     = opslevel_scorecard.test.owner_id == var.owner_id
    error_message = "wrong owner_id for opslevel_scorecard resource"
  }

}

run "resource_scorecard_create_with_empty_optional_fields" {

  variables {
    affects_overall_service_levels = var.affects_overall_service_levels
    description                    = ""
    owner_id                       = run.from_team_get_owner_id.first_team.id
    name                           = "New ${var.name} with empty fields"
  }

  module {
    source = "./scorecard"
  }

  assert {
    condition     = opslevel_scorecard.test.description == ""
    error_message = var.error_expected_empty_string
  }

}

run "resource_scorecard_update_unset_optional_fields" {

  variables {
    description = null
    filter_id   = null
    owner_id    = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./scorecard"
  }

  assert {
    condition     = opslevel_scorecard.test.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_scorecard.test.filter_id == null
    error_message = var.error_expected_null_field
  }

}

run "resource_scorecard_update_set_all_fields" {

  variables {
    affects_overall_service_levels = !var.affects_overall_service_levels
    description                    = "${var.description} updated"
    filter_id                      = run.from_filter_get_filter_id.first_filter.id
    name                           = "${var.name} updated"
    owner_id                       = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./scorecard"
  }

  assert {
    condition     = opslevel_scorecard.test.affects_overall_service_levels == var.affects_overall_service_levels
    error_message = "wrong affects_overall_service_levels for opslevel_scorecard resource"
  }

  assert {
    condition     = opslevel_scorecard.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.scorecard_one)
  }

  assert {
    condition     = opslevel_scorecard.test.filter_id == var.filter_id
    error_message = "wrong filter_id for opslevel_scorecard resource"
  }

  assert {
    condition     = opslevel_scorecard.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.scorecard_one)
  }

  assert {
    condition     = opslevel_scorecard.test.owner_id == var.owner_id
    error_message = "wrong owner_id for opslevel_scorecard resource"
  }


}

run "datasource_scorecards_all" {

  variables {
    owner_id = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./scorecard"
  }

  assert {
    condition     = can(data.opslevel_scorecards.all.scorecards)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.scorecards_all)
  }

  assert {
    condition     = length(data.opslevel_scorecards.all.scorecards) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.scorecards_all)
  }

}

run "datasource_scorecard_first" {

  variables {
    owner_id = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./scorecard"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_scorecard.first_scorecard_by_id.affects_overall_service_levels),
      can(data.opslevel_scorecard.first_scorecard_by_id.aliases),
      can(data.opslevel_scorecard.first_scorecard_by_id.categories),
      can(data.opslevel_scorecard.first_scorecard_by_id.description),
      can(data.opslevel_scorecard.first_scorecard_by_id.filter_id),
      can(data.opslevel_scorecard.first_scorecard_by_id.id),
      can(data.opslevel_scorecard.first_scorecard_by_id.name),
      can(data.opslevel_scorecard.first_scorecard_by_id.owner_id),
      can(data.opslevel_scorecard.first_scorecard_by_id.passing_checks),
      can(data.opslevel_scorecard.first_scorecard_by_id.service_count),
      can(data.opslevel_scorecard.first_scorecard_by_id.total_checks),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.scorecard_one)
  }

  assert {
    condition     = data.opslevel_scorecard.first_scorecard_by_id.id == data.opslevel_scorecards.all.scorecards[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.scorecard_one)
  }

}
