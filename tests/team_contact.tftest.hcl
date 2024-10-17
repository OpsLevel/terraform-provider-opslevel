variables {
  resource_name = "opslevel_team_contact"

  # required fields
  name  = "TF Test Team Contact"
  team  = null # sourced from module
  type  = "slack"
  value = "#devs"

  # optional fields - none
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_teams.all,
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_team_contact_create_slack_channel" {

  variables {
    # other fields from file scoped variables block
    team = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/team/contact"
  }

  assert {
    condition = alltrue([
      can(opslevel_team_contact.this.id),
      can(opslevel_team_contact.this.name),
      can(opslevel_team_contact.this.team),
      can(opslevel_team_contact.this.type),
      can(opslevel_team_contact.this.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_team_contact.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_team_contact.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_team_contact.this.name,
    )
  }

  assert {
    condition = opslevel_team_contact.this.team == var.team
    error_message = format(
      "expected '%v' but got '%v'",
      var.team,
      opslevel_team_contact.this.team,
    )
  }

  assert {
    condition = opslevel_team_contact.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.this.type,
    )
  }

  assert {
    condition = opslevel_team_contact.this.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.this.value,
    )
  }

}

run "delete_team_contact_outside_of_terraform" {

  variables {
    resource_id   = run.resource_team_contact_create_slack_channel.this.id
    resource_type = "contact"
  }

  module {
    source = "./provisioner"
  }
}

run "resource_team_contact_create_slack_handle" {

  variables {
    # other fields from file scoped variables block
    team  = run.from_data_module.first_team.id
    type  = "slack_handle"
    value = "@platform"
  }

  module {
    source = "./opslevel_modules/modules/team/contact"
  }

  assert {
    condition = run.resource_team_contact_create_slack_channel.this.id != opslevel_team_contact.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_team_contact_create_slack_channel.this.id,
      opslevel_team_contact.this.id,
    )
  }

  assert {
    condition = opslevel_team_contact.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.this.type,
    )
  }

  assert {
    condition = opslevel_team_contact.this.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.this.value,
    )
  }

}

run "resource_team_contact_update_email" {

  variables {
    # other fields from file scoped variables block
    team  = run.from_data_module.first_team.id
    type  = "email"
    value = "test@example.com"
  }

  module {
    source = "./opslevel_modules/modules/team/contact"
  }

  assert {
    condition = opslevel_team_contact.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.this.type,
    )
  }

  assert {
    condition = opslevel_team_contact.this.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.this.value,
    )
  }

}

run "resource_team_contact_update_github" {

  variables {
    # other fields from file scoped variables block
    team  = run.from_data_module.first_team.id
    type  = "github"
    value = "opslevel"
  }

  module {
    source = "./opslevel_modules/modules/team/contact"
  }

  assert {
    condition = opslevel_team_contact.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.this.type,
    )
  }

  assert {
    condition = opslevel_team_contact.this.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.this.value,
    )
  }

}

run "resource_team_contact_update_web" {

  variables {
    # other fields from file scoped variables block
    team  = run.from_data_module.first_team.id
    type  = "web"
    value = "https://platform.opslevel.com"
  }

  module {
    source = "./opslevel_modules/modules/team/contact"
  }

  assert {
    condition = opslevel_team_contact.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_team_contact.this.type,
    )
  }

  assert {
    condition = opslevel_team_contact.this.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_team_contact.this.value,
    )
  }

}
