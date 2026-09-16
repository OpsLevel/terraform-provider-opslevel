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

  #assert {
  #  condition     = data.opslevel_team.mock_team_with_alias.members[0].email == "person1@opslevel.com" && data.opslevel_team.mock_team_with_alias.members[0].role == "manager"
  #  error_message = "wrong members data on opslevel_team"
  #}

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

run "datasource_team_tags_and_properties" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_id.tags == tolist(["key1:value2", "key2:value2"])
    error_message = "wrong tags list in opslevel_team mock"
  }

  assert {
    condition     = data.opslevel_team.mock_team_with_id.properties != null
    error_message = "null properties on opslevel_team"
  }

  # NOTE: mock_provider does not populate nested list attributes, same as opslevel_service
  #assert {
  #  condition     = data.opslevel_team.mock_team_with_id.properties[0].definition.id == "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8yODk"
  #  error_message = "wrong definition.id in properties.definition in opslevel_team mock"
  #}

  #assert {
  #  condition     = data.opslevel_team.mock_team_with_id.properties[0].value == "mock-property-definition"
  #  error_message = "wrong value in properties list in opslevel_team mock"
  #}
}
