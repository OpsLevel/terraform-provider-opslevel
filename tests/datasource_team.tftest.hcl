mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_team_with_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_alias.alias == "platform"
    error_message = "wrong alias on opslevel_team"
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_alias.id != null && data.opslevel_team.mock_team_with_alias.id != ""
    error_message = "empty id on opslevel_team"
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_alias.name == "Platform"
    error_message = "wrong name on opslevel_team"
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_alias.parent_alias == "engineering"
    error_message = "wrong parent_alias on opslevel_team"
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_alias.parent_id == "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xMDI0Mg"
    error_message = "wrong parent_id on opslevel_team"
  }
}

run "datasource_team_with_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_id.alias == "platform"
    error_message = "wrong alias on opslevel_team"
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_id.id != null && data.opslevel_team.mock_team_with_id.id != ""
    error_message = "empty id on opslevel_team"
  }
}
