resource "opslevel_integration_google_cloud" "this" {
  client_email            = var.client_email
  name                    = var.name
  ownership_tag_keys      = var.ownership_tag_keys
  ownership_tag_overrides = var.ownership_tag_overrides
  private_key             = var.private_key
}
