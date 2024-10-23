variables {
  resource_name = "opslevel_check_code_issue"

  # -- check_code_issue fields --
  # required fields
  constraint = "any"

  # optional fields
  issue_name  = "idk"
  issue_type  = ["snyk:code"]
  max_allowed = 5
  resolution_time = {
    unit  = "week"
    value = 3
  }
  severity = ["snyk:high"]

  # -- check base fields --
  # required fields
  category = null # sourced from module
  level    = null # sourced from module
  name     = "TF Test Check Code Issue"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null # sourced from module
  notes     = "Notes on Code Issue Check"
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

run "resource_check_code_issue_create_with_constraint_any" {

  variables {
    # other fields from file scoped variables block
    category   = run.from_data_module.first_rubric_category.id
    constraint = "any"
    issue_name = null # not allowed when constraint is "any"
    filter     = run.from_data_module.first_filter.id
    level      = run.from_data_module.max_index_rubric_level.id
    owner      = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/check/code_issue"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_code_issue.this.category),
      can(opslevel_check_code_issue.this.constraint),
      can(opslevel_check_code_issue.this.description),
      can(opslevel_check_code_issue.this.enable_on),
      can(opslevel_check_code_issue.this.enabled),
      can(opslevel_check_code_issue.this.filter),
      can(opslevel_check_code_issue.this.id),
      can(opslevel_check_code_issue.this.issue_name),
      can(opslevel_check_code_issue.this.issue_type),
      can(opslevel_check_code_issue.this.level),
      can(opslevel_check_code_issue.this.max_allowed),
      can(opslevel_check_code_issue.this.name),
      can(opslevel_check_code_issue.this.notes),
      can(opslevel_check_code_issue.this.owner),
      can(opslevel_check_code_issue.this.resolution_time),
      can(opslevel_check_code_issue.this.severity),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }


  assert {
    condition = opslevel_check_code_issue.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_code_issue.this.category,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.constraint == var.constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.constraint,
      opslevel_check_code_issue.this.constraint,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_code_issue.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_code_issue.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_code_issue.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_code_issue.this.issue_name == var.issue_name
    error_message = format(
      "expected '%v' but got '%v'",
      var.issue_name,
      opslevel_check_code_issue.this.issue_name,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.issue_type == var.issue_type
    error_message = format(
      "expected '%v' but got '%v'",
      var.issue_type,
      opslevel_check_code_issue.this.issue_type,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_code_issue.this.filter,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_code_issue.this.level,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.max_allowed == var.max_allowed
    error_message = format(
      "expected '%v' but got '%v'",
      var.max_allowed,
      opslevel_check_code_issue.this.max_allowed,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_code_issue.this.name,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_code_issue.this.notes,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_code_issue.this.owner,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.resolution_time == var.resolution_time
    error_message = format(
      "expected '%v' but got '%v'",
      var.resolution_time,
      opslevel_check_code_issue.this.resolution_time,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.severity == var.severity
    error_message = format(
      "expected '%v' but got '%v'",
      var.severity,
      opslevel_check_code_issue.this.severity,
    )
  }

}

run "resource_check_code_issue_unset_optional_fields" {

  variables {
    # other fields from file scoped variables block
    category    = run.from_data_module.first_rubric_category.id
    constraint  = "any"
    enable_on   = null
    enabled     = null
    issue_name  = null # not allowed when constraint is "any"
    issue_type  = null
    filter      = null
    level       = run.from_data_module.max_index_rubric_level.id
    max_allowed = null
    notes       = null
    owner       = null
    # resolution_time = null
  }

  module {
    source = "./opslevel_modules/modules/check/code_issue"
  }

  assert {
    condition = opslevel_check_code_issue.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_code_issue.this.category,
    )
  }

  assert {
    condition     = opslevel_check_code_issue.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_code_issue.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_code_issue resource"
  }

  assert {
    condition     = opslevel_check_code_issue.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_code_issue.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_code_issue.this.level,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_code_issue.this.name,
    )
  }

  assert {
    condition     = opslevel_check_code_issue.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_code_issue.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "delete_check_code_issue_with_constraint_any_outside_of_terraform" {

  plan_options {
    target = [
      terraform_data.opslevel_cli,
    ]
  }

  variables {
    command = "delete check ${run.resource_check_code_issue_create_with_constraint_any.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_check_code_issue_create_with_constraint_contains" {

  variables {
    # other fields from file scoped variables block
    category   = run.from_data_module.first_rubric_category.id
    constraint = "contains"
    enable_on  = null
    enabled    = null
    filter     = null
    issue_name = "w"  # required when constraint is "contains"
    issue_type = null # not allowed when constraint is "contains"
    level      = run.from_data_module.max_index_rubric_level.id
    notes      = null
    owner      = null
    severity   = null # not allowed when constraint is "contains"
  }

  module {
    source = "./opslevel_modules/modules/check/code_issue"
  }

  # assert {
  #   condition = run.resource_check_code_issue_create_with_constraint_any.this.id != opslevel_check_code_issue.this.id
  #   error_message = format(
  #     "expected old id '%v' to be different from new id '%v'",
  #     run.resource_check_code_issue_create_with_constraint_any.this.id,
  #     opslevel_check_code_issue.this.id,
  #   )
  # }

  assert {
    condition = opslevel_check_code_issue.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_code_issue.this.category,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.constraint == var.constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.constraint,
      opslevel_check_code_issue.this.constraint,
    )
  }

  assert {
    condition     = opslevel_check_code_issue.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_code_issue.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_code_issue resource"
  }

  assert {
    condition     = opslevel_check_code_issue.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_check_code_issue.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_code_issue.this.level,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_code_issue.this.name,
    )
  }

  assert {
    condition     = opslevel_check_code_issue.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_code_issue.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "delete_check_code_issue_with_constraint_contains_outside_of_terraform" {

  plan_options {
    target = [
      terraform_data.delete_command,
    ]

  }
  variables {
    resource_id   = run.resource_check_code_issue_create_with_constraint_contains.this.id
    resource_type = "check"
  }

  module {
    source = "./cli"
  }
}

run "resource_check_code_issue_create_with_constraint_exact" {

  variables {
    # other fields from file scoped variables block
    category    = run.from_data_module.first_rubric_category.id
    constraint  = "exact"
    filter      = run.from_data_module.first_filter.id
    issue_type  = null # not allowed when constraint is "exact"
    level       = run.from_data_module.max_index_rubric_level.id
    max_allowed = null # not allowed when constraint is "exact"
    owner       = run.from_data_module.first_team.id
    severity    = null # not allowed when constraint is "exact"
  }

  module {
    source = "./opslevel_modules/modules/check/code_issue"
  }

  assert {
    condition = run.resource_check_code_issue_create_with_constraint_contains.this.id != opslevel_check_code_issue.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_check_code_issue_create_with_constraint_contains.this.id,
      opslevel_check_code_issue.this.id,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_code_issue.this.category,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.constraint == var.constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.constraint,
      opslevel_check_code_issue.this.constraint,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_code_issue.this.enable_on,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_code_issue.this.enabled,
    )
  }

  assert {
    condition     = startswith(opslevel_check_code_issue.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_code_issue.this.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_code_issue.this.filter,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_code_issue.this.level,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_code_issue.this.name,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_code_issue.this.notes,
    )
  }

  assert {
    condition = opslevel_check_code_issue.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_code_issue.this.owner,
    )
  }

}
