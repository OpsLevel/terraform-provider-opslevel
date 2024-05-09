mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_user" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_user.mock_user.email == "mock-user-name@example.com"
    error_message = "wrong email address of opslevel_user"
  }

  assert {
    condition     = data.opslevel_user.mock_user.id != null && data.opslevel_user.mock_user.id != ""
    error_message = "opslevel_user id should not be empty"
  }

  assert {
    condition     = contains([data.opslevel_user.mock_user.email, data.opslevel_user.mock_user.id], data.opslevel_user.mock_user.identifier)
    error_message = "wrong identifier of opslevel_user, should match id or email"
  }

  assert {
    condition     = data.opslevel_user.mock_user.name == "mock-user-name"
    error_message = "wrong name of opslevel_user"
  }

  assert {
    condition     = data.opslevel_user.mock_user.role == "mock-user-role"
    error_message = "wrong role of opslevel_user"
  }
}

