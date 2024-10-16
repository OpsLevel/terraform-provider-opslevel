variables {
  resource_name = "opslevel_team"

  # required fields
  name = "TF Test Team"

  # optional fields
  aliases          = ["test_team_foo_bar_baz"]
  parent           = null # sourced from module
  responsibilities = "Team responsibilities"
  members          = [] # sourced from module
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_lifecycles.all,
      data.opslevel_teams.all,
      data.opslevel_users.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_team_create_with_all_fields" {

  variables {
    members = [
      {
        email = run.from_data_module.all_users.users[0].email
        role  = "manager"
      },
      {
        email = run.from_data_module.all_users.users[1].email
        role  = "contributor"
      },
    ]
    parent = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = alltrue([
      can(opslevel_team.this.aliases),
      can(opslevel_team.this.id),
      can(opslevel_team.this.member),
      can(opslevel_team.this.name),
      can(opslevel_team.this.parent),
      can(opslevel_team.this.responsibilities),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_team.this.aliases == toset(var.aliases)
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_team.this.aliases,
    )
  }

  assert {
    condition     = startswith(opslevel_team.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = length(opslevel_team.this.member) == length(var.members)
    error_message = format(
      "expected '%v' but got '%v'",
      var.members,
      opslevel_team.this.member,
    )
  }

  assert {
    condition = opslevel_team.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_team.this.name,
    )
  }

  assert {
    condition = opslevel_team.this.parent == var.parent
    error_message = format(
      "expected '%v' but got '%v'",
      var.parent,
      opslevel_team.this.parent,
    )
  }

  assert {
    condition = opslevel_team.this.responsibilities == var.responsibilities
    error_message = format(
      "expected '%v' but got '%v'",
      var.responsibilities,
      opslevel_team.this.responsibilities,
    )
  }

}

run "resource_team_unset_optional_fields" {

  variables {
    aliases          = null
    parent           = null
    responsibilities = null
    members          = []
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition     = opslevel_team.this.aliases == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = length(opslevel_team.this.member) == 0
    error_message = format(
      "expected 'opslevel_team.this.member' to be empty but got '%v'",
      opslevel_team.this.member,
    )
  }

  assert {
    condition     = opslevel_team.this.parent == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team.this.responsibilities == null
    error_message = var.error_expected_null_field
  }

}

run "delete_team_outside_of_terraform" {

  variables {
    resource_id   = run.resource_team_create_with_all_fields.this.id
    resource_type = "team"
  }

  module {
    source = "./provisioner"
  }
}

run "resource_team_create_with_required_fields" {

  variables {
    aliases          = null
    parent           = null
    responsibilities = null
    members          = []
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = run.resource_team_create_with_all_fields.this.id != opslevel_team.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_team_create_with_all_fields.this.id,
      opslevel_team.this.id,
    )
  }

  assert {
    condition     = opslevel_team.this.aliases == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = startswith(opslevel_team.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = length(opslevel_team.this.member) == 0
    error_message = format(
      "expected 'opslevel_team.this.member' to be empty but got '%v'",
      opslevel_team.this.member,
    )
  }

  assert {
    condition = opslevel_team.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_team.this.name,
    )
  }

  assert {
    condition     = opslevel_team.this.parent == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_team.this.responsibilities == null
    error_message = var.error_expected_null_field
  }

}

run "resource_team_set_all_fields" {

  variables {
    members = [
      {
        email = run.from_data_module.all_users.users[0].email
        role  = "manager"
      },
      {
        email = run.from_data_module.all_users.users[1].email
        role  = "contributor"
      },
    ]
    parent = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = opslevel_team.this.aliases == toset(var.aliases)
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_team.this.aliases,
    )
  }

  assert {
    condition = length(opslevel_team.this.member) == length(var.members)
    error_message = format(
      "expected '%v' but got '%v'",
      var.members,
      opslevel_team.this.member,
    )
  }

  assert {
    condition = opslevel_team.this.parent == var.parent
    error_message = format(
      "expected '%v' but got '%v'",
      var.parent,
      opslevel_team.this.parent,
    )
  }

  assert {
    condition = opslevel_team.this.responsibilities == var.responsibilities
    error_message = format(
      "expected '%v' but got '%v'",
      var.responsibilities,
      opslevel_team.this.responsibilities,
    )
  }

}
