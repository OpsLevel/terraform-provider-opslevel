mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_service_dependency_with_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_dependency.with_alias.depends_upon == var.test_id
    error_message = "wrong depends_upon in opslevel_service_dependency.with_alias"
  }

  assert {
    condition     = can(opslevel_service_dependency.with_alias.id)
    error_message = "id attribute missing from filter in opslevel_service_dependency.with_alias"
  }

  assert {
    condition     = opslevel_service_dependency.with_alias.service == var.test_id
    error_message = "wrong service in opslevel_service_dependency.with_alias"
  }

}

run "resource_service_dependency_with_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_dependency.with_id.depends_upon == var.test_id
    error_message = "wrong depends_upon in opslevel_service_dependency.with_id"
  }

  assert {
    condition     = can(opslevel_service_dependency.with_id.id)
    error_message = "id attribute missing from filter in opslevel_service_dependency.with_id"
  }

  assert {
    condition     = opslevel_service_dependency.with_id.note == <<-EOT
    This is an example of notes on a service dependency
  EOT
    error_message = "wrong note in opslevel_service_dependency.with_id"
  }

  assert {
    condition     = opslevel_service_dependency.with_id.service == var.test_id
    error_message = "wrong service in opslevel_service_dependency.with_id"
  }

}
