resource "opslevel_integration_aws" "this" {
  name                    = var.name
  iam_role                = var.iam_role
  external_id             = var.external_id
  ownership_tag_overrides = var.ownership_tag_overrides
  ownership_tag_keys      = var.ownership_tag_keys
  region_override         = var.region_override
}
