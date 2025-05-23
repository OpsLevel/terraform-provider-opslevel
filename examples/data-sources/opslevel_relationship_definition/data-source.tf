data "opslevel_relationship_definition" "example" {
  identifier = "Z2lkOi8vb3BzbGV2ZWwvUmVsYXRpb25zaGlwRGVmaW5pdGlvbi80Mg" # Only supports ID for now
}

output "relationship_definition_id" {
  value = data.opslevel_relationship_definition.example.id
}
