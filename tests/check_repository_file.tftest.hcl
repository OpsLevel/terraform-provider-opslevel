variables {
  check_repository_file = "opslevel_check_repository_file"

  # -- check_repository_file fields --
  # required fields
  directory_search  = true
  filepaths         = tolist(["one/two.py"])
  use_absolute_root = true

  # optional fields
  file_contents_predicate = {
    type  = "exists"
    value = null
  }

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check repository_file"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on repository_file check"
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

run "resource_check_repository_file_create_with_all_fields" {

  variables {
    category                = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on               = var.enable_on
    enabled                 = var.enabled
    file_contents_predicate = var.file_contents_predicate
    filter                  = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/repository_file"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_repository_file.test.category),
      can(opslevel_check_repository_file.test.description),
      can(opslevel_check_repository_file.test.enable_on),
      can(opslevel_check_repository_file.test.enabled),
      can(opslevel_check_repository_file.test.file_contents_predicate),
      can(opslevel_check_repository_file.test.filter),
      can(opslevel_check_repository_file.test.id),
      can(opslevel_check_repository_file.test.level),
      can(opslevel_check_repository_file.test.name),
      can(opslevel_check_repository_file.test.notes),
      can(opslevel_check_repository_file.test.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.test.category == var.category
    error_message = "wrong category of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_repository_file resource"
  }

  assert {
    condition     = startswith(opslevel_check_repository_file.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.test.file_contents_predicate == var.file_contents_predicate
    error_message = "wrong file_contents_predicate of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.level == var.level
    error_message = "wrong level ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_repository_file resource"
  }

}

run "resource_check_repository_file_update_unset_optional_fields" {

  variables {
    category                = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on               = null
    enabled                 = null
    file_contents_predicate = null
    filter                  = null
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    notes = null
    owner = null
  }

  module {
    source = "./opslevel_modules/modules/check/repository_file"
  }

  assert {
    condition     = opslevel_check_repository_file.test.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.test.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.file_contents_predicate == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.test.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.test.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.test.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_repository_file_update_all_fields" {

  variables {
    category                = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on               = var.enable_on
    enabled                 = var.enabled
    file_contents_predicate = var.file_contents_predicate
    filter                  = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/repository_file"
  }

  assert {
    condition     = opslevel_check_repository_file.test.category == var.category
    error_message = "wrong category of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.file_contents_predicate == var.file_contents_predicate
    error_message = "wrong file_contents_predicate of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.level == var.level
    error_message = "wrong level ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_repository_file resource"
  }

}
