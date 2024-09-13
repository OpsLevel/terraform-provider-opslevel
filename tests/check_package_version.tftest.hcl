variables {
  resource_name = "opslevel_check_package_version"

  # -- check_package_version fields --
  # required fields
  missing_package_result = "passed"
  package_constraint     = "matches_version"
  package_manager        = "docker"
  package_name           = "foobar"

  # optional fields
  package_name_is_regex = null
  version_constraint_predicate = {
    type  = "does_not_match_regex"
    value = "^$"
  }

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check package_version"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on package_version check"
  owner     = null
}

run "from_filter_module" {
  command = plan

  module {
    source = "./data/filter"
  }
}

run "from_rubric_category_module" {
  command = plan

  module {
    source = "./opslevel_modules/modules/rubric_category"
  }
}

run "from_rubric_level_module" {
  command = plan

  module {
    source = "./opslevel_modules/modules/rubric_level"
  }
}

run "from_team_module" {
  command = plan

  module {
    source = "./data/team"
  }
}

run "resource_check_package_version_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.first.id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    missing_package_result       = var.missing_package_result
    name                         = var.name
    notes                        = var.notes
    owner                        = run.from_team_module.first.id
    package_constraint           = var.package_constraint
    package_manager              = var.package_manager
    package_name                 = var.package_name
    package_name_is_regex        = var.package_name_is_regex
    version_constraint_predicate = var.version_constraint_predicate
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
