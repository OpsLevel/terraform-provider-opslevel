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

  module {
    source = "./data/filter"
  }
}

run "from_rubric_category_module" {
  command = plan

  module {
    source = "./opslevel_modules/modules/rubric_category/data"
  }
}

run "from_rubric_level_module" {
  command = plan

  module {
    source = "./opslevel_modules/modules/rubric_level/data"
  }
}

run "from_team_module" {
  command = plan

  module {
    source = "./data/team"
  }
}

run "resource_check_repository_file_create_with_all_fields" {

  variables {
    category                = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on               = var.enable_on
    enabled                 = var.enabled
    file_contents_predicate = var.file_contents_predicate
    filter                  = run.from_filter_module.first.id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.first.id
  }

  module {
    source = "./opslevel_modules/modules/check/repository_file"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_repository_file.this.category),
      can(opslevel_check_repository_file.this.description),
      can(opslevel_check_repository_file.this.enable_on),
      can(opslevel_check_repository_file.this.enabled),
      can(opslevel_check_repository_file.this.file_contents_predicate),
      can(opslevel_check_repository_file.this.filter),
      can(opslevel_check_repository_file.this.id),
      can(opslevel_check_repository_file.this.level),
      can(opslevel_check_repository_file.this.name),
      can(opslevel_check_repository_file.this.notes),
      can(opslevel_check_repository_file.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.this.category == var.category
    error_message = "wrong category of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_repository_file resource"
  }

  assert {
    condition     = startswith(opslevel_check_repository_file.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.this.file_contents_predicate == var.file_contents_predicate
    error_message = "wrong file_contents_predicate of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.level == var.level
    error_message = "wrong level ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.owner == var.owner
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
    condition     = opslevel_check_repository_file.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.file_contents_predicate == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_file.this.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_repository_file_update_all_fields" {

  variables {
    category                = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on               = var.enable_on
    enabled                 = var.enabled
    file_contents_predicate = var.file_contents_predicate
    filter                  = run.from_filter_module.first.id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.first.id
  }

  module {
    source = "./opslevel_modules/modules/check/repository_file"
  }

  assert {
    condition     = opslevel_check_repository_file.this.category == var.category
    error_message = "wrong category of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.file_contents_predicate == var.file_contents_predicate
    error_message = "wrong file_contents_predicate of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.level == var.level
    error_message = "wrong level ID of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_repository_file)
  }

  assert {
    condition     = opslevel_check_repository_file.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_repository_file resource"
  }

  assert {
    condition     = opslevel_check_repository_file.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_repository_file resource"
  }

}
