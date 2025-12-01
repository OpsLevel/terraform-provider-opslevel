# Infrastructure resource relationship using
resource "opslevel_relationship_assignment" "basic" {
  source = opslevel_service.example.id        # Component ID
  target = opslevel_infrastructure.example.id # Infrastructure ID
  type   = "belongs_to"                       # One of: belongs_to, depends_on
}

# Component relationship using a relationship definition
resource "opslevel_relationship_definition" "example" {
  name               = "Example"
  alias              = "example"
  component_type     = "service"
  allowed_categories = ["default"]
  allowed_types      = ["service"]
}

resource "opslevel_relationship_assignment" "custom" {
  source     = opslevel_service.example.id
  target     = opslevel_service.other.id
  type       = "related_to"
  definition = opslevel_relationship_definition.example.id
}
