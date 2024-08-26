# variables {
#   resource_name = "opslevel_check_custom_event"


#   # -- check_custom_event fields --
#   # required fields
#   integration       = "(Sample) Cart"
#   pass_pending      = true
#   service_selector  = ".[] | select(.service == \"datadog\") | .service"
#   success_condition = ".[] | select(.service == \"datadog\") | .success"

#   # optional fields
#   message = null

#   # -- check base fields --
#   # required fields
#   category = null
#   level    = null
#   name     = "TF Test Check Custom Event"

#   # optional fields
#   enable_on = null
#   enabled   = true
#   filter    = null
#   notes     = "Notes on TF Check Custom Event"
#   owner     = null
# }

# run "from_filter_get_filter_id" {
#   command = plan

#   variables {
#     connective = null
#   }

#   module {
#     source = "./filter"
#   }
# }

# run "from_integration_module" {
#   command = plan

#   variables {
#     name = "(Sample) Cart"
#   }

#   module {
#     source = "./integration"
#   }
# }

# run "from_rubric_category_get_category_id" {
#   command = plan

#   variables {
#     name = ""
#   }

#   module {
#     source = "./rubric_category"
#   }
# }

# run "from_rubric_level_get_level_id" {
#   command = plan

#   variables {
#     description = null
#     index       = null
#     name        = ""
#   }

#   module {
#     source = "./rubric_level"
#   }
# }

# run "from_team_get_owner_id" {
#   command = plan

#   variables {
#     aliases          = null
#     name             = ""
#     parent           = null
#     responsibilities = null
#   }

#   module {
#     source = "./team"
#   }
# }

# run "resource_check_custom_event_create_with_all_fields" {

#   variables {
#     integration       = run.from_integration_module.first_integration.id
#     message           = var.message
#     pass_pending      = var.pass_pending
#     service_selector  = var.service_selector
#     success_condition = var.success_condition

#     category  = run.from_rubric_category_get_category_id.first_category.id
#     enable_on = var.enable_on
#     enabled   = var.enabled
#     filter    = run.from_filter_get_filter_id.first_filter.id
#     level     = run.from_rubric_level_get_level_id.greatest_level.id
#     name      = var.name
#     notes     = var.notes
#     owner     = run.from_team_get_owner_id.first_team.id
#   }

#   module {
#     source = "./check_custom_event"
#   }

#   assert {
#     condition = alltrue([
#       can(opslevel_check_custom_event.test.category),
#       can(opslevel_check_custom_event.test.enable_on),
#       can(opslevel_check_custom_event.test.enabled),
#       can(opslevel_check_custom_event.test.filter),
#       can(opslevel_check_custom_event.test.id),
#       can(opslevel_check_custom_event.test.integration),
#       can(opslevel_check_custom_event.test.level),
#       can(opslevel_check_custom_event.test.message),
#       can(opslevel_check_custom_event.test.name),
#       can(opslevel_check_custom_event.test.notes),
#       can(opslevel_check_custom_event.test.owner),
#       can(opslevel_check_custom_event.test.pass_pending),
#       can(opslevel_check_custom_event.test.service_selector),
#       can(opslevel_check_custom_event.test.success_condition),
#     ])
#     error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.category == var.category
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.category,
#       opslevel_service_tag.test.category,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.enable_on == var.enable_on
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.enable_on,
#       opslevel_service_tag.test.enable_on,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.enabled == var.enabled
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.enabled,
#       opslevel_service_tag.test.enabled,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.filter == var.filter
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.filter,
#       opslevel_service_tag.test.filter,
#     )
#   }

#   assert {
#     condition     = startswith(opslevel_check_custom_event.test.id, var.id_prefix)
#     error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.integration == var.integration
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.integration,
#       opslevel_service_tag.test.integration,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.level == var.level
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.level,
#       opslevel_service_tag.test.level,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.message == var.message
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.message,
#       opslevel_service_tag.test.message,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.name == var.name
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.name,
#       opslevel_service_tag.test.name,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.notes == var.notes
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.notes,
#       opslevel_service_tag.test.notes,
#     )
#   }

#   assert {
#     condition = opslevel_check_custom_event.test.owner == var.owner
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.owner,
#       opslevel_service_tag.test.owner,
#     )
#   }

# }
