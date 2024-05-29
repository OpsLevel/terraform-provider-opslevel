variables {
  webhook_action_one = "opslevel_webhook_action"
  webhook_action_all = "opslevel_webhook_actions"

  # required fields
  method  = "GET"
  name    = "TF Test Webhook Action"
  payload = "{\"operation:\": \"create\"}"
  url     = "https://webhook.url"

  # optional fields
  description = "Webhook Action description"
  headers     = tomap({ operation = "create", fields = "all" })
}

run "resource_webhook_action_create_with_all_fields" {

  variables {
    description = var.description
    headers     = var.headers
    method      = var.method
    name        = var.name
    payload     = var.payload
    url         = var.url
  }

  module {
    source = "./webhook_action"
  }

  assert {
    condition = alltrue([
      can(opslevel_webhook_action.test.description),
      can(opslevel_webhook_action.test.headers),
      can(opslevel_webhook_action.test.id),
      can(opslevel_webhook_action.test.method),
      can(opslevel_webhook_action.test.name),
      can(opslevel_webhook_action.test.payload),
      can(opslevel_webhook_action.test.url),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.webhook_action_one)
  }

  assert {
    condition     = opslevel_webhook_action.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.webhook_action_one)
  }

  assert {
    condition     = opslevel_webhook_action.test.headers == var.headers
    error_message = "wrong headers in opslevel_webhook_action resource"
  }

  assert {
    condition     = startswith(opslevel_webhook_action.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.webhook_action_one)
  }

  assert {
    condition     = opslevel_webhook_action.test.method == var.method
    error_message = "wrong method in opslevel_webhook_action resource"
  }

  assert {
    condition     = opslevel_webhook_action.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.webhook_action_one)
  }

  assert {
    condition     = opslevel_webhook_action.test.payload == var.payload
    error_message = "wrong payload in opslevel_webhook_action resource"
  }

  assert {
    condition     = opslevel_webhook_action.test.url == var.url
    error_message = "wrong url in opslevel_webhook_action resource"
  }

}

run "resource_webhook_action_update_unset_optional_fields" {

  variables {
    description = null
    headers     = null
  }

  module {
    source = "./webhook_action"
  }

  assert {
    condition     = opslevel_webhook_action.test.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_webhook_action.test.headers == var.headers
    error_message = var.error_expected_null_field
  }

}

run "resource_webhook_action_update_all_fields" {

  variables {
    description = "${var.description} updated"
    headers     = tomap({ operation = "update", fields = "all" })
    method      = "POST"
    name        = "${var.name} updated"
    payload     = "{\"payload:\": \"updated\"}"
    url         = "${var.url}/updated"
  }

  module {
    source = "./webhook_action"
  }

  assert {
    condition     = opslevel_webhook_action.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.webhook_action_one)
  }

  assert {
    condition     = opslevel_webhook_action.test.headers == var.headers
    error_message = "wrong headers in opslevel_webhook_action resource"
  }

  assert {
    condition     = startswith(opslevel_webhook_action.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.webhook_action_one)
  }

  assert {
    condition     = opslevel_webhook_action.test.method == var.method
    error_message = "wrong method in opslevel_webhook_action resource"
  }

  assert {
    condition     = opslevel_webhook_action.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.webhook_action_one)
  }

  assert {
    condition     = opslevel_webhook_action.test.payload == var.payload
    error_message = "wrong payload in opslevel_webhook_action resource"
  }

  assert {
    condition     = opslevel_webhook_action.test.url == var.url
    error_message = "wrong url in opslevel_webhook_action resource"
  }

}

run "datasource_webhook_actions_all" {

  module {
    source = "./webhook_action"
  }

  assert {
    condition     = can(data.opslevel_webhook_actions.all.webhook_actions)
    error_message = "cannot set all expected webhook_action datasource fields"
  }

  assert {
    condition     = length(data.opslevel_webhook_actions.all.webhook_actions) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.webhook_action_all)
  }

}

run "datasource_webhook_action_first" {

  module {
    source = "./webhook_action"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_webhook_action.first_webhook_action_by_id.aliases),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.description),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.headers),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.id),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.identifier),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.method),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.name),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.payload),
      can(data.opslevel_webhook_action.first_webhook_action_by_id.url),
    ])
    error_message = "cannot set all expected webhook_action datasource fields"
  }

  assert {
    condition     = data.opslevel_webhook_action.first_webhook_action_by_id.id == data.opslevel_webhook_actions.all.webhook_actions[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.webhook_action_one)
  }

}
