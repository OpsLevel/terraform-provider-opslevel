run "datasource_rubric_levels_all" {

  variables {
    datasource_type = "opslevel_rubric_levels"
  }

  assert {
    condition     = length(data.opslevel_rubric_levels.all.rubric_levels) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

  assert {
    condition = alltrue([
      can(data.opslevel_rubric_levels.all.rubric_levels[0].alias),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].id),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].index),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].name),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

}

run "datasource_rubric_level_first" {

  variables {
    datasource_type = "opslevel_rubric_level"
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.alias == data.opslevel_rubric_levels.all.rubric_levels[0].alias
    error_message = replace(var.wrong_alias_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.id == data.opslevel_rubric_levels.all.rubric_levels[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.index == data.opslevel_rubric_levels.all.rubric_levels[0].index
    error_message = replace(var.wrong_index_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.name == data.opslevel_rubric_levels.all.rubric_levels[0].name
    error_message = replace(var.wrong_name_error, "TYPE", var.datasource_type)
  }

}
