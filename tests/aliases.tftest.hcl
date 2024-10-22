variables {
  resource_name = "opslevel_alias"
  # required fields
  # optional fields
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_domains.all,
      data.opslevel_services.all,
      data.opslevel_systems.all,
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_create_aliases" {
  variables {
    resource_type       = "domain"
    resource_identifier = run.from_data_module.first_domain.id
    aliases             = toset(["one", "two", "three"])
  }

  module {
    source = "./opslevel_modules/modules/aliases"
  }

  assert {
    condition = alltrue([
      can(opslevel_alias.this.aliases),
      can(opslevel_alias.this.id),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_alias.this.aliases == var.aliases
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_alias.this.aliases,
    )
  }
}

run "resource_modify_managed_aliases" {
  variables {
    resource_type       = "domain"
    resource_identifier = run.from_data_module.first_domain.id
    aliases = toset(["one", "four", "three"])
  }

  module {
    source = "./opslevel_modules/modules/aliases"
  }

  assert {
    condition = opslevel_alias.this.aliases == var.aliases
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_alias.this.aliases,
    )
  }
}

run "delete_delete_alias_outside_of_terraform" {

  variables {
    command = "delete alias -t domain four"
  }

  module {
    source = "./cli"
  }
}

run "resource_ensure_managed_aliases" {
  variables {
    resource_type       = "domain"
    resource_identifier = run.from_data_module.first_domain.id
    aliases = toset(["one", "four", "three"])
  }

  module {
    source = "./opslevel_modules/modules/aliases"
  }

  assert {
    condition = opslevel_alias.this.aliases == var.aliases
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_alias.this.aliases,
    )
  }
}