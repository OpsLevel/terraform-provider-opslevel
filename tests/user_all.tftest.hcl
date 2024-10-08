variables {
  resource_name = "opslevel_user"

  # required fields
  email = "foo@bar.com"
  name  = "TF Test User"

  # optional fields
  role               = "user" # make required??
  send_invite        = true
  skip_welcome_email = false
}

run "resource_user_create_with_all_fields" {

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
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
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
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
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
      var.email,
      opslevel_user.this.email,
    )
  }

}

run "resource_user_unset_skip_welcome_email_defaults_to_true" {

  variables {
    skip_welcome_email = null
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition     = opslevel_user.this.skip_welcome_email == true
    error_message = "expected 'true' default for skip_welcome_email in opslevel_user resource"
  }

}
