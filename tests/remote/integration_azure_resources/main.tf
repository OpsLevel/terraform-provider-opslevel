resource "opslevel_integration_azure_resources" "test" {
  client_id               = var.client_id
  client_secret           = var.client_secret
  name                    = var.name
  ownership_tag_keys      = var.ownership_tag_keys
  ownership_tag_overrides = var.ownership_tag_overrides
  subscription_id         = var.subscription_id
  tenant_id               = var.tenant_id
}

