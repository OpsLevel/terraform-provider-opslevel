run "datasource_lifecycles_all" {

  assert {
    condition     = length(data.opslevel_lifecycles.all.lifecycles) > 0
    error_message = "zero lifecycles found in data.opslevel_lifecycles"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_lifecycles.all.lifecycles[0].alias),
      can(data.opslevel_lifecycles.all.lifecycles[0].id),
      can(data.opslevel_lifecycles.all.lifecycles[0].index),
      can(data.opslevel_lifecycles.all.lifecycles[0].name),
    ])
    error_message = "cannot set all expected lifecycle datasource fields"
  }

}

run "datasource_lifecycle_first" {

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.alias == data.opslevel_lifecycles.all.lifecycles[0].alias
    error_message = "wrong alias on opslevel_lifecycle"
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.id == data.opslevel_lifecycles.all.lifecycles[0].id
    error_message = "wrong ID on opslevel_lifecycle"
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.index == data.opslevel_lifecycles.all.lifecycles[0].index
    error_message = "wrong index on opslevel_lifecycle"
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_name.name == data.opslevel_lifecycles.all.lifecycles[0].name
    error_message = "wrong name on opslevel_lifecycle"
  }

}
