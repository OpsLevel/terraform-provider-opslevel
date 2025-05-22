resource "opslevel_relationship_definition" "example" {
  name           = "Example Relationship"
  alias          = "example_relationship"
  description    = "An example relationship definition"
  component_type = "service"  # This should be a valid component type alias from your OpsLevel account
  allowed_types  = ["service", "library", "team"]  # Valid types this relationship can target, component alias or 'team'
}
