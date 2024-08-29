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

run "from_filter_get_filter_id" {
  command = plan

  variables {
    connective = null
  }

  module {
    source = "./filter"
  }
}

run "from_rubric_category_get_category_id" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./rubric_category"
  }
}

run "from_rubric_level_get_level_id" {
  command = plan

  variables {
    description = null
    index       = null
    name        = ""
  }

  module {
    source = "./rubric_level"
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

run "resource_check_package_version_create_with_all_fields" {

  variables {
    category                     = run.from_rubric_category_get_category_id.first_category.id
    enable_on                    = var.enable_on
    enabled                      = var.enabled
    filter                       = run.from_filter_get_filter_id.first_filter.id
    level                        = run.from_rubric_level_get_level_id.greatest_level.id
    missing_package_result       = var.missing_package_result
    name                         = var.name
    notes                        = var.notes
    owner                        = run.from_team_get_owner_id.first_team.id
    package_constraint           = var.package_constraint
    package_manager              = var.package_manager
    package_name                 = var.package_name
    package_name_is_regex        = var.package_name_is_regex
    version_constraint_predicate = var.version_constraint_predicate
  }

  module {
    source = "./check_package_version"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_package_version.test.category),
      can(opslevel_check_package_version.test.description),
      can(opslevel_check_package_version.test.enable_on),
      can(opslevel_check_package_version.test.enabled),
      can(opslevel_check_package_version.test.filter),
      can(opslevel_check_package_version.test.id),
      can(opslevel_check_package_version.test.level),
      can(opslevel_check_package_version.test.missing_package_result),
      can(opslevel_check_package_version.test.name),
      can(opslevel_check_package_version.test.notes),
      can(opslevel_check_package_version.test.owner),
      can(opslevel_check_package_version.test.package_constraint),
      can(opslevel_check_package_version.test.package_manager),
      can(opslevel_check_package_version.test.package_name),
      can(opslevel_check_package_version.test.package_name_is_regex),
      can(opslevel_check_package_version.test.version_constraint_predicate),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_package_version.test.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_check_package_version.test.category,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.enable_on == var.enable_on
    error_message = format(
      "expected '%v' but got '%v'",
      var.enable_on,
      opslevel_check_package_version.test.enable_on,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.enabled == var.enabled
    error_message = format(
      "expected '%v' but got '%v'",
      var.enabled,
      opslevel_check_package_version.test.enabled,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_check_package_version.test.filter,
    )
  }

  assert {
    condition     = startswith(opslevel_check_package_version.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_check_package_version.test.level == var.level
    error_message = format(
      "expected '%v' but got '%v'",
      var.level,
      opslevel_check_package_version.test.level,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.missing_package_result == var.missing_package_result
    error_message = format(
      "expected '%v' but got '%v'",
      var.missing_package_result,
      opslevel_check_package_version.test.missing_package_result,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_check_package_version.test.name,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.notes == var.notes
    error_message = format(
      "expected '%v' but got '%v'",
      var.notes,
      opslevel_check_package_version.test.notes,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_check_package_version.test.owner,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.package_constraint == var.package_constraint
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_constraint,
      opslevel_check_package_version.test.package_constraint,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.package_manager == var.package_manager
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_manager,
      opslevel_check_package_version.test.package_manager,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.package_name == var.package_name
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name,
      opslevel_check_package_version.test.package_name,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.package_name_is_regex == var.package_name_is_regex
    error_message = format(
      "expected '%v' but got '%v'",
      var.package_name_is_regex,
      opslevel_check_package_version.test.package_name_is_regex,
    )
  }

  assert {
    condition = opslevel_check_package_version.test.version_constraint_predicate == var.version_constraint_predicate
    error_message = format(
      "expected '%v' but got '%v'",
      var.version_constraint_predicate,
      opslevel_check_package_version.test.version_constraint_predicate,
    )
  }

}
