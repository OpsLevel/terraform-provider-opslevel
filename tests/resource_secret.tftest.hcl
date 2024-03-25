mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_secret" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_secret.mock_secret.alias == "secret-alias"
    error_message = "wrong alias in opslevel_secret"
  }

  assert {
    condition     = opslevel_secret.mock_secret.created_at == "2022-02-24T13:50:07Z"
    error_message = "wrong created_at timestamp in opslevel_secret"
  }

  assert {
    condition     = contains([-1, 0], timecmp(opslevel_secret.mock_secret.created_at, opslevel_secret.mock_secret.updated_at))
    error_message = "created_at timestamp should not be after updated_at timestamp in opslevel_secret"
  }

  assert {
    condition     = opslevel_secret.mock_secret.id != null && opslevel_secret.mock_secret.id != ""
    error_message = "opslevel_secret id should not be empty"
  }

  assert {
    condition     = opslevel_secret.mock_secret.owner == "Developers"
    error_message = "wrong owner of opslevel_secret"
  }

  assert {
    condition     = opslevel_secret.mock_secret.value == "too_many_passwords"
    error_message = "wrong value of opslevel_secret"
  }

}

