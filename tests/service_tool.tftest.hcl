variables {
  resource_name = "opslevel_service_tool"

  # required fields
  category = "observability"
  name     = "TF Test Service Tool"
  url      = "https://example.com"

  # optional fields
  environment   = "production"
  service       = null
  service_alias = null
}

run "from_service_module" {
  command = plan

  module {
    source = "./data/service"
  }
}

run "resource_service_tool_create_with_service_id" {

  variables {
    category      = var.category
    environment   = var.environment
    name          = var.name
    service       = run.from_service_module.first.id
    service_alias = null
    url           = var.url
  }

  module {
    source = "./opslevel_modules/modules/service/tool"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_tool.this.category),
      can(opslevel_service_tool.this.environment),
      can(opslevel_service_tool.this.id),
      can(opslevel_service_tool.this.name),
      can(opslevel_service_tool.this.service),
      can(opslevel_service_tool.this.service_alias),
      can(opslevel_service_tool.this.url),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_tool.this.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_service_tool.this.category,
    )
  }

  assert {
    condition = opslevel_service_tool.this.environment == var.environment
    error_message = format(
      "expected '%v' but got '%v'",
      var.environment,
      opslevel_service_tool.this.environment,
    )
  }

  assert {
    condition     = startswith(opslevel_service_tool.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_tool.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_service_tool.this.name,
    )
  }

  assert {
    condition = opslevel_service_tool.this.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_tool.this.service,
    )
  }

  assert {
    condition     = opslevel_service_tool.this.service_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_tool.this.url == var.url
    error_message = format(
      "expected '%v' but got '%v'",
      var.url,
      opslevel_service_tool.this.url,
    )
  }

}

# run "resource_service_tool_create_with_service_alias" {

#   variables {
#     category      = var.category
#     environment   = var.environment
#     name          = var.name
#     service       = null
#     service_alias = run.from_service_module.first.aliases[0]
#     url           = var.url
#   }

#   module {
#     source = "./opslevel_modules/modules/service/tool"
#   }

#   assert {
#     condition = opslevel_service_tool.this.category == var.category
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.category,
#       opslevel_service_tool.this.category,
#     )
#   }

#   assert {
#     condition = opslevel_service_tool.this.environment == var.environment
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.environment,
#       opslevel_service_tool.this.environment,
#     )
#   }

#   assert {
#     condition     = startswith(opslevel_service_tool.this.id, var.id_prefix)
#     error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
#   }

#   assert {
#     condition     = opslevel_service_tool.this.service == null
#     error_message = var.error_expected_null_field
#   }

#   assert {
#     condition = opslevel_service_tool.this.service_alias == var.service_alias
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.service_alias,
#       opslevel_service_tool.this.service_alias,
#     )
#   }

#   assert {
#     condition = opslevel_service_tool.this.url == var.url
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.url,
#       opslevel_service_tool.this.url,
#     )
#   }

# }
