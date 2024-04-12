mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_service_repository_with_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_repository.with_alias.base_directory == "/home/opslevel"
    error_message = "wrong base_directory for opslevel_service_repository.with_alias"
  }

  assert {
    condition     = opslevel_service_repository.with_alias.name == "Service Repo Name"
    error_message = "wrong name for opslevel_service_repository.with_alias"
  }

  assert {
    condition     = opslevel_service_repository.with_alias.repository_alias == "github.com:OpsLevel/terraform-provider-opslevel"
    error_message = "wrong repository alias for opslevel_service_repository.with_alias"
  }

  assert {
    condition     = opslevel_service_repository.with_alias.service_alias == "service-1"
    error_message = "wrong service alias for opslevel_service_repository.with_alias"
  }

}

run "resource_service_repository_with_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_repository.with_id.repository == var.test_id
    error_message = "wrong repository id for opslevel_service_repository.with_id"
  }

  assert {
    condition     = opslevel_service_repository.with_id.service == var.test_id
    error_message = "wrong service id for opslevel_service_repository.with_id"
  }

  assert {
    condition     = can(opslevel_service_repository.with_id.id)
    error_message = "id attribute missing from in opslevel_service_repository.example"
  }

}
