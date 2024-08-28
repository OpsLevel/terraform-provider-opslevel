variables {
  resource_name = "opslevel_tag"

  # required fields
  key                 = "test-tag-key"
  value               = "test-tag-value"
  resource_identifier = null
  resource_type       = "Team"

  # optional fields - none
}

run "from_team_module" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }
}

run "resource_tag_create_with_all_fields" {

  variables {
    key                 = var.key
    resource_identifier = run.from_team_module.first_team.id
    resource_type       = var.resource_type
    value               = var.value
  }

  module {
    source = "./tag"
  }

  assert {
    condition = alltrue([
      can(opslevel_tag.test.key),
      can(opslevel_tag.test.id),
      can(opslevel_tag.test.resource_identifier),
      can(opslevel_tag.test.resource_type),
      can(opslevel_tag.test.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_tag.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_tag.test.key == var.key
    error_message = format(
      "expected '%v' but got '%v'",
      var.key,
      opslevel_tag.test.key,
    )
  }

  assert {
    condition = opslevel_tag.test.resource_identifier == var.resource_identifier
    error_message = format(
      "expected '%v' but got '%v'",
      var.resource_identifier,
      opslevel_tag.test.resource_identifier,
    )
  }

  assert {
    condition = opslevel_tag.test.resource_type == var.resource_type
    error_message = format(
      "expected '%v' but got '%v'",
      var.resource_type,
      opslevel_tag.test.resource_type,
    )
  }

  assert {
    condition = opslevel_tag.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_tag.test.value,
    )
  }

}

# BUG: https://github.com/OpsLevel/team-platform/issues/460
# run "resource_tag_update_key_and_value" {

#   variables {
#     key                 = "${var.key}-updated"
#     resource_identifier = run.from_team_module.first_team.id
#     resource_type       = var.resource_type
#     value               = "${var.value}-updated"
#   }

#   module {
#     source = "./tag"
#   }

#   assert {
#     condition = opslevel_tag.test.key == var.key
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.key,
#       opslevel_tag.test.key,
#     )
#   }

#   assert {
#     condition = opslevel_tag.test.value == var.value
#     error_message = format(
#       "expected '%v' but got '%v'",
#       var.value,
#       opslevel_tag.test.value,
#     )
#   }

# }
