variables {
  check_service_ownership = "opslevel_check_service_ownership"

  # -- check_service_ownership fields --
  # optional fields
  contact_method         = "any"
  require_contact_method = true
  tag_key                = "test-tag-key"
  tag_predicate = {
    type  = "does_not_match_regex"
    value = "abckjasldkfj"
  }

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Service Ownership"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Service Ownership Check"
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

run "resource_check_service_ownership_create_with_all_fields" {

  variables {
    category       = run.from_rubric_category_module.all.rubric_categories[0].id
    contact_method = var.contact_method
    enable_on      = var.enable_on
    enabled        = var.enabled
    filter         = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name          = var.name
    notes         = var.notes
    owner         = run.from_team_module.first.id
    tag_predicate = var.tag_predicate
  }

  module {
    source = "./opslevel_modules/modules/check/service_ownership"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_service_ownership.this.category),
      can(opslevel_check_service_ownership.this.contact_method),
      can(opslevel_check_service_ownership.this.description),
      can(opslevel_check_service_ownership.this.enable_on),
      can(opslevel_check_service_ownership.this.enabled),
      can(opslevel_check_service_ownership.this.filter),
      can(opslevel_check_service_ownership.this.id),
      can(opslevel_check_service_ownership.this.level),
      can(opslevel_check_service_ownership.this.name),
      can(opslevel_check_service_ownership.this.notes),
      can(opslevel_check_service_ownership.this.owner),
      can(opslevel_check_service_ownership.this.require_contact_method),
      can(opslevel_check_service_ownership.this.tag_key),
      can(opslevel_check_service_ownership.this.tag_predicate),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_service_ownership)
  }

  assert {
    condition     = opslevel_check_service_ownership.this.category == var.category
    error_message = "wrong category of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = lower(opslevel_check_service_ownership.this.contact_method) == lower(var.contact_method)
    error_message = "wrong contact_method of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = startswith(opslevel_check_service_ownership.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_service_ownership)
  }

  assert {
    condition     = opslevel_check_service_ownership.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.level == var.level
    error_message = "wrong level ID of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_service_ownership)
  }

  assert {
    condition     = opslevel_check_service_ownership.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_service_ownership resource"
  }

}

run "resource_check_service_ownership_update_unset_optional_fields" {

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
    source = "./opslevel_modules/modules/check/service_ownership"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_service_ownership.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_service_ownership.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_service_ownership.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_service_ownership.this.tag_predicate == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_service_ownership_update_all_fields" {

  variables {
    category       = run.from_rubric_category_module.all.rubric_categories[0].id
    contact_method = var.contact_method
    enable_on      = var.enable_on
    enabled        = var.enabled
    filter         = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name          = var.name
    notes         = var.notes
    owner         = run.from_team_module.first.id
    tag_predicate = var.tag_predicate
  }

  module {
    source = "./opslevel_modules/modules/check/service_ownership"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.category == var.category
    error_message = "wrong category of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = lower(opslevel_check_service_ownership.this.contact_method) == lower(var.contact_method)
    error_message = "wrong contact_method of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.level == var.level
    error_message = "wrong level ID of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_service_ownership)
  }

  assert {
    condition     = opslevel_check_service_ownership.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_service_ownership resource"
  }

  assert {
    condition     = opslevel_check_service_ownership.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_service_ownership resource"
  }

}
