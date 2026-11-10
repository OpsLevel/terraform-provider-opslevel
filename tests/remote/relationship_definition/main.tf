resource "opslevel_component_type" "this" {
  name        = var.component_type_name
  alias       = var.component_type_alias
  description = "Component type for the relationship definition tests"
  properties  = {}

  icon = {
    color = "#f759ab"
    name  = "PhDeviceMobile"
  }
}

locals {
  # The same component type named three ways, chosen by var.component_type_reference.
  component_type = {
    alias       = opslevel_component_type.this.alias
    id          = opslevel_component_type.this.id
    upper_alias = upper(opslevel_component_type.this.alias)
  }
}

resource "opslevel_relationship_definition" "this" {
  name               = var.name
  alias              = var.alias
  component_type     = local.component_type[var.component_type_reference]
  allowed_categories = ["infrastructure"]
  management_rules = [{
    source_property      = "alias"
    operator             = "EQUALS"
    target_category      = "infrastructure"
    target_property      = "tag"
    target_tag_key       = "business_unit"
    target_tag_operation = "equals"
  }]
}
