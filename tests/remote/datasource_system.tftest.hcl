run "datasource_systems_all" {

  variables {
    datasource_type = "opslevel_systems"
  }

  assert {
    condition     = can(data.opslevel_systems.all.systems)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_systems.all.systems) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_system_first" {

  variables {
    datasource_type = "opslevel_system"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_system.first_system_by_id.aliases),
      can(data.opslevel_system.first_system_by_id.description),
      can(data.opslevel_system.first_system_by_id.domain),
      can(data.opslevel_system.first_system_by_id.id),
      can(data.opslevel_system.first_system_by_id.identifier),
      can(data.opslevel_system.first_system_by_id.name),
      can(data.opslevel_system.first_system_by_id.owner),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_system.first_system_by_id.id == data.opslevel_systems.all.systems[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

}
