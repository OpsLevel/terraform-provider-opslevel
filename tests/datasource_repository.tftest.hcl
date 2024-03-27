mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_repository_with_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_repository.mock_repository_with_alias.alias == "github.com:OpsLevel/opslevel-go"
    error_message = "wrong alias on opslevel_repository"
  }

  assert {
    condition     = data.opslevel_repository.mock_repository_with_alias.id != null && data.opslevel_repository.mock_repository_with_alias.id != ""
    error_message = "empty id on opslevel_repository"
  }
}

run "datasource_repository_with_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_repository.mock_repository_with_id.alias == "github.com:OpsLevel/opslevel-go"
    error_message = "wrong alias on opslevel_repository"
  }

  assert {
    condition     = data.opslevel_repository.mock_repository_with_id.id != null && data.opslevel_repository.mock_repository_with_id.id != ""
    error_message = "empty id on opslevel_repository"
  }
}
