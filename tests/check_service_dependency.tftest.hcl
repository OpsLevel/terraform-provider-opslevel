variables {
  check_service_dependency = "opslevel_check_service_dependency"

  # -- check_service_dependency fields --
  # only uses base fields

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Service Dependency"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Service Dependency Check"
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
    name             = ""
  }

  module {
    source = "./opslevel_modules/modules/team"
  }
}

run "resource_check_service_dependency_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name      = var.name
    notes     = var.notes
    owner = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/service_dependency"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_service_dependency.this.category),
      can(opslevel_check_service_dependency.this.description),
      can(opslevel_check_service_dependency.this.enable_on),
      can(opslevel_check_service_dependency.this.enabled),
      can(opslevel_check_service_dependency.this.filter),
      can(opslevel_check_service_dependency.this.id),
      can(opslevel_check_service_dependency.this.level),
      can(opslevel_check_service_dependency.this.name),
      can(opslevel_check_service_dependency.this.notes),
      can(opslevel_check_service_dependency.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_service_dependency)
  }

  assert {
    condition     = opslevel_check_service_dependency.this.category == var.category
    error_message = "wrong category of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = startswith(opslevel_check_service_dependency.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_service_dependency)
  }

  assert {
    condition     = opslevel_check_service_dependency.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.level == var.level
    error_message = "wrong level ID of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_service_dependency)
  }

  assert {
    condition     = opslevel_check_service_dependency.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_service_dependency resource"
  }

}

run "resource_check_service_dependency_update_unset_optional_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = null
    enabled   = null
    filter    = null
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    notes     = null
    owner     = null
  }

  module {
    source = "./opslevel_modules/modules/check/service_dependency"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_service_dependency.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_service_dependency.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_service_dependency.this.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_service_dependency_update_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name      = var.name
    notes     = var.notes
    owner = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/service_dependency"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.category == var.category
    error_message = "wrong category of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.level == var.level
    error_message = "wrong level ID of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_service_dependency)
  }

  assert {
    condition     = opslevel_check_service_dependency.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_service_dependency resource"
  }

  assert {
    condition     = opslevel_check_service_dependency.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_service_dependency resource"
  }

}
