mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_team_contact" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_team_contact.tc_1.id)
    error_message = "expected to have an ID"
  }

  assert {
    condition     = can(opslevel_team_contact.tc_2.id)
    error_message = "expected to have an ID"
  }

  assert {
    condition     = opslevel_team_contact.tc_1.team == "team_platform_3"
    error_message = "has unexpected team"
  }

  assert {
    condition     = opslevel_team_contact.tc_2.team == "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQxMg"
    error_message = "has unexpected team"
  }

  assert {
    condition     = opslevel_team_contact.tc_1.type == "slack"
    error_message = "has unexpected type"
  }

  assert {
    condition     = opslevel_team_contact.tc_2.type == "email"
    error_message = "has unexpected type"
  }

  assert {
    condition     = opslevel_team_contact.tc_1.value == "#platform-3"
    error_message = "has unexpected value"
  }

  assert {
    condition     = opslevel_team_contact.tc_2.value == "team-platform-3-3-3@opslevel.com"
    error_message = "has unexpected value"
  }
}