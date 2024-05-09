run "datasource_rubric_levels_all" {

  assert {
    condition     = length(data.opslevel_rubric_levels.all.rubric_levels) > 0
    error_message = "zero rubric_levels found in data.opslevel_rubric_levels"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_rubric_levels.all.rubric_levels[0].alias),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].id),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].index),
      can(data.opslevel_rubric_levels.all.rubric_levels[0].name),
    ])
    error_message = "cannot set all expected rubric_level datasource fields"
  }

}

run "datasource_rubric_level_first" {

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.alias == data.opslevel_rubric_levels.all.rubric_levels[0].alias
    error_message = "wrong alias opslevel_rubric_level"
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.index == data.opslevel_rubric_levels.all.rubric_levels[0].index
    error_message = "wrong index on opslevel_rubric_level"
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.id == data.opslevel_rubric_levels.all.rubric_levels[0].id
    error_message = "wrong ID on opslevel_rubric_level"
  }

  assert {
    condition     = data.opslevel_rubric_level.first_level_by_id.name == data.opslevel_rubric_levels.all.rubric_levels[0].name
    error_message = "wrong name on opslevel_rubric_level"
  }

}
