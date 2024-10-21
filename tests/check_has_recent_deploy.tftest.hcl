variables {
  check_has_recent_deploy = "opslevel_check_has_recent_deploy"

  # -- check_has_recent_deploy fields --
  # required fields
  days = 5

  # -- check base fields --
  # required fields
  category = null # sourced from module
  level    = null # sourced from module
  name     = "TF Test Check Has Recent Deploy"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null # sourced from module
  notes     = "Notes on Has Recent Deploy Check"
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

run "resource_check_has_recent_deploy_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/has_recent_deploy"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_has_recent_deploy.this.category),
      can(opslevel_check_has_recent_deploy.this.days),
      can(opslevel_check_has_recent_deploy.this.description),
      can(opslevel_check_has_recent_deploy.this.enable_on),
      can(opslevel_check_has_recent_deploy.this.enabled),
      can(opslevel_check_has_recent_deploy.this.filter),
      can(opslevel_check_has_recent_deploy.this.id),
      can(opslevel_check_has_recent_deploy.this.level),
      can(opslevel_check_has_recent_deploy.this.name),
      can(opslevel_check_has_recent_deploy.this.notes),
      can(opslevel_check_has_recent_deploy.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_has_recent_deploy)
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_has_recent_deploy.this.category,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.days == var.days
    error_message = format(
      "expected '%v' but got '%v'",
      var.days,
      opslevel_check_has_recent_deploy.this.days,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_has_recent_deploy.this.enabled,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_has_recent_deploy.this.enable_on,
    )
  }

  assert {
    condition     = startswith(opslevel_check_has_recent_deploy.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_has_recent_deploy)
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_has_recent_deploy.this.filter,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_has_recent_deploy.this.level,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_has_recent_deploy.this.name,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_has_recent_deploy.this.notes,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_has_recent_deploy.this.owner,
    )
  }

}

run "resource_check_has_recent_deploy_unset_optional_fields" {

  variables {
    # other fields from file scoped variables block
    category  = run.from_data_module.first_rubric_category.id
    enable_on = null
    enabled   = null
    filter    = null
    level     = run.from_data_module.max_index_rubric_level.id
    notes     = null
    owner     = null
  }

  module {
    source = "./opslevel_modules/modules/check/has_recent_deploy"
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_has_recent_deploy.this.category,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.days == var.days
    error_message = format(
      "expected '%v' but got '%v'",
      var.days,
      opslevel_check_has_recent_deploy.this.days,
    )
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_has_recent_deploy.this.level,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_has_recent_deploy.this.name,
    )
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.owner == null
    error_message = var.error_expected_null_field
  }


}

run "delete_check_has_recent_deploy_outside_of_terraform" {

  variables {
    resource_id   = run.resource_check_has_recent_deploy_create_with_all_fields.this.id
    resource_type = "check"
  }

  module {
    source = "./provisioner"
  }
}

run "resource_check_has_recent_deploy_create_with_required_fields" {

  variables {
    # other fields from file scoped variables block
    category  = run.from_data_module.first_rubric_category.id
    enable_on = null
    enabled   = null
    filter    = null
    level     = run.from_data_module.max_index_rubric_level.id
    notes     = null
    owner     = null
  }

  assert {
    condition = run.resource_check_has_recent_deploy_create_with_all_fields.this.id != opslevel_check_has_recent_deploy.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_check_has_recent_deploy_create_with_all_fields.this.id,
      opslevel_check_has_recent_deploy.this.id,
    )
  }

  module {
    source = "./opslevel_modules/modules/check/has_recent_deploy"
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_has_recent_deploy.this.category,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.days == var.days
    error_message = format(
      "expected '%v' but got '%v'",
      var.days,
      opslevel_check_has_recent_deploy.this.days,
    )
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_has_recent_deploy.this.level,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_has_recent_deploy.this.name,
    )
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.this.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_has_recent_deploy_set_all_fields" {

  variables {
    # other fields from file scoped variables block
    category = run.from_data_module.first_rubric_category.id
    filter   = run.from_data_module.first_filter.id
    level    = run.from_data_module.max_index_rubric_level.id
    owner    = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/has_recent_deploy"
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_has_recent_deploy.this.category,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.days == var.days
    error_message = format(
      "expected '%v' but got '%v'",
      var.days,
      opslevel_check_has_recent_deploy.this.days,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_has_recent_deploy.this.enabled,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_has_recent_deploy.this.enable_on,
    )
  }

  assert {
    condition     = startswith(opslevel_check_has_recent_deploy.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_has_recent_deploy)
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_has_recent_deploy.this.filter,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_has_recent_deploy.this.level,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_has_recent_deploy.this.name,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_has_recent_deploy.this.notes,
    )
  }

  assert {
    condition = opslevel_check_has_recent_deploy.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_has_recent_deploy.this.owner,
    )
  }


}
