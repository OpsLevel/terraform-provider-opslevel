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

  module {
    source = "./data/repository"
  }
}

run "from_service_module" {
  command = plan

  module {
    source = "./data/service"
  }
}

run "resource_service_repository_create_with_ids" {

  variables {
    base_directory   = "base/path"
    name             = "TF test service repository"
    repository       = run.from_repository_module.first.id
    repository_alias = null
    service          = run.from_service_module.first.id
    service_alias    = null
  }

  module {
    source = "./opslevel_modules/modules/service/repository"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_repository.this.base_directory),
      can(opslevel_service_repository.this.id),
      can(opslevel_service_repository.this.name),
      can(opslevel_service_repository.this.repository),
      can(opslevel_service_repository.this.repository_alias),
      can(opslevel_service_repository.this.service),
      can(opslevel_service_repository.this.service_alias),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_repository.this.base_directory == var.base_directory
    error_message = format(
      "expected '%v' but got '%v'",
      var.base_directory,
      opslevel_service_repository.this.base_directory,
    )
  }

  assert {
    condition     = startswith(opslevel_service_repository.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_repository.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_service_repository.this.name,
    )
  }

  assert {
    condition = opslevel_service_repository.this.repository == var.repository
    error_message = format(
      "expected '%v' but got '%v'",
      var.repository,
      opslevel_service_repository.this.repository,
    )
  }

  assert {
    condition     = opslevel_service_repository.this.repository_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_repository.this.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_repository.this.service,
    )
  }

  assert {
    condition     = opslevel_service_repository.this.service_alias == null
    error_message = var.error_expected_null_field
  }

}

# run "resource_service_repository_create_with_aliases" {

#   variables {
#     base_directory   = "base/path"
#     name             = "TF test service repository"
#     repository       = null
#     repository_alias = run.from_repository_module.first.alias
#     service          = null
#     service_alias    = run.from_service_module.first.aliases[0]
#   }

#   module {
#     source = "./opslevel_modules/modules/service/repository"
#   }

#   assert {
#     condition = opslevel_service_repository.this.base_directory == var.base_directory
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.base_directory,
#       opslevel_service_repository.this.base_directory,
#     )
#   }

#   assert {
#     condition = opslevel_service_repository.this.name == var.name
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.name,
#       opslevel_service_repository.this.name,
#     )
#   }

#   assert {
#     condition     = opslevel_service_repository.this.repository == null
#     error_message = var.error_expected_null_field
#   }

#   assert {
#     condition = opslevel_service_repository.this.repository_alias == var.repository_alias
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.repository_alias,
#       opslevel_service_repository.this.repository_alias,
#     )
#   }

#   assert {
#     condition     = opslevel_service_repository.this.service == null
#     error_message = var.error_expected_null_field
#   }

#   assert {
#     condition = opslevel_service_repository.this.service_alias == var.service_alias
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.service_alias,
#       opslevel_service_repository.this.service_alias,
#     )
#   }

# }

# run "resource_service_repository_update_unset_optional_fields" {

#   variables {
#     base_directory   = null
#     name             = "TF test service repository"
#     repository       = null
#     repository_alias = run.from_repository_module.first.alias
#     service          = null
#     service_alias    = run.from_service_module.first.aliases[0]
#   }

#   module {
#     source = "./opslevel_modules/modules/service/repository"
#   }

#   assert {
#     condition     = opslevel_service_repository.this.base_directory == null
#     error_message = var.error_expected_null_field
#   }

#   # TODO: unable to unset 'name' field for now
#   # assert {
#   #   condition     = opslevel_service_repository.this.name == null
#   #   error_message = var.error_expected_null_field
#   # }

# }
