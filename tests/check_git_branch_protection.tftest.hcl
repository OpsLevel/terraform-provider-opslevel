variables {
  check_git_branch_protection = "opslevel_check_git_branch_protection"

  # -- check_git_branch_protection fields --
  # only uses base fields

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Git Branch Protection"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Git Branch Protection Check"
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

run "resource_check_git_branch_protection_create_with_all_fields" {

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
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_git_branch_protection)
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.category == var.category
    error_message = "wrong category of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = startswith(opslevel_check_git_branch_protection.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_git_branch_protection)
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.level == var.level
    error_message = "wrong level ID of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_git_branch_protection)
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_git_branch_protection resource"
  }

}

run "resource_check_git_branch_protection_update_unset_optional_fields" {

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

run "resource_check_git_branch_protection_update_all_fields" {

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
    source = "./opslevel_modules/modules/check/git_branch_protection"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.category == var.category
    error_message = "wrong category of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.level == var.level
    error_message = "wrong level ID of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_git_branch_protection)
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_git_branch_protection resource"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_git_branch_protection resource"
  }

}

