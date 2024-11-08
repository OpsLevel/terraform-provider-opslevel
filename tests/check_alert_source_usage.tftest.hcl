variables {
  resource_name = "opslevel_check_alert_source_usage"

  # -- check_alert_source_usage fields --
  # required fields
  alert_type = "datadog"

  # optional fields
  alert_name_predicate = {
    type  = "exists"
    value = null
  }

  # -- check base fields --
  # required fields
  category = null # sourced from module
  level    = null # sourced from module
  name     = "TF Test Check Alert Source Usage"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null # sourced from module
  notes     = "Notes on Alert Source Usage Check"
  owner     = null # sourced from module
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_filters.all,
      data.opslevel_rubric_categories.all,
      data.opslevel_rubric_levels.all,
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_check_alert_source_usage_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    alert_name_predicate = var.alert_name_predicate
    category             = run.from_data_module.first_rubric_category.id
    filter               = run.from_data_module.first_filter.id
    level                = run.from_data_module.max_index_rubric_level.id
    owner                = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/alert_source_usage"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_alert_source_usage.this.alert_type),
      can(opslevel_check_alert_source_usage.this.category),
      can(opslevel_check_alert_source_usage.this.description),
      can(opslevel_check_alert_source_usage.this.enable_on),
      can(opslevel_check_alert_source_usage.this.enabled),
      can(opslevel_check_alert_source_usage.this.filter),
      can(opslevel_check_alert_source_usage.this.id),
      can(opslevel_check_alert_source_usage.this.level),
      can(opslevel_check_alert_source_usage.this.name),
      can(opslevel_check_alert_source_usage.this.notes),
      can(opslevel_check_alert_source_usage.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.alert_name_predicate == var.alert_name_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.alert_name_predicate,
      opslevel_check_alert_source_usage.this.alert_name_predicate,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.alert_type == var.alert_type
    error_message = format(
      "expected '%v' but got '%v'",
      var.alert_type,
      opslevel_check_alert_source_usage.this.alert_type,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_alert_source_usage.this.category,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_alert_source_usage.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_alert_source_usage.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_alert_source_usage.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_alert_source_usage.this.filter,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_alert_source_usage.this.level,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_alert_source_usage.this.name,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_alert_source_usage.this.notes,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_alert_source_usage.this.owner,
    )
  }

}

run "resource_check_alert_source_usage_unset_optional_fields" {

  variables {
    # other fields from file scoped variables block
    alert_name_predicate = null
    category             = run.from_data_module.first_rubric_category.id
    enable_on            = null
    enabled              = null
    filter               = null
    level                = run.from_data_module.max_index_rubric_level.id
    notes                = null
    owner                = null
  }

  module {
    source = "./opslevel_modules/modules/check/alert_source_usage"
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.alert_type == var.alert_type
    error_message = format(
      "expected '%v' but got '%v'",
      var.alert_type,
      opslevel_check_alert_source_usage.this.alert_type,
    )
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.alert_name_predicate == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_alert_source_usage.this.category,
    )
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_alert_source_usage.this.level,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_alert_source_usage.this.name,
    )
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "delete_check_alert_source_usage_outside_of_terraform" {

  variables {
    command = "delete check ${run.resource_check_alert_source_usage_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_check_alert_source_usage_create_with_required_fields" {

  variables {
    # other fields from file scoped variables block
    alert_name_predicate = null
    category             = run.from_data_module.first_rubric_category.id
    enable_on            = null
    enabled              = null
    filter               = null
    level                = run.from_data_module.max_index_rubric_level.id
    notes                = null
    owner                = null
  }

  module {
    source = "./opslevel_modules/modules/check/alert_source_usage"
  }

  assert {
    condition = run.resource_check_alert_source_usage_create_with_all_fields.this.id != opslevel_check_alert_source_usage.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_check_alert_source_usage_create_with_all_fields.this.id,
      opslevel_check_alert_source_usage.this.id,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.alert_type == var.alert_type
    error_message = format(
      "expected '%v' but got '%v'",
      var.alert_type,
      opslevel_check_alert_source_usage.this.alert_type,
    )
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.alert_name_predicate == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_alert_source_usage.this.category,
    )
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_alert_source_usage.this.level,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_alert_source_usage.this.name,
    )
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_alert_source_usage_set_all_fields" {

  variables {
    # other fields from file scoped variables block
    alert_name_predicate = var.alert_name_predicate
    category             = run.from_data_module.first_rubric_category.id
    filter               = run.from_data_module.first_filter.id
    level                = run.from_data_module.max_index_rubric_level.id
    owner                = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/alert_source_usage"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_alert_source_usage.this.category),
      can(opslevel_check_alert_source_usage.this.description),
      can(opslevel_check_alert_source_usage.this.enable_on),
      can(opslevel_check_alert_source_usage.this.enabled),
      can(opslevel_check_alert_source_usage.this.filter),
      can(opslevel_check_alert_source_usage.this.id),
      can(opslevel_check_alert_source_usage.this.level),
      can(opslevel_check_alert_source_usage.this.name),
      can(opslevel_check_alert_source_usage.this.notes),
      can(opslevel_check_alert_source_usage.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.alert_type == var.alert_type
    error_message = format(
      "expected '%v' but got '%v'",
      var.alert_type,
      opslevel_check_alert_source_usage.this.alert_type,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.alert_name_predicate == var.alert_name_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.alert_name_predicate,
      opslevel_check_alert_source_usage.this.alert_name_predicate,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_alert_source_usage.this.category,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_alert_source_usage.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_alert_source_usage.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_alert_source_usage.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_alert_source_usage.this.filter,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_alert_source_usage.this.level,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_alert_source_usage.this.name,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_alert_source_usage.this.notes,
    )
  }

  assert {
    condition = opslevel_check_alert_source_usage.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_alert_source_usage.this.owner,
    )
  }

}
