variables {
  resource_name = "opslevel_user"

  # required fields
  email = "foo@bar.com"
  name  = "TF Test User"

  # optional fields
  role               = "user" # make required??
  send_invite        = null
  skip_welcome_email = null
}

run "resource_user_create_with_required_fields" {

  module {
    source = "./opslevel_modules/modules/user"
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
    condition = opslevel_user.this.role == "user"
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_user.this.name,
    )
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
    condition     = opslevel_user.this.send_invite == false
    error_message = "expected 'false' default for send_invite in opslevel_user resource"
  }

  assert {
    condition     = opslevel_user.this.skip_welcome_email == true
    error_message = "expected 'true' default for skip_welcome_email in opslevel_user resource"
  }

}

run "resource_user_set_role_to_admin" {

  variables {
    role = "admin"
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition = opslevel_user.this.role == var.role
    error_message = format(
      "expected '%v' but got '%v'",
      var.role,
      opslevel_user.this.role,
    )
  }

}

run "resource_user_set_role_to_standards_admin" {

  variables {
    role = "standards_admin"
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition = opslevel_user.this.role == var.role
    error_message = format(
      "expected '%v' but got '%v'",
      var.role,
      opslevel_user.this.role,
    )
  }

}

run "resource_user_set_role_to_team_member" {

  variables {
    role = "team_member"
  }

  module {
    source = "./opslevel_modules/modules/user"
  }

  assert {
    condition = opslevel_user.this.role == var.role
    error_message = format(
      "expected '%v' but got '%v'",
      var.role,
      opslevel_user.this.role,
    )
  }

}
