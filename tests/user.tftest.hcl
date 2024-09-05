variables {
  user_one  = "opslevel_user"
  users_all = "opslevel_users"

  # required fields
  email = "foo@bar.com"
  name  = "TF Test User"

  # optional fields
  force_send_invite  = true
  role               = "user"
  skip_welcome_email = false
}

run "resource_user_create_with_all_fields" {

  variables {
    email              = var.email
    force_send_invite  = var.force_send_invite
    name               = var.name
    role               = var.role
    skip_welcome_email = var.skip_welcome_email
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition = alltrue([
      can(opslevel_user.this.email),
      can(opslevel_user.this.force_send_invite),
      can(opslevel_user.this.id),
      can(opslevel_user.this.name),
      can(opslevel_user.this.role),
      can(opslevel_user.this.skip_welcome_email),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.user_one)
  }

  assert {
    condition     = opslevel_user.this.force_send_invite == var.force_send_invite
    error_message = "wrong force_send_invite value for opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.email == var.email
    error_message = "wrong email for opslevel_user resource"
  }

  assert {
    condition     = startswith(opslevel_user.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.user_one)
  }

  assert {
    condition     = opslevel_user.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.user_one)
  }

  assert {
    condition     = opslevel_user.this.role == var.role
    error_message = "wrong role for opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.skip_welcome_email == var.skip_welcome_email
    error_message = "wrong email for opslevel_user resource"
  }

}

run "resource_user_update_unset_fields_return_default_value" {

  variables {
    force_send_invite  = null
    skip_welcome_email = null
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition     = opslevel_user.this.force_send_invite == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_user.this.skip_welcome_email == true
    error_message = "expected 'true' default for skip_welcome_email in opslevel_user resource"
  }

}

run "resource_user_update_set_all_fields" {

  variables {
    email              = var.email
    force_send_invite  = !var.force_send_invite
    name               = "${var.name} updated"
    role               = var.role == "user" ? "admin" : var.role
    skip_welcome_email = !var.skip_welcome_email
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition     = opslevel_user.this.email == var.email
    error_message = "wrong email for opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.force_send_invite == var.force_send_invite
    error_message = "wrong force_send_invite for opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.user_one)
  }

  assert {
    condition     = opslevel_user.this.role == var.role
    error_message = "wrong role for opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.skip_welcome_email == var.skip_welcome_email
    error_message = "wrong skip_welcome_email for opslevel_user resource"
  }

}

run "datasource_users_all" {

  module {
    source = "./opslevel_modules/modules/user"
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
