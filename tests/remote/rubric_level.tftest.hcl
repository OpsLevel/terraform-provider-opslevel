variables {
  rubric_level_one  = "opslevel_rubric_level"
  rubric_levels_all = "opslevel_rubric_levels"

  # required fields
  name = "TF Rubric Level"

  # optional fields
  description = "TF Rubric Level description"
  index       = null # be careful not to overwrite existing Rubric Level
}

run "resource_rubric_level_create_with_all_fields" {

  variables {
    description = var.description
    name        = var.name
    index       = null # only test index via "terraform plan"
  }

  module {
    source = "./rubric_level"
  }

  assert {
    condition = alltrue([
      can(opslevel_rubric_level.test.description),
      can(opslevel_rubric_level.test.id),
      can(opslevel_rubric_level.test.index),
      can(opslevel_rubric_level.test.name),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.rubric_level_one)
  }

  assert {
    condition     = opslevel_rubric_level.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.rubric_level_one)
  }

  assert {
    condition     = startswith(opslevel_rubric_level.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.rubric_level_one)
  }

  assert {
    condition     = opslevel_rubric_level.test.index >= 0
    error_message = "wrong index for ${var.rubric_level_one}"
  }

  assert {
    condition     = opslevel_rubric_level.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.rubric_level_one)
  }

}

run "resource_rubric_level_create_with_empty_optional_fields" {

  variables {
    description             = ""
    name                    = "New ${var.name} with empty fields"
  }

  module {
    source = "./rubric_level"
  }

  assert {
    condition     = opslevel_rubric_level.test.description == ""
    error_message = var.error_expected_empty_string
  }

}

run "resource_rubric_level_update_unset_optional_fields" {

  variables {
    description = null
  }

  module {
    source = "./rubric_level"
  }

  assert {
    condition     = opslevel_rubric_level.test.description == null
    error_message = var.error_expected_null_field
  }

}

run "resource_rubric_level_update_all_fields" {

  variables {
    description = "${var.description} updated"
    name        = "${var.name} updated"
    index       = null # only test index via "terraform plan"
  }

  module {
    source = "./rubric_level"
  }

  assert {
    condition     = opslevel_rubric_level.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.rubric_level_one)
  }

  assert {
    condition     = opslevel_rubric_level.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.rubric_level_one)
  }

}

run "datasource_rubric_levels_all" {

  module {
    source = "./rubric_level"
  }

  assert {
    condition     = length(data.opslevel_rubric_levels.all.rubric_levels) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.rubric_levels_all)
  }

  assert {
    condition = alltrue([
      can(data.opslevel_rubric_levels.all.rubric_levels[0].alias),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].id),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].index),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.rubric_levels_all)
  }

}

run "datasource_rubric_level_first" {

  module {
    source = "./rubric_level"
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.alias == data.opslevel_rubric_levels.all.rubric_levels[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.rubric_level_one)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.id == data.opslevel_rubric_levels.all.rubric_levels[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.rubric_level_one)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.index == data.opslevel_rubric_levels.all.rubric_levels[0].index
    error_message = replace(var.error_wrong_index, "TYPE", var.rubric_level_one)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.name == data.opslevel_rubric_levels.all.rubric_levels[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.rubric_level_one)
  }

}
