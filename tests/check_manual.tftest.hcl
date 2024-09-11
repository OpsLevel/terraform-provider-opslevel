variables {
  check_manual = "opslevel_check_manual"

  # -- check_manual fields --
  # required fields
  update_requires_comment = true

  # optional fields
  update_frequency = {
    starting_date = "2020-02-12T06:36:13Z"
    time_scale    = "week"
    value         = 1
  }

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Manual"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on manual check"
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

  variables {
    name = ""
  }

  module {
    source = "./opslevel_modules/modules/team"
  }
}

run "resource_check_manual_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name                    = var.name
    notes                   = var.notes
    owner                   = run.from_team_module.all.teams[0].id
    update_requires_comment = var.update_requires_comment
    update_frequency        = var.update_frequency
  }

  module {
    source = "./opslevel_modules/modules/check/manual"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_manual.this.category),
      can(opslevel_check_manual.this.description),
      can(opslevel_check_manual.this.enable_on),
      can(opslevel_check_manual.this.enabled),
      can(opslevel_check_manual.this.filter),
      can(opslevel_check_manual.this.id),
      can(opslevel_check_manual.this.level),
      can(opslevel_check_manual.this.name),
      can(opslevel_check_manual.this.notes),
      can(opslevel_check_manual.this.owner),
      can(opslevel_check_manual.this.update_frequency),
      can(opslevel_check_manual.this.update_requires_comment),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.this.category == var.category
    error_message = "wrong category of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_manual resource"
  }

  assert {
    condition     = startswith(opslevel_check_manual.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.level == var.level
    error_message = "wrong level ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.update_frequency == var.update_frequency
    error_message = "wrong update_frequency of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.update_requires_comment == var.update_requires_comment
    error_message = "wrong update_requires_comment of opslevel_check_manual resource"
  }

}

run "resource_check_manual_update_unset_optional_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = null
    enabled   = null
    filter    = null
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    notes            = null
    owner            = null
    update_frequency = null
  }

  module {
    source = "./opslevel_modules/modules/check/manual"
  }

  assert {
    condition     = opslevel_check_manual.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.enabled == true
    error_message = "expected 'enabled' to be unchanged from create, field included in 'ignore_changes' lifecycle"
  }

  assert {
    condition     = opslevel_check_manual.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.this.update_frequency == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_manual_update_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = !var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name                    = "${var.name} updated"
    notes                   = "${var.notes} updated"
    owner                   = run.from_team_module.all.teams[0].id
    update_requires_comment = !var.update_requires_comment
    update_frequency        = var.update_frequency
  }

  module {
    source = "./opslevel_modules/modules/check/manual"
  }

  assert {
    condition     = opslevel_check_manual.this.category == var.category
    error_message = "wrong category of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.enabled == !var.enabled
    error_message = "wrong enabled of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.level == var.level
    error_message = "wrong level ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.update_frequency == var.update_frequency
    error_message = "wrong update_frequency of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.this.update_requires_comment == var.update_requires_comment
    error_message = "wrong update_requires_comment of opslevel_check_manual resource"
  }

}
