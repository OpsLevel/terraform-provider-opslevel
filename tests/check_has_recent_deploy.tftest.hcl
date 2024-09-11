variables {
  check_has_recent_deploy = "opslevel_check_has_recent_deploy"

  # -- check_has_recent_deploy fields --
  # required fields
  days = 5

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Has Recent Deploy"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Has Recent Deploy Check"
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

run "resource_check_has_recent_deploy_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/has_recent_deploy"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_has_recent_deploy.test.category),
      can(opslevel_check_has_recent_deploy.test.description),
      can(opslevel_check_has_recent_deploy.test.enable_on),
      can(opslevel_check_has_recent_deploy.test.enabled),
      can(opslevel_check_has_recent_deploy.test.filter),
      can(opslevel_check_has_recent_deploy.test.id),
      can(opslevel_check_has_recent_deploy.test.level),
      can(opslevel_check_has_recent_deploy.test.name),
      can(opslevel_check_has_recent_deploy.test.notes),
      can(opslevel_check_has_recent_deploy.test.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_has_recent_deploy)
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.category == var.category
    error_message = "wrong category of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = startswith(opslevel_check_has_recent_deploy.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_has_recent_deploy)
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.level == var.level
    error_message = "wrong level ID of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_has_recent_deploy)
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_has_recent_deploy resource"
  }

}

run "resource_check_has_recent_deploy_update_unset_optional_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = null
    enabled   = null
    filter    = null
    level     = run.from_rubric_level_module.greatest_level.id
    notes     = null
    owner     = null
  }

  module {
    source = "./opslevel_modules/modules/check/has_recent_deploy"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_has_recent_deploy_update_all_fields" {

  variables {
    category  = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_module.all.filters[0].id
    level     = run.from_rubric_level_module.greatest_level.id
    name      = var.name
    notes     = var.notes
    owner     = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/has_recent_deploy"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.category == var.category
    error_message = "wrong category of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.level == var.level
    error_message = "wrong level ID of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_has_recent_deploy)
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_has_recent_deploy resource"
  }

  assert {
    condition     = opslevel_check_has_recent_deploy.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_has_recent_deploy resource"
  }

}
