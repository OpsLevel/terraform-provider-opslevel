variables {
  resource_name = "opslevel_check_package_version"

  # -- check_package_version fields --
  # required fields
  package_constraint = "matches_version"
  package_manager    = "docker"
  package_name       = "foobar"

  # optional fields
  missing_package_result = "passed"
  package_name_is_regex  = true
  version_constraint_predicate = {
    type  = "does_not_match_regex"
    value = "^$"
  }

  # -- check base fields --
  # required fields
  category = null # sourced from module
  level    = null # sourced from module
  name     = "TF Test Check package_version"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null # sourced from module
  notes     = "Notes on package_version check"
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

run "resource_check_package_version_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_package_version.this.category),
      can(opslevel_check_package_version.this.description),
      can(opslevel_check_package_version.this.enable_on),
      can(opslevel_check_package_version.this.enabled),
      can(opslevel_check_package_version.this.filter),
      can(opslevel_check_package_version.this.id),
      can(opslevel_check_package_version.this.level),
      can(opslevel_check_package_version.this.missing_package_result),
      can(opslevel_check_package_version.this.name),
      can(opslevel_check_package_version.this.notes),
      can(opslevel_check_package_version.this.owner),
      can(opslevel_check_package_version.this.package_constraint),
      can(opslevel_check_package_version.this.package_manager),
      can(opslevel_check_package_version.this.package_name),
      can(opslevel_check_package_version.this.package_name_is_regex),
      can(opslevel_check_package_version.this.version_constraint_predicate),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_package_version.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_package_version.this.category,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_package_version.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_package_version.this.enabled,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_package_version.this.filter,
    )
  }

  assert {
    condition     = startswith(opslevel_check_package_version.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_package_version.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_package_version.this.level,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.missing_package_result == var.missing_package_result
    error_message = format(
      "expected '%v' but got '%v'",
      var.missing_package_result,
      opslevel_check_package_version.this.missing_package_result,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_package_version.this.name,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_package_version.this.notes,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_package_version.this.owner,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_constraint == var.package_constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_constraint,
      opslevel_check_package_version.this.package_constraint,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_manager == var.package_manager
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_manager,
      opslevel_check_package_version.this.package_manager,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_name == var.package_name
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name,
      opslevel_check_package_version.this.package_name,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_name_is_regex == var.package_name_is_regex
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name_is_regex,
      opslevel_check_package_version.this.package_name_is_regex,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.version_constraint_predicate == var.version_constraint_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.version_constraint_predicate,
      opslevel_check_package_version.this.version_constraint_predicate,
    )
  }

}

run "resource_check_package_version_unset_optional_fields" {

  variables {
    # other fields from file scoped variables block
    category              = run.from_data_module.first_rubric_category.id
    enable_on             = null
    enabled               = null
    filter                = null
    level                 = run.from_data_module.max_index_rubric_level.id
    notes                 = null
    owner                 = null
    package_name_is_regex = null
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition = opslevel_check_package_version.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_package_version.this.category,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_package_version.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_package_version resource"
  }

  assert {
    condition     = opslevel_check_package_version.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_package_version.this.level,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.missing_package_result == var.missing_package_result
    error_message = format(
      "expected '%v' but got '%v'",
      var.missing_package_result,
      opslevel_check_package_version.this.missing_package_result,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_package_version.this.name,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_package_version.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.package_constraint == var.package_constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_constraint,
      opslevel_check_package_version.this.package_constraint,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_manager == var.package_manager
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_manager,
      opslevel_check_package_version.this.package_manager,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_name == var.package_name
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name,
      opslevel_check_package_version.this.package_name,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.package_name_is_regex == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.version_constraint_predicate == var.version_constraint_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.version_constraint_predicate,
      opslevel_check_package_version.this.version_constraint_predicate,
    )
  }

}

