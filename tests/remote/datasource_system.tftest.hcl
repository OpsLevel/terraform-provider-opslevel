run "datasource_systems_all" {

  assert {
    condition     = length(data.opslevel_systems.all.systems) > 0
    error_message = "zero systems found in data.opslevel_systems"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_systems.all.systems[0].id),
    ])
    error_message = "cannot set all expected system datasource fields"
  }

}

run "datasource_system_first" {

  assert {
    condition     = data.opslevel_system.first_system_by_id.id == data.opslevel_systems.all.systems[0].id
    error_message = "wrong ID on opslevel_system"
  }

}
