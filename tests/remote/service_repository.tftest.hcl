variables {
  resource_name = "opslevel_service_repository"

  # required fields

  # optional fields
  base_directory   = null
  name             = null
  repository       = null
  repository_alias = null
  service          = null
  service_alias    = null
}

run "from_repository_module" {
  command = plan

  variables {
    identifier = ""
  }

  module {
    source = "./repository"
  }
}

run "from_service_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./service"
  }
}

run "resource_service_repository_create_with_ids" {

  variables {
    base_directory   = "base/path"
    name             = "TF test service repository"
    repository       = run.from_repository_module.first_repository.id
    repository_alias = null
    service          = run.from_service_module.first_service.id
    service_alias    = null
  }

  module {
    source = "./service_repository"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_repository.test.base_directory),
      can(opslevel_service_repository.test.id),
      can(opslevel_service_repository.test.name),
      can(opslevel_service_repository.test.repository),
      can(opslevel_service_repository.test.repository_alias),
      can(opslevel_service_repository.test.service),
      can(opslevel_service_repository.test.service_alias),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_repository.test.base_directory == var.base_directory
    error_message = format(
      "expected '%v' but got '%v'",
      var.base_directory,
      opslevel_service_repository.test.base_directory,
    )
  }

  assert {
    condition     = startswith(opslevel_service_repository.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_repository.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_service_repository.test.name,
    )
  }

  assert {
    condition = opslevel_service_repository.test.repository == var.repository
    error_message = format(
      "expected '%v' but got '%v'",
      var.repository,
      opslevel_service_repository.test.repository,
    )
  }

  assert {
    condition     = opslevel_service_repository.test.repository_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_repository.test.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_repository.test.service,
    )
  }

  assert {
    condition     = opslevel_service_repository.test.service_alias == null
    error_message = var.error_expected_null_field
  }

}

run "resource_service_repository_create_with_aliases" {

  variables {
    base_directory   = "base/path"
    name             = "TF test service repository"
    repository       = null
    repository_alias = run.from_repository_module.first_repository.alias
    service          = null
    service_alias    = run.from_service_module.first_service.aliases[0]
  }

  module {
    source = "./service_repository"
  }

  assert {
    condition = opslevel_service_repository.test.base_directory == var.base_directory
    error_message = format(
      "expected '%v' but got '%v'",
      var.base_directory,
      opslevel_service_repository.test.base_directory,
    )
  }

  assert {
    condition = opslevel_service_repository.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_service_repository.test.name,
    )
  }

  assert {
    condition     = opslevel_service_repository.test.repository == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_repository.test.repository_alias == var.repository_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.repository_alias,
      opslevel_service_repository.test.repository_alias,
    )
  }

  assert {
    condition     = opslevel_service_repository.test.service == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_repository.test.service_alias == var.service_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.service_alias,
      opslevel_service_repository.test.service_alias,
    )
  }

}

run "resource_service_repository_update_unset_optional_fields" {

  variables {
    base_directory   = null
    name             = "TF test service repository"
    repository       = null
    repository_alias = run.from_repository_module.first_repository.alias
    service          = null
    service_alias    = run.from_service_module.first_service.aliases[0]
  }

  module {
    source = "./service_repository"
  }

  assert {
    condition     = opslevel_service_repository.test.base_directory == null
    error_message = var.error_expected_null_field
  }

  # TODO: unable to unset 'name' field for now
  # assert {
  #   condition     = opslevel_service_repository.test.name == null
  #   error_message = var.error_expected_null_field
  # }

}
