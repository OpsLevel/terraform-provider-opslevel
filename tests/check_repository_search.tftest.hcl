variables {
  check_repository_search = "opslevel_check_repository_search"

  # -- check_repository_search fields --
  # required fields
  file_contents_predicate = {
    type  = "does_not_contain",
    value = "something_unlikely.txt",
  }

  # optional fields
  file_extensions = toset(["go", "py", "rs"])

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Repository Search"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Repository Search Check"
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

run "resource_check_repository_search_create_with_all_fields" {

  variables {
    category                = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on               = var.enable_on
    enabled                 = var.enabled
    file_contents_predicate = var.file_contents_predicate
    file_extensions         = var.file_extensions
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
    source = "./opslevel_modules/modules/check/repository_search"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_repository_search.this.category),
      can(opslevel_check_repository_search.this.description),
      can(opslevel_check_repository_search.this.enable_on),
      can(opslevel_check_repository_search.this.enabled),
      can(opslevel_check_repository_search.this.file_contents_predicate),
      can(opslevel_check_repository_search.this.file_extensions),
      can(opslevel_check_repository_search.this.filter),
      can(opslevel_check_repository_search.this.id),
      can(opslevel_check_repository_search.this.level),
      can(opslevel_check_repository_search.this.name),
      can(opslevel_check_repository_search.this.notes),
      can(opslevel_check_repository_search.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_repository_search)
  }

  assert {
    condition     = opslevel_check_repository_search.this.category == var.category
    error_message = "wrong category of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_repository_search resource"
  }

  assert {
    condition     = startswith(opslevel_check_repository_search.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_repository_search)
  }

  assert {
    condition     = opslevel_check_repository_search.this.file_extensions == var.file_extensions
    error_message = "wrong file_extensions of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.level == var.level
    error_message = "wrong level ID of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_repository_search)
  }

  assert {
    condition     = opslevel_check_repository_search.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_repository_search resource"
  }

}

run "resource_check_repository_search_update_unset_optional_fields" {

  variables {
    category        = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on       = null
    enabled         = null
    filter          = null
    file_extensions = null
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    notes = null
    owner = null
  }

  module {
    source = "./opslevel_modules/modules/check/repository_search"
  }

  assert {
    condition     = opslevel_check_repository_search.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_search.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.file_extensions == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_search.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_search.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_repository_search.this.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_repository_search_update_all_fields" {

  variables {
    category                = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on               = var.enable_on
    enabled                 = var.enabled
    file_contents_predicate = var.file_contents_predicate
    file_extensions         = setunion(var.file_extensions, ["yaml"])
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
    source = "./opslevel_modules/modules/check/repository_search"
  }

  assert {
    condition     = opslevel_check_repository_search.this.category == var.category
    error_message = "wrong category of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.file_extensions == var.file_extensions
    error_message = "wrong file_extensions of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.level == var.level
    error_message = "wrong level ID of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_repository_search)
  }

  assert {
    condition     = opslevel_check_repository_search.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_repository_search resource"
  }

  assert {
    condition     = opslevel_check_repository_search.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_repository_search resource"
  }

}
