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

run "resource_check_has_documentation_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.first.id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.first.id
  }

  module {
    source = "./opslevel_modules/modules/check/has_documentation"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_has_documentation.this.category),
      can(opslevel_check_has_documentation.this.description),
      can(opslevel_check_has_documentation.this.enable_on),
      can(opslevel_check_has_documentation.this.enabled),
      can(opslevel_check_has_documentation.this.filter),
      can(opslevel_check_has_documentation.this.id),
      can(opslevel_check_has_documentation.this.level),
      can(opslevel_check_has_documentation.this.name),
      can(opslevel_check_has_documentation.this.notes),
      can(opslevel_check_has_documentation.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.this.category == var.category
    error_message = "wrong category of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = startswith(opslevel_check_has_documentation.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.level == var.level
    error_message = "wrong level ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_has_documentation resource"
  }

}

run "resource_check_has_documentation_update_unset_optional_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = null
    enabled   = null
    filter    = null
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    notes = null
    owner = null
  }

  module {
    source = "./opslevel_modules/modules/check/has_documentation"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_documentation.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_documentation.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_documentation.this.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_has_documentation_update_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.first.id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.first.id
  }

  module {
    source = "./opslevel_modules/modules/check/has_documentation"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.category == var.category
    error_message = "wrong category of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.level == var.level
    error_message = "wrong level ID of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_has_documentation)
  }

  assert {
    condition     = opslevel_check_has_documentation.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_has_documentation resource"
  }

  assert {
    condition     = opslevel_check_has_documentation.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_has_documentation resource"
  }

}