run "delete_check_package_version_outside_of_terraform" {

  variables {
    command = "delete check ${run.resource_check_package_version_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_check_package_version_create_with_required_fields" {

  variables {
    # other fields from file scoped variables block
    category              = run.from_data_module.first_rubric_category.id
    enable_on             = null
    enabled               = null
    filter                = null
    level                 = run.from_data_module.max_index_rubric_level.id
    notes                 = null
    owner                 = null
    package_name_is_regex = null
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition = run.resource_check_package_version_create_with_all_fields.this.id != opslevel_check_package_version.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_check_package_version_create_with_all_fields.this.id,
      opslevel_check_package_version.this.id,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_package_version.this.category,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_package_version.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_package_version resource"
  }

  assert {
    condition     = opslevel_check_package_version.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_package_version.this.level,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.missing_package_result == var.missing_package_result
    error_message = format(
      "expected '%v' but got '%v'",
      var.missing_package_result,
      opslevel_check_package_version.this.missing_package_result,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_package_version.this.name,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_package_version.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.package_constraint == var.package_constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_constraint,
      opslevel_check_package_version.this.package_constraint,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_manager == var.package_manager
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_manager,
      opslevel_check_package_version.this.package_manager,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_name == var.package_name
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name,
      opslevel_check_package_version.this.package_name,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.package_name_is_regex == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.version_constraint_predicate == var.version_constraint_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.version_constraint_predicate,
      opslevel_check_package_version.this.version_constraint_predicate,
    )
  }

}

run "resource_check_package_version_set_all_fields" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition = opslevel_check_package_version.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_package_version.this.category,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_package_version.this.enabled,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.filter == var.filter
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_package_version.this.level,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.missing_package_result == var.missing_package_result
    error_message = format(
      "expected '%v' but got '%v'",
      var.missing_package_result,
      opslevel_check_package_version.this.missing_package_result,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_package_version.this.name,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_package_version.this.notes,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_package_version.this.owner,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_constraint == var.package_constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_constraint,
      opslevel_check_package_version.this.package_constraint,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_manager == var.package_manager
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_manager,
      opslevel_check_package_version.this.package_manager,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_name == var.package_name
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name,
      opslevel_check_package_version.this.package_name,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.package_name_is_regex == var.package_name_is_regex
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name_is_regex,
      opslevel_check_package_version.this.package_name_is_regex,
    )
  }

  assert {
    condition = opslevel_check_package_version.this.version_constraint_predicate == var.version_constraint_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.version_constraint_predicate,
      opslevel_check_package_version.this.version_constraint_predicate,
    )
  }

}

run "resource_check_package_version_set_package_constraint_does_not_exist" {

  variables {
    # other fields from file scoped variables block
    category                     = run.from_data_module.first_rubric_category.id
    filter                       = run.from_data_module.first_filter.id
    level                        = run.from_data_module.max_index_rubric_level.id
    owner                        = run.from_data_module.first_team.id
    missing_package_result       = null
    package_constraint           = "does_not_exist"
    version_constraint_predicate = null
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition     = opslevel_check_package_version.this.missing_package_result == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.package_constraint == var.package_constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_constraint,
      opslevel_check_package_version.this.package_constraint,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.version_constraint_predicate == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_package_version_set_package_constraint_exists" {

  variables {
    # other fields from file scoped variables block
    category                     = run.from_data_module.first_rubric_category.id
    filter                       = run.from_data_module.first_filter.id
    level                        = run.from_data_module.max_index_rubric_level.id
    owner                        = run.from_data_module.first_team.id
    missing_package_result       = null
    package_constraint           = "exists"
    version_constraint_predicate = null
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition     = opslevel_check_package_version.this.missing_package_result == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_package_version.this.package_constraint == var.package_constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_constraint,
      opslevel_check_package_version.this.package_constraint,
    )
  }

  assert {
    condition     = opslevel_check_package_version.this.version_constraint_predicate == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_package_version_set_missing_package_result_failed" {

  variables {
    # other fields from file scoped variables block
    category               = run.from_data_module.first_rubric_category.id
    filter                 = run.from_data_module.first_filter.id
    level                  = run.from_data_module.max_index_rubric_level.id
    missing_package_result = "failed"
    owner                  = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition = opslevel_check_package_version.this.missing_package_result == var.missing_package_result
    error_message = format(
      "expected '%v' but got '%v'",
      var.missing_package_result,
      opslevel_check_package_version.this.missing_package_result,
    )
  }

}

run "resource_check_package_version_set_predicate_matches_regex" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
    version_constraint_predicate = {
      type  = "matches_regex"
      value = "^$"
    }
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition = opslevel_check_package_version.this.version_constraint_predicate == var.version_constraint_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.version_constraint_predicate,
      opslevel_check_package_version.this.version_constraint_predicate,
    )
  }

}

run "resource_check_package_version_set_predicate_satisfies_version_constraint" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
    version_constraint_predicate = {
      type  = "satisfies_version_constraint"
      value = "^$"
    }
  }

  module {
    source = "./opslevel_modules/modules/check/package_version"
  }

  assert {
    condition = opslevel_check_package_version.this.version_constraint_predicate == var.version_constraint_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.version_constraint_predicate,
      opslevel_check_package_version.this.version_constraint_predicate,
    )
  }

}
