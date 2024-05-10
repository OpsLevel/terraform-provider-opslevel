data "opslevel_integrations" "all" {}

data "opslevel_integration" "first_integration_by_id" {
  filter {
    field = "id"
    value = data.opslevel_integrations.all.integrations[0].id
  }
}

data "opslevel_integration" "first_integration_by_name" {
  filter {
    field = "name"
    value = data.opslevel_integrations.all.integrations[0].name
  }
}

