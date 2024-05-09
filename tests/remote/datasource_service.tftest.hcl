run "datasource_services_all" {

  assert {
    condition     = length(data.opslevel_services.all.services) > 0
    error_message = "zero services found in data.opslevel_services"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_services.all.services[0].id),
    ])
    error_message = "cannot set all expected service datasource fields"
  }

}

run "datasource_service_first" {

  assert {
    condition     = data.opslevel_service.first_service_by_id.id == data.opslevel_services.all.services[0].id
    error_message = "wrong ID on opslevel_service"
  }

}
