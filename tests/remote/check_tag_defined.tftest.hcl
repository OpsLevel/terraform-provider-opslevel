variables {
  check_tag_defined = "opslevel_check_tag_defined"

  # -- check_tag_defined fields --
  # required fields
  tag_key = "test-tag-key"

  # optional fields
  tag_predicate = null

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

run "from_filter_get_filter_id" {
  command = plan

  variables {
    connective     = null
    predicate_list = null
  }

  module {
    source = "./filter"
  }
}

run "from_rubric_category_get_category_id" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./rubric_category"
  }
}

run "from_rubric_level_get_level_id" {
  command = plan

  variables {
    description = null
    index       = null
    name        = ""
  }

  module {
    source = "./rubric_level"
  }
}

run "from_team_get_owner_id" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }
}

run "resource_check_tag_defined_create_with_all_fields" {

  variables {
    category  = run.from_rubric_category_get_category_id.first_category.id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_get_filter_id.first_filter.id
    level     = run.from_rubric_level_get_level_id.greatest_level.id
    name      = var.name
    notes     = var.notes
    owner     = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./check_tag_defined"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_tag_defined.test.category),
      can(opslevel_check_tag_defined.test.description),
      can(opslevel_check_tag_defined.test.enable_on),
      can(opslevel_check_tag_defined.test.enabled),
      can(opslevel_check_tag_defined.test.filter),
      can(opslevel_check_tag_defined.test.id),
      can(opslevel_check_tag_defined.test.last_updated),
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
    category  = run.from_rubric_category_get_category_id.first_category.id
    enable_on = null
    enabled   = null
    filter    = null
    level     = run.from_rubric_level_get_level_id.greatest_level.id
    notes     = null
    owner     = null
  }

  module {
    source = "./check_tag_defined"
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


}

run "resource_check_tag_defined_update_all_fields" {

  variables {
    category  = run.from_rubric_category_get_category_id.first_category.id
    enable_on = var.enable_on
    enabled   = var.enabled
    filter    = run.from_filter_get_filter_id.first_filter.id
    level     = run.from_rubric_level_get_level_id.greatest_level.id
    name      = var.name
    notes     = var.notes
    owner     = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./check_tag_defined"
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