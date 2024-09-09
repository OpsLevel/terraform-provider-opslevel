variables {
  resource_name = "opslevel_team"

  # required fields
  name = "TF Test Team"

  # optional fields
  aliases          = null
  parent           = null # sourced from module
  responsibilities = null
  members          = [] # sourced from module
}

run "resource_team_create_with_required_fields" {

  module {
    source = "./opslevel_modules/modules/team"
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

run "resource_team_set_aliases" {

  variables {
    aliases = ["test_team_123", "test_team_321"]
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

}

run "from_user_module" {
  command = plan

  variables {
    email = ""
  }

  module {
    source = "./opslevel_modules/modules/user"
  }
}

run "resource_team_set_members" {

  variables {
    members = [
      {
        email = run.from_user_module.all.users[0].email
        role  = "manager"
      },
      {
        email = run.from_user_module.all.users[1].email
        role  = "contributor"
      },
    ]
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = length(opslevel_team.this.member) == length(var.members)
    error_message = format(
      "expected '%v' but got '%v'",
      var.members,
      opslevel_team.this.member,
    )
  }

}

run "from_team_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./opslevel_modules/modules/team"
  }
}

run "resource_team_set_parent" {

  variables {
    parent = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = opslevel_team.this.parent == var.parent
    error_message = format(
      "expected '%v' but got '%v'",
      var.parent,
      opslevel_team.this.parent,
    )
  }
}

run "resource_team_set_responsibilities" {

  variables {
    responsibilities = "Team responsibilities"
  }

  module {
    source = "./opslevel_modules/modules/team"
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
