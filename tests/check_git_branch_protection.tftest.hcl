variables {
  resource_name = "opslevel_check_git_branch_protection"

  # -- check_git_branch_protection fields --
  # only uses base fields

  # -- check base fields --
  # required fields
  category = null # sourced from module
  level    = null # sourced from module
  name     = "TF Test Check Git Branch Protection"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null # sourced from module
  notes     = "Notes on Git Branch Protection Check"
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

run "resource_check_git_branch_protection_create_with_all_fields" {

  variables {
    category  = run.from_data_module.first_rubric_category.id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_data_module.first_filter.id
    level     = run.from_data_module.max_index_rubric_level.id
    name      = var.name
    notes     = var.notes
    owner     = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/git_branch_protection"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_git_branch_protection.this.category),
      can(opslevel_check_git_branch_protection.this.description),
      can(opslevel_check_git_branch_protection.this.enable_on),
      can(opslevel_check_git_branch_protection.this.enabled),
      can(opslevel_check_git_branch_protection.this.filter),
      can(opslevel_check_git_branch_protection.this.id),
      can(opslevel_check_git_branch_protection.this.level),
      can(opslevel_check_git_branch_protection.this.name),
      can(opslevel_check_git_branch_protection.this.notes),
      can(opslevel_check_git_branch_protection.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_git_branch_protection.this.category,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_git_branch_protection.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_git_branch_protection.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_git_branch_protection.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_git_branch_protection.this.filter,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_git_branch_protection.this.level,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_git_branch_protection.this.name,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_git_branch_protection.this.notes,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_git_branch_protection.this.owner,
    )
  }

}

run "resource_check_git_branch_protection_update_unset_optional_fields" {

  variables {
    category  = run.from_data_module.first_rubric_category.id
    enable_on = null
    enabled   = null
    filter    = null
    level     = run.from_data_module.max_index_rubric_level.id
    notes     = null
    owner     = null
  }

  module {
    source = "./opslevel_modules/modules/check/git_branch_protection"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "delete_check_alert_source_usage_outside_of_terraform" {

  variables {
    resource_id   = run.resource_check_git_branch_protection_create_with_all_fields.this.id
    resource_type = "check"
  }

  module {
    source = "./provisioner"
  }
}

run "resource_check_git_branch_protection_create_with_required_fields" {

  variables {
    category  = run.from_data_module.first_rubric_category.id
    enable_on = null
    enabled   = null
    filter    = null
    level     = run.from_data_module.max_index_rubric_level.id
    notes     = null
    owner     = null
  }

  module {
    source = "./opslevel_modules/modules/check/git_branch_protection"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_git_branch_protection_update_set_all_fields" {

  variables {
    category  = run.from_data_module.first_rubric_category.id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_data_module.first_filter.id
    level     = run.from_data_module.max_index_rubric_level.id
    name      = var.name
    notes     = var.notes
    owner     = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/git_branch_protection"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_git_branch_protection.this.category),
      can(opslevel_check_git_branch_protection.this.description),
      can(opslevel_check_git_branch_protection.this.enable_on),
      can(opslevel_check_git_branch_protection.this.enabled),
      can(opslevel_check_git_branch_protection.this.filter),
      can(opslevel_check_git_branch_protection.this.id),
      can(opslevel_check_git_branch_protection.this.level),
      can(opslevel_check_git_branch_protection.this.name),
      can(opslevel_check_git_branch_protection.this.notes),
      can(opslevel_check_git_branch_protection.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_git_branch_protection.this.category,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_git_branch_protection.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_git_branch_protection.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_git_branch_protection.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_git_branch_protection.this.filter,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_git_branch_protection.this.level,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_git_branch_protection.this.name,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_git_branch_protection.this.notes,
    )
  }

  assert {
    condition = opslevel_check_git_branch_protection.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_git_branch_protection.this.owner,
    )
  }

}
