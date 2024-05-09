run "datasource_webhook_actions_all" {

  variables {
    datasource_type = "opslevel_webhook_actions"
  }

  assert {
    condition     = can(data.opslevel_webhook_actions.all.webhook_actions)
    error_message = "cannot set all expected webhook_action datasource fields"
  }

  assert {
    condition     = length(data.opslevel_webhook_actions.all.webhook_actions) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_webhook_action_first" {

  variables {
    datasource_type = "opslevel_webhook_action"
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
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

}
