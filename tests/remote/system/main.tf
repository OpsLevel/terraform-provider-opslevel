data "opslevel_systems" "all" {}

data "opslevel_system" "first_system_by_alias" {
  identifier = data.opslevel_systems.all.systems[0].aliases[0]
}

data "opslevel_system" "first_system_by_id" {
  identifier = data.opslevel_systems.all.systems[0].id
}

resource "opslevel_system" "test" {
  description = var.description
  domain      = var.domain_id
  name        = var.name
  note        = var.note
  owner       = var.owner_id
}
