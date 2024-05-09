run "datasource_integrations_all" {

  assert {
    condition     = length(data.opslevel_integrations.all.integrations) > 0
    error_message = "zero integrations found in data.opslevel_integrations"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_integrations.all.integrations[0].id),
      can(data.opslevel_integrations.all.integrations[0].name),
    ])
    error_message = "cannot set all expected integration datasource fields"
  }

}

run "datasource_integration_first" {

  assert {
    condition     = data.opslevel_integration.first_integration_by_id.id == data.opslevel_integrations.all.integrations[0].id
    error_message = "wrong ID on opslevel_integration"
  }

  assert {
    condition     = data.opslevel_integration.first_integration_by_name.name == data.opslevel_integrations.all.integrations[0].name
    error_message = "wrong name on opslevel_integration"
  }

}
