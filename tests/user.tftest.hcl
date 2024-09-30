variables {
  user_one  = "opslevel_user"
  users_all = "opslevel_users"

  # required fields
  email = "foo@bar.com"
  name  = "TF Test User"

  # optional fields
  role               = "user"
  send_invite        = true
  skip_welcome_email = false
}

run "resource_user_create_with_all_fields" {

  variables {
    email              = var.email
    name               = var.name
    role               = var.role
    send_invite        = var.send_invite
    skip_welcome_email = var.skip_welcome_email
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition = alltrue([
      can(opslevel_user.this.email),
      can(opslevel_user.this.id),
      can(opslevel_user.this.name),
      can(opslevel_user.this.role),
      can(opslevel_user.this.send_invite),
      can(opslevel_user.this.skip_welcome_email),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.user_one)
  }

  assert {
    condition = opslevel_user.this.email == var.email
    error_message = format(
      "expected '%v' but got '%v'",
      var.email,
      opslevel_user.this.email,
    )
  }

  assert {
    condition     = startswith(opslevel_user.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.user_one)
  }

  assert {
    condition = opslevel_user.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_user.this.name,
    )
  }

  assert {
    condition = opslevel_user.this.role == var.role
    error_message = format(
      "expected '%v' but got '%v'",
      var.role,
      opslevel_user.this.role,
    )
  }

  assert {
    condition = opslevel_user.this.send_invite == var.send_invite
    error_message = format(
      "expected '%v' but got '%v'",
      var.send_invite,
      opslevel_user.this.send_invite,
    )
  }

  assert {
    condition = opslevel_user.this.skip_welcome_email == var.skip_welcome_email
    error_message = format(
      "expected '%v' but got '%v'",
      var.skip_welcome_email,
      opslevel_user.this.skip_welcome_email,
    )
  }

}

run "resource_user_update_unset_fields_return_default_value" {

  variables {
    send_invite        = null
    skip_welcome_email = null
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition     = opslevel_user.this.send_invite == true
    error_message = "expected 'true' default for send_invite in opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.skip_welcome_email == true
    error_message = "expected 'true' default for skip_welcome_email in opslevel_user resource"
  }

}

run "resource_user_update_set_all_fields" {

  variables {
    email              = var.email
    name               = "${var.name} updated"
    role               = var.role == "user" ? "admin" : var.role
    send_invite        = !var.send_invite
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
    condition     = opslevel_user.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.user_one)
  }

  assert {
    condition     = opslevel_user.this.role == var.role
    error_message = "wrong role for opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.send_invite == var.send_invite
    error_message = "wrong send_invite for opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.skip_welcome_email == var.skip_welcome_email
    error_message = "wrong skip_welcome_email for opslevel_user resource"
  }

}
