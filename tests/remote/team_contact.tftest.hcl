variables {
  resource_name = "opslevel_team_contact"

  # required fields
  name  = "TF Test Team Contact"
  team  = null
  type  = null
  value = null

  # optional fields - none
}

run "from_team_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./team"
  }
}

run "resource_team_contact_create_slack_channel" {

  variables {
    name  = var.name
    team  = run.from_team_module.first_team.id
    type  = "slack"
    value = "#devs"
  }

  module {
    source = "./team_contact"
  }

  assert {
    condition = alltrue([
      can(opslevel_team_contact.test.id),
      can(opslevel_team_contact.test.name),
      can(opslevel_team_contact.test.team),
      can(opslevel_team_contact.test.type),
      can(opslevel_team_contact.test.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_team_contact.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_team_contact.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_team_contact.test.name,
    )
  }

  assert {
    condition = opslevel_team_contact.test.team == var.team
    error_message = format(
      "expected '%v' but got '%v'",
      var.team,
      opslevel_team_contact.test.team,
    )
  }

  assert {
    condition = opslevel_team_contact.test.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.test.type,
    )
  }

  assert {
    condition = opslevel_team_contact.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.test.value,
    )
  }

}

run "resource_team_contact_update_slack_handle" {

  variables {
    name  = var.name
    team  = run.from_team_module.first_team.id
    type  = "slack_handle"
    value = "@platform"
  }

  module {
    source = "./team_contact"
  }

  assert {
    condition = opslevel_team_contact.test.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.test.type,
    )
  }

  assert {
    condition = opslevel_team_contact.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.test.value,
    )
  }

}

run "resource_team_contact_update_email" {

  variables {
    name  = var.name
    team  = run.from_team_module.first_team.id
    type  = "email"
    value = "test@example.com"
  }

  module {
    source = "./team_contact"
  }

  assert {
    condition = opslevel_team_contact.test.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.test.type,
    )
  }

  assert {
    condition = opslevel_team_contact.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.test.value,
    )
  }

}

run "resource_team_contact_update_github" {

  variables {
    name  = var.name
    team  = run.from_team_module.first_team.id
    type  = "github"
    value = "opslevel"
  }

  module {
    source = "./team_contact"
  }

  assert {
    condition = opslevel_team_contact.test.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.test.type,
    )
  }

  assert {
    condition = opslevel_team_contact.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.test.value,
    )
  }

}

run "resource_team_contact_update_web" {

  variables {
    name  = var.name
    team  = run.from_team_module.first_team.id
    type  = "web"
    value = "https://platform.opslevel.com"
  }

  module {
    source = "./team_contact"
  }

  assert {
    condition = opslevel_team_contact.test.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.test.type,
    )
  }

  assert {
    condition = opslevel_team_contact.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.test.value,
    )
  }

}
