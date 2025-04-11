variables {
  resource_name = "opslevel_trigger_definition"

  # required fields
  access_control = "service_owners"
  action         = null
  name           = "TF Test Trigger Definition"
  owner          = null
  published      = true

  # optional fields
  approval_required        = true
  approval_users           = []
  description              = "TF Test Trigger Definition description"
  entity_type              = "SERVICE"
  extended_team_access     = [] # team aliases
  filter                   = null
  manual_inputs_definition = <<EOT
---
version: 1
inputs:
  - identifier: IncidentTitle
    displayName: Title
    description: Title of the incident to trigger
    type: text_input
    required: true
    maxLength: 60
    defaultValue: Service Incident Manual Trigger
  - identifier: IncidentDescription
    displayName: Incident Description
    description: The description of the incident
    type: text_area
    required: true
  EOT
  response_template        = <<EOT
{% if response.status >= 200 and response.status < 300 %}
## Congratulations!
Your request for {{ service.name }} has succeeded. See the incident here: {{response.body.incident.html_url}}
{% else %}
## Oops something went wrong!
Please contact [{{ action_owner.name }}]({{ action_owner.href }}) for more help.
{% endif %}
  EOT
}

run "from_filter_get_filter_id" {
  command = plan

  variables {
    connective = null
  }

  module {
    source = "./filter"
  }
}

run "from_team_get_owner_id" {
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

run "from_webhook_action_get_webhook_action_id" {
  command = plan

  variables {
    method  = "GET"
    name    = ""
    payload = ""
    url     = ""
  }

  module {
    source = "./webhook_action"
  }
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_teams.all,
      data.opslevel_users.all
    ]
  }

  module {
    source = "../data"
  }
}
run "resource_trigger_definition_create_with_all_fields" {

  variables {
    access_control           = var.access_control
    action                   = run.from_webhook_action_get_webhook_action_id.first_webhook_action.id
    approval_required        = var.approval_required
    approval_users           = [run.from_data_module.all_users.users[0].email]
    description              = var.description
    entity_type              = var.entity_type
    extended_team_access     = var.extended_team_access
    filter                   = run.from_filter_get_filter_id.first_filter.id
    manual_inputs_definition = var.manual_inputs_definition
    name                     = var.name
    owner                    = run.from_team_get_owner_id.first_team.id
    published                = var.published
    response_template        = var.response_template
  }

  module {
    source = "./trigger_definition"
  }

  assert {
    condition = alltrue([
      can(opslevel_trigger_definition.test.access_control),
      can(opslevel_trigger_definition.test.action),
      can(opslevel_trigger_definition.test.approval_required),
      can(opslevel_trigger_definition.test.description),
      can(opslevel_trigger_definition.test.entity_type),
      can(opslevel_trigger_definition.test.extended_team_access),
      can(opslevel_trigger_definition.test.filter),
      can(opslevel_trigger_definition.test.manual_inputs_definition),
      can(opslevel_trigger_definition.test.name),
      can(opslevel_trigger_definition.test.owner),
      can(opslevel_trigger_definition.test.published),
      can(opslevel_trigger_definition.test.response_template),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_trigger_definition.test.access_control == var.access_control
    error_message = format(
      "expected '%v' but got '%v'",
      var.access_control,
      opslevel_trigger_definition.test.access_control,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.action == run.from_webhook_action_get_webhook_action_id.first_webhook_action.id
    error_message = format(
      "expected '%v' but got '%v'",
      run.from_webhook_action_get_webhook_action_id.first_webhook_action.id,
      opslevel_trigger_definition.test.action,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_trigger_definition.test.description,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.entity_type == var.entity_type
    error_message = format(
      "expected '%v' but got '%v'",
      var.entity_type,
      opslevel_trigger_definition.test.entity_type,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.extended_team_access == var.extended_team_access
    error_message = format(
      "expected '%v' but got '%v'",
      var.extended_team_access,
      opslevel_trigger_definition.test.extended_team_access,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.filter == var.filter
    error_message = format(
      "expected '%v' but got '%v'",
      var.filter,
      opslevel_trigger_definition.test.filter,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.manual_inputs_definition == var.manual_inputs_definition
    error_message = format(
      "expected '%v' but got '%v'",
      var.manual_inputs_definition,
      opslevel_trigger_definition.test.manual_inputs_definition,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_trigger_definition.test.name,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_trigger_definition.test.owner,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.published == var.published
    error_message = format(
      "expected '%v' but got '%v'",
      var.published,
      opslevel_trigger_definition.test.published,
    )
  }

  assert {
    condition = opslevel_trigger_definition.test.response_template == var.response_template
    error_message = format(
      "expected '%v' but got '%v'",
      var.response_template,
      opslevel_trigger_definition.test.response_template,
    )
  }

}

run "resource_trigger_definition_update_unset_fields" {

  variables {
    access_control           = var.access_control
    action                   = run.from_webhook_action_get_webhook_action_id.first_webhook_action.id
    description              = null
    entity_type              = null # TODO: explicitly set default to match API
    extended_team_access     = null
    filter                   = null
    manual_inputs_definition = null
    name                     = var.name
    owner                    = run.from_team_get_owner_id.first_team.id
    published                = var.published
    response_template        = null
  }

  module {
    source = "./trigger_definition"
  }

  assert {
    condition     = opslevel_trigger_definition.test.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_trigger_definition.test.entity_type == "SERVICE"
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_trigger_definition.test.extended_team_access == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_trigger_definition.test.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_trigger_definition.test.manual_inputs_definition == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_trigger_definition.test.response_template == null
    error_message = var.error_expected_null_field
  }

}
