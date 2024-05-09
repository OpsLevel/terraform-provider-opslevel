run "datasource_webhook_actions_all" {

  assert {
    condition     = length(data.opslevel_webhook_actions.all.webhook_actions) > 0
    error_message = "zero webhook_actions found in data.opslevel_webhook_actions"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_webhook_actions.all.webhook_actions[0].id),
    ])
    error_message = "cannot set all expected webhook_action datasource fields"
  }

}

run "datasource_webhook_action_first" {

  assert {
    condition     = data.opslevel_webhook_action.first_webhook_action_by_id.id == data.opslevel_webhook_actions.all.webhook_actions[0].id
    error_message = "wrong ID on opslevel_webhook_action"
  }

}

