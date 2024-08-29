resource "opslevel_infrastructure" "test" {
  aliases       = var.aliases
  data          = var.data
  owner         = var.owner
  provider_data = var.provider_data
  schema        = var.schema
}
