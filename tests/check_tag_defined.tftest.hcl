variables {
  check_tag_defined = "opslevel_check_tag_defined"

  # -- check_tag_defined fields --
  # required fields
  tag_key = "test-tag-key"

  # optional fields
  tag_predicate = {
    type  = "exists"
    value = null
  }

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Tag Defined"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Tag Defined Check"
  owner     = null
}

run "from_filter_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }
}

run "from_rubric_category_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./rubric_category"
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

  variables {
    name = ""
  }

  module {
    source = "./opslevel_modules/modules/team"
  }
}

run "resource_check_tag_defined_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name          = var.name
    notes         = var.notes
    owner         = run.from_team_module.all.teams[0].id
    tag_predicate = var.tag_predicate
  }

  module {
    source = "./opslevel_modules/modules/check/tag_defined"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_tag_defined.test.category),
      can(opslevel_check_tag_defined.test.description),
      can(opslevel_check_tag_defined.test.enable_on),
      can(opslevel_check_tag_defined.test.enabled),
      can(opslevel_check_tag_defined.test.filter),
      can(opslevel_check_tag_defined.test.id),
      can(opslevel_check_tag_defined.test.level),
      can(opslevel_check_tag_defined.test.name),
      can(opslevel_check_tag_defined.test.notes),
      can(opslevel_check_tag_defined.test.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_tag_defined)
  }

  assert {
    condition     = opslevel_check_tag_defined.test.category == var.category
    error_message = "wrong category of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = startswith(opslevel_check_tag_defined.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_tag_defined)
  }

  assert {
    condition     = opslevel_check_tag_defined.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.level == var.level
    error_message = "wrong level ID of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_tag_defined)
  }

  assert {
    condition     = opslevel_check_tag_defined.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_tag_defined resource"
  }

}

run "resource_check_tag_defined_update_unset_optional_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = null
    enabled   = null
    filter    = null
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    notes         = null
    owner         = null
    tag_predicate = null
  }

  module {
    source = "./opslevel_modules/modules/check/tag_defined"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_tag_defined.test.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_tag_defined.test.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_tag_defined.test.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_tag_defined.test.tag_predicate == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_tag_defined_update_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name          = var.name
    notes         = var.notes
    owner         = run.from_team_module.all.teams[0].id
    tag_predicate = var.tag_predicate
  }

  module {
    source = "./opslevel_modules/modules/check/tag_defined"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.category == var.category
    error_message = "wrong category of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.level == var.level
    error_message = "wrong level ID of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_tag_defined)
  }

  assert {
    condition     = opslevel_check_tag_defined.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_tag_defined resource"
  }

  assert {
    condition     = opslevel_check_tag_defined.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_tag_defined resource"
  }

}
