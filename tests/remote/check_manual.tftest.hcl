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

run "resource_check_manual_create_with_all_fields" {

  variables {
    category                = run.from_rubric_category_get_category_id.first_category.id
    enable_on               = var.enable_on
    enabled                 = var.enabled
    filter                  = run.from_filter_get_filter_id.first_filter.id
    level                   = run.from_rubric_level_get_level_id.greatest_level.id
    name                    = var.name
    notes                   = var.notes
    owner                   = run.from_team_get_owner_id.first_team.id
    update_requires_comment = var.update_requires_comment
    update_frequency        = var.update_frequency
  }

  module {
    source = "./check_manual"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_manual.test.category),
      can(opslevel_check_manual.test.description),
      can(opslevel_check_manual.test.enable_on),
      can(opslevel_check_manual.test.enabled),
      can(opslevel_check_manual.test.filter),
      can(opslevel_check_manual.test.id),
      can(opslevel_check_manual.test.last_updated),
      can(opslevel_check_manual.test.level),
      can(opslevel_check_manual.test.name),
      can(opslevel_check_manual.test.notes),
      can(opslevel_check_manual.test.owner),
      can(opslevel_check_manual.test.update_frequency),
      can(opslevel_check_manual.test.update_requires_comment),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.test.category == var.category
    error_message = "wrong category of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_manual resource"
  }

  assert {
    condition     = startswith(opslevel_check_manual.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.level == var.level
    error_message = "wrong level ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.update_frequency == var.update_frequency
    error_message = "wrong update_frequency of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.update_requires_comment == var.update_requires_comment
    error_message = "wrong update_requires_comment of opslevel_check_manual resource"
  }

}

run "resource_check_manual_update_unset_optional_fields" {

  variables {
    category         = run.from_rubric_category_get_category_id.first_category.id
    enable_on        = null
    enabled          = !var.enabled
    filter           = null
    level            = run.from_rubric_level_get_level_id.greatest_level.id
    notes            = null
    owner            = null
    update_frequency = null
  }

  module {
    source = "./check_manual"
  }

  assert {
    condition     = opslevel_check_manual.test.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.test.enabled == var.enabled
    error_message = "expected 'enabled' to be unchanged from create, field included in 'ignore_changes' lifecycle"
  }

  assert {
    condition     = opslevel_check_manual.test.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.test.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.test.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_manual.test.update_frequency == null
    error_message = var.error_expected_null_field
  }

}

run "resource_check_manual_update_all_fields" {

  variables {
    category                = run.from_rubric_category_get_category_id.first_category.id
    enable_on               = var.enable_on
    enabled                 = !var.enabled
    filter                  = run.from_filter_get_filter_id.first_filter.id
    level                   = run.from_rubric_level_get_level_id.greatest_level.id
    name                    = "${var.name} updated"
    notes                   = "${var.notes} updated"
    owner                   = run.from_team_get_owner_id.first_team.id
    update_requires_comment = !var.update_requires_comment
    update_frequency        = var.update_frequency
  }

  module {
    source = "./check_manual"
  }

  assert {
    condition     = opslevel_check_manual.test.category == var.category
    error_message = "wrong category of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.enabled == !var.enabled
    error_message = "wrong enabled of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.level == var.level
    error_message = "wrong level ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_manual)
  }

  assert {
    condition     = opslevel_check_manual.test.notes == var.notes
    error_message = "wrong notes of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.update_frequency == var.update_frequency
    error_message = "wrong update_frequency of opslevel_check_manual resource"
  }

  assert {
    condition     = opslevel_check_manual.test.update_requires_comment == var.update_requires_comment
    error_message = "wrong update_requires_comment of opslevel_check_manual resource"
  }

}
