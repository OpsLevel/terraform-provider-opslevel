variables {
  check_tool_usage = "opslevel_check_tool_usage"

  # -- check_tool_usage fields --
  # required fields
  tool_category = "api_documentation"

  # optional fields
  environment_predicate = {
    type  = "exists"
    value = null
  }
  tool_name_predicate = {
    type  = "exists"
    value = null
  }
  tool_url_predicate = {
    type  = "exists"
    value = null
  }

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Tool Usage"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Tool Usage Check"
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

run "resource_check_tool_usage_create_with_all_fields" {

  variables {
    category              = run.from_rubric_category_get_category_id.first_category.id
    enable_on             = var.enable_on
    enabled               = var.enabled
    environment_predicate = var.environment_predicate
    filter                = run.from_filter_get_filter_id.first_filter.id
    level                 = run.from_rubric_level_get_level_id.greatest_level.id
    name                  = var.name
    notes                 = var.notes
    owner                 = run.from_team_get_owner_id.first_team.id
    tool_category         = var.tool_category
    tool_name_predicate   = var.tool_name_predicate
    tool_url_predicate    = var.tool_url_predicate
  }

  module {
    source = "./check_tool_usage"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_tool_usage.test.category),
      can(opslevel_check_tool_usage.test.description),
      can(opslevel_check_tool_usage.test.enable_on),
      can(opslevel_check_tool_usage.test.enabled),
      can(opslevel_check_tool_usage.test.filter),
      can(opslevel_check_tool_usage.test.id),
      can(opslevel_check_tool_usage.test.level),
      can(opslevel_check_tool_usage.test.name),
      can(opslevel_check_tool_usage.test.notes),
      can(opslevel_check_tool_usage.test.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_tool_usage)
  }

  assert {
    condition     = opslevel_check_tool_usage.test.category == var.category
    error_message = "wrong category of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = startswith(opslevel_check_tool_usage.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_tool_usage)
  }

  assert {
    condition     = opslevel_check_tool_usage.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.level == var.level
    error_message = "wrong level ID of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_tool_usage)
  }

  assert {
    condition     = opslevel_check_tool_usage.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_tool_usage resource"
  }

}

run "resource_check_tool_usage_update_unset_optional_fields" {

  variables {
    category  = run.from_rubric_category_get_category_id.first_category.id
    enable_on = null
    enabled   = null
    filter    = null
    level     = run.from_rubric_level_get_level_id.greatest_level.id
    notes     = null
    owner     = null
  }

  module {
    source = "./check_tool_usage"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_tool_usage.test.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_tool_usage.test.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_tool_usage.test.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_tool_usage_update_all_fields" {

  variables {
    category              = run.from_rubric_category_get_category_id.first_category.id
    enable_on             = var.enable_on
    enabled               = var.enabled
    environment_predicate = var.environment_predicate
    filter                = run.from_filter_get_filter_id.first_filter.id
    level                 = run.from_rubric_level_get_level_id.greatest_level.id
    name                  = var.name
    notes                 = var.notes
    owner                 = run.from_team_get_owner_id.first_team.id
    tool_category         = var.tool_category
    tool_name_predicate   = var.tool_name_predicate
    tool_url_predicate    = var.tool_url_predicate
  }

  module {
    source = "./check_tool_usage"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.category == var.category
    error_message = "wrong category of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.level == var.level
    error_message = "wrong level ID of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_tool_usage)
  }

  assert {
    condition     = opslevel_check_tool_usage.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_tool_usage resource"
  }

  assert {
    condition     = opslevel_check_tool_usage.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_tool_usage resource"
  }

}
