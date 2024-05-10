run "datasource_users_all" {

  variables {
    datasource_type = "opslevel_users"
  }

  module {
    source = "./user"
  }

  assert {
    condition     = can(data.opslevel_users.all.users)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_users.all.users) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

}

run "datasource_user_first" {

  variables {
    datasource_type = "opslevel_user"
  }

  module {
    source = "./user"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_user.first_user_by_id.email),
      can(data.opslevel_user.first_user_by_id.id),
      can(data.opslevel_user.first_user_by_id.identifier),
      can(data.opslevel_user.first_user_by_id.name),
      can(data.opslevel_user.first_user_by_id.role),
    ])
    error_message = "cannot reference all expected user datasource fields"
  }

  assert {
    condition     = data.opslevel_user.first_user_by_id.email == data.opslevel_users.all.users[0].email
    error_message = "wrong email on opslevel_user"
  }

  assert {
    condition     = data.opslevel_user.first_user_by_id.id == data.opslevel_users.all.users[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

}

