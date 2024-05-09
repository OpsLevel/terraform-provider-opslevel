# TODO: PropertyDefinition tests works on orange. Need to add to PAT acct.
#run "datasource_property_definitions_all" {
#
#  assert {
#    condition     = length(data.opslevel_property_definitions.all.property_definitions) > 0
#    error_message = "zero property_definitions found in data.opslevel_property_definitions"
#  }
#
#  assert {
#    condition = alltrue([
#      can(data.opslevel_property_definitions.all.property_definitions[0].id),
#      can(data.opslevel_property_definitions.all.property_definitions[0].name),
#    ])
#    error_message = "cannot set all expected property_definition datasource fields"
#  }
#
#}
#
#run "datasource_property_definition_first" {
#
#  assert {
#    condition     = data.opslevel_property_definition.first_property_definition_by_id.id == data.opslevel_property_definitions.all.property_definitions[0].id
#    error_message = "wrong ID on opslevel_property_definition"
#  }
#
#}
