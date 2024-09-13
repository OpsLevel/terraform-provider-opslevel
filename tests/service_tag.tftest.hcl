variables {
  resource_name = "opslevel_service_tag"

  # required fields
  key   = "test-service-tag-key"
  value = "test-service-tag-value"

  # optional fields
  service       = null
  service_alias = null
}

run "from_service_module" {
  command = plan

  module {
    source = "./data/service"
  }
}

run "resource_service_tag_create_with_service_id" {

  variables {
    key           = var.key
    service       = run.from_service_module.first.id
    service_alias = null
    value         = var.value
  }

  module {
    source = "./opslevel_modules/modules/service/tag"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_tag.this.key),
      can(opslevel_service_tag.this.id),
      can(opslevel_service_tag.this.service),
      can(opslevel_service_tag.this.service_alias),
      can(opslevel_service_tag.this.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_service_tag.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_tag.this.key == var.key
    error_message = format(
      "expected '%v' but got '%v'",
      var.key,
      opslevel_service_tag.this.key,
    )
  }

  assert {
    condition = opslevel_service_tag.this.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_tag.this.service,
    )
  }

  assert {
    condition     = opslevel_service_tag.this.service_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_tag.this.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_service_tag.this.value,
    )
  }

}

# run "resource_service_tag_create_with_service_alias" {

#   variables {
#     key           = var.key
#     service       = null
#     service_alias = run.from_service_module.first.aliases[0]
#     value         = var.value
#   }

#   module {
#     source = "./opslevel_modules/modules/service/tag"
#   }

#   assert {
#     condition = opslevel_service_tag.this.key == var.key
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.key,
#       opslevel_service_tag.this.key,
#     )
#   }

#   assert {
#     condition     = opslevel_service_tag.this.service == null
#     error_message = var.error_expected_null_field
#   }

#   assert {
#     condition = opslevel_service_tag.this.service_alias == var.service_alias
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.service_alias,
#       opslevel_service_tag.this.service_alias,
#     )
#   }

#   assert {
#     condition = opslevel_service_tag.this.value == var.value
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.value,
#       opslevel_service_tag.this.value,
#     )
#   }

# }
