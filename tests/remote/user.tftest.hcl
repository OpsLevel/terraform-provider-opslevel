variables {
  user_one  = "opslevel_user"
  users_all = "opslevel_users"

  # required fields
  email = "foo@bar.com"
  name  = "TF Test User"

  # optional fields
  role = "user"
}

run "resource_user_create_with_all_fields" {

  variables {
    email = var.email
    name  = var.name
    role  = var.role
  }

  module {
    source = "./user"
  }

}

run "resource_user_update_unset_optional_fields" {

  variables {
    role  = null
  }

  module {
    source = "./user"
  }

}

run "resource_user_update_set_optional_fields" {

  variables {
    email = var.email
    name  = "${var.name} updated"
    role  = "admin"
  }

  module {
    source = "./user"
  }

}

run "datasource_users_all" {

  module {
    source = "./user"
  }

  assert {
    condition     = can(data.opslevel_users.all.users)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.users_all)
  }

  assert {
    condition     = length(data.opslevel_users.all.users) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.users_all)
  }

}

run "datasource_user_first" {

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
    error_message = replace(var.error_wrong_id, "TYPE", var.user_one)
  }

}
