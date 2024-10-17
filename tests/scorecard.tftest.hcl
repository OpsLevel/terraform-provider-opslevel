variables {
  resource_name = "opslevel_scorecard"

  # required fields
  affects_overall_service_levels = true
  name                           = "TF Test Scorecard"
  owner_id                       = null # sourced from module

  # optional fields
  description = "TF Scorecard description"
  filter_id   = null # sourced from module
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_filters.all,
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_scorecard_create_with_all_fields" {

  variables {
    filter_id = run.from_data_module.first_filter.id
    owner_id  = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/scorecard"
  }

  assert {
    condition = alltrue([
      can(opslevel_scorecard.this.affects_overall_service_levels),
      can(opslevel_scorecard.this.aliases),
      can(opslevel_scorecard.this.categories),
      can(opslevel_scorecard.this.description),
      can(opslevel_scorecard.this.filter_id),
      can(opslevel_scorecard.this.id),
      can(opslevel_scorecard.this.name),
      can(opslevel_scorecard.this.owner_id),
      can(opslevel_scorecard.this.passing_checks),
      can(opslevel_scorecard.this.service_count),
      can(opslevel_scorecard.this.total_checks),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_scorecard.this.affects_overall_service_levels == var.affects_overall_service_levels
    error_message = format(
      "expected '%v' but got '%v'",
      var.affects_overall_service_levels,
      opslevel_scorecard.this.affects_overall_service_levels,
    )
  }

  assert {
    condition = opslevel_scorecard.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_scorecard.this.description,
    )
  }

  assert {
    condition = opslevel_scorecard.this.filter_id == var.filter_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter_id,
      opslevel_scorecard.this.filter_id,
    )
  }

  assert {
    condition     = startswith(opslevel_scorecard.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_scorecard.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_scorecard.this.name,
    )
  }

  assert {
    condition = opslevel_scorecard.this.owner_id == var.owner_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner_id,
      opslevel_scorecard.this.owner_id,
    )
  }

}

run "resource_scorecard_unset_optional_fields" {

  variables {
    description = null
    filter_id   = null
    owner_id    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/scorecard"
  }

  assert {
    condition     = opslevel_scorecard.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_scorecard.this.filter_id == null
    error_message = var.error_expected_null_field
  }

}

run "delete_scorecard_outside_of_terraform" {

  variables {
    resource_id   = run.resource_scorecard_create_with_all_fields.this.id
    resource_type = "scorecard"
  }

  module {
    source = "./provisioner"
  }
}

run "resource_scorecard_create_with_required_fields" {

  variables {
    # other fields from file scoped variables block
    description = null
    filter_id   = null
    owner_id    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/scorecard"
  }

  assert {
    condition = run.resource_scorecard_create_with_all_fields.this.id != opslevel_scorecard.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_scorecard_create_with_all_fields.this.id,
      opslevel_scorecard.this.id,
    )
  }

  assert {
    condition     = opslevel_scorecard.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_scorecard.this.filter_id == null
    error_message = var.error_expected_null_field
  }

}

run "resource_scorecard_set_all_fields" {

  variables {
    # other fields from file scoped variables block
    affects_overall_service_levels = !var.affects_overall_service_levels
    filter_id                      = run.from_data_module.first_filter.id
    owner_id                       = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/scorecard"
  }

  assert {
    condition = opslevel_scorecard.this.affects_overall_service_levels == var.affects_overall_service_levels
    error_message = format(
      "expected '%v' but got '%v'",
      var.affects_overall_service_levels,
      opslevel_scorecard.this.affects_overall_service_levels,
    )
  }

  assert {
    condition = opslevel_scorecard.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_scorecard.this.description,
    )
  }

  assert {
    condition = opslevel_scorecard.this.filter_id == var.filter_id
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter_id,
      opslevel_scorecard.this.filter_id,
    )
  }

}
