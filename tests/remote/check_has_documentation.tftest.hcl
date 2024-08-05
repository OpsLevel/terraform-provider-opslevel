variables {
  check_has_documentation = "opslevel_check_has_documentation"

  # -- check_has_documentation fields --
  # required fields
  document_type    = "api"
  document_subtype = "openapi"

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Has Documentation"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Has Documentation Check"
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

run "resource_check_has_documentation_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_get_category_id.first_category.id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_get_filter_id.first_filter.id
    level     = run.from_rubric_level_get_level_id.greatest_level.id
    name      = var.name
    notes     = var.notes
    owner     = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./check_has_documentation"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_has_documentation.test.category),
      can(opslevel_check_has_documentation.test.description),
      can(opslevel_check_has_documentation.test.enable_on),
      can(opslevel_check_has_documentation.test.enabled),
      can(opslevel_check_has_documentation.test.filter),
      can(opslevel_check_has_documentation.test.id),
      can(opslevel_check_has_documentation.test.level),
      can(opslevel_check_has_documentation.test.name),
      can(opslevel_check_has_documentation.test.notes),
      can(opslevel_check_has_documentation.test.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.test.category == var.category
    error_message = "wrong category of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = startswith(opslevel_check_has_documentation.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.level == var.level
    error_message = "wrong level ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_has_documentation resource"
  }

}

run "resource_check_has_documentation_update_unset_optional_fields" {

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
    source = "./check_has_documentation"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_documentation.test.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_documentation.test.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_documentation.test.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_has_documentation_update_all_fields" {

  variables {
    category  = run.from_rubric_category_get_category_id.first_category.id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_get_filter_id.first_filter.id
    level     = run.from_rubric_level_get_level_id.greatest_level.id
    name      = var.name
    notes     = var.notes
    owner     = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./check_has_documentation"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.category == var.category
    error_message = "wrong category of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.level == var.level
    error_message = "wrong level ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_has_documentation resource"
  }

}
