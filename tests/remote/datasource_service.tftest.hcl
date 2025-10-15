variables {
  api_token = ""
}

run "datasource_service_reads_fields" {
  command = plan

  variables {
    api_token = var.api_token
    name      = "Test Service for Data Source"
  }

  module {
    source = "./service"
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.id != null
    error_message = "BUG: service id should not be null"
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.name != null
    error_message = "BUG: service name should not be null - this was the reported issue"
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.description != null
    error_message = "BUG: service description should not be null - this was the reported issue"
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.owner != null
    error_message = "BUG: service owner should not be null"
  }

  assert {
    condition     = data.opslevel_service.last_service.id != null
    error_message = "BUG: service id should not be null"
  }

  assert {
    condition     = data.opslevel_service.last_service.name != null
    error_message = "BUG: service name should not be null"
  }
}

