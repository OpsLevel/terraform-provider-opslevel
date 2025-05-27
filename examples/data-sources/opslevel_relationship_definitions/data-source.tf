data "opslevel_relationship_definitions" "all" {
}

output "relationship_definitions" {
  value = data.opslevel_relationship_definitions.all
}

# Example of filtering relationship definitions by component type
output "service_relationships" {
  value = [
    for rd in data.opslevel_relationship_definitions.all.all :
    rd if rd.component_type == "Z2lkOi8vb3BzbGV2ZWwvUmVsYXRpb25zaGlwRGVmaW5pdGlvbi80Mg"
  ]
}
