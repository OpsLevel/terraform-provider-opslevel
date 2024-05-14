run "datasource_rubric_levels_all" {

  variables {
    datasource_type = "opslevel_rubric_levels"
  }

  module {
    source = "./rubric_level"
  }

  assert {
    condition     = length(data.opslevel_rubric_levels.all.rubric_levels) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

  assert {
    condition = alltrue([
      can(data.opslevel_rubric_levels.all.rubric_levels[0].alias),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].id),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].index),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

}

run "datasource_rubric_level_first" {

  variables {
    datasource_type = "opslevel_rubric_level"
  }

  module {
    source = "./rubric_level"
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.alias == data.opslevel_rubric_levels.all.rubric_levels[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.id == data.opslevel_rubric_levels.all.rubric_levels[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.index == data.opslevel_rubric_levels.all.rubric_levels[0].index
    error_message = replace(var.error_wrong_index, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.name == data.opslevel_rubric_levels.all.rubric_levels[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.datasource_type)
  }

}
