mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_team_tag_using_team_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_team_tag.using_team_id.key == "hello_with_id"
    error_message = "wrong team tag key"
  }

  assert {
    condition     = opslevel_team_tag.using_team_id.value == "world_with_id"
    error_message = "wrong team tag value"
  }

  assert {
    condition     = opslevel_team_tag.using_team_id.team == "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQxMg"
    error_message = "expected team identifier to be an id"
  }

  assert {
    condition     = can(opslevel_team_tag.using_team_id.id)
    error_message = "expected team tag to have an ID"
  }
}

run "resource_team_tag_using_team_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_team_tag.using_team_alias.key == "hello_with_alias"
    error_message = "wrong team tag key"
  }

  assert {
    condition     = opslevel_team_tag.using_team_alias.value == "world_with_alias"
    error_message = "wrong team tag value"
  }

  assert {
    condition     = opslevel_team_tag.using_team_alias.team_alias == "team_platform_3"
    error_message = "expected team identifier to be an alias"
  }

  assert {
    condition     = can(opslevel_team_tag.using_team_alias.id)
    error_message = "expected team tag to have an ID"
  }
}