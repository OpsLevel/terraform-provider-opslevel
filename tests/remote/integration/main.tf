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


#resource "opslevel_integration" "test" {
#  allowed_in_config_files = var.allowed_in_config_files
#  description             = var.description
#  name                    = var.name
#  property_display_status = var.property_display_status
#  schema                  = var.schema
#  note                    = var.note
#}
