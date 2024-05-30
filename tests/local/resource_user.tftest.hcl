mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_user" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_user.mock_user.email == "mock_user@mock.com"
    error_message = "wrong email for opslevel_user"
  }

  assert {
    condition     = opslevel_user.mock_user.id != null && opslevel_user.mock_user.id != ""
    error_message = "opslevel_user id should not be empty"
  }

  assert {
    condition     = opslevel_user.mock_user.name == "Mock User"
    error_message = "wrong name for opslevel_user"
  }

  assert {
    condition     = contains(["admin", "user"], opslevel_user.mock_user.role)
    error_message = "wrong role for opslevel_user"
  }

  assert {
    condition     = opslevel_user.mock_user.role == "user"
    error_message = "wrong role for opslevel_user"
  }

}

run "resource_user_no_role" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_user.mock_user_no_role.role == null
    error_message = "omitted role should be null for opslevel_user"
  }

}

run "resource_user_admin" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_user.mock_user_admin.role == "admin"
    error_message = "wrong role for opslevel_user"
  }

}
