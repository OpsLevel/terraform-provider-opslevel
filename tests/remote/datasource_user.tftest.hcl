run "datasource_users_all" {

  assert {
    condition = alltrue([
      can(data.opslevel_users.all.users),
    ])
    error_message = "cannot set all expected users datasource fields"
  }

  assert {
    condition     = length(data.opslevel_users.all.users) > 0
    error_message = "zero users found in data.opslevel_users"
  }

}

run "datasource_user_first" {

  assert {
    condition = alltrue([
      can(data.opslevel_user.first_user_by_id.email),
      can(data.opslevel_user.first_user_by_id.id),
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
    error_message = "wrong ID on opslevel_user"
  }

}

