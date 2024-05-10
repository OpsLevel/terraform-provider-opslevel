run "datasource_systems_all" {

  variables {
    datasource_type = "opslevel_systems"
  }

  module {
    source = "./system"
  }

  assert {
    condition     = can(data.opslevel_systems.all.systems)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_systems.all.systems) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

}

run "datasource_system_first" {

  variables {
    datasource_type = "opslevel_system"
  }

  module {
    source = "./system"
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
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_system.first_system_by_id.id == data.opslevel_systems.all.systems[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

}
