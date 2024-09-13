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

run "from_team_module" {
  command = plan

  module {
    source = "./data/team"
  }
}

run "from_user_module" {
  command = plan

  module {
    source = "./data/user"
  }
}

run "resource_team_create_with_all_fields" {

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
    parent = run.from_team_module.first.id
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

run "resource_team_unset_aliases" {

  variables {
    aliases = null
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition     = opslevel_team.this.aliases == null
    error_message = var.error_expected_null_field
  }

}

run "resource_team_unset_members" {

  variables {
    members = []
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition = length(opslevel_team.this.member) == 0
    error_message = format(
      "expected 'opslevel_team.this.member' to be empty but got '%v'",
      opslevel_team.this.member,
    )
  }

}

run "resource_team_unset_parent" {

  variables {
    parent = null
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition     = opslevel_team.this.parent == null
    error_message = var.error_expected_null_field
  }

}

run "resource_team_unset_responsibilities" {

  variables {
    responsibilities = null
  }

  module {
    source = "./opslevel_modules/modules/team"
  }

  assert {
    condition     = opslevel_team.this.responsibilities == null
    error_message = var.error_expected_null_field
  }

}
