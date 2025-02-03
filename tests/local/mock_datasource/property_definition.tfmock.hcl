mock_data "opslevel_property_definition" {
  defaults = {
    allowed_in_config_files = true
    description             = "mock-property-definition-description"
    # id intentionally omitted - will be assigned a random string
    name                    = "mock-property-definition-name"
    property_display_status = "visible"
    locked_status           = "unlocked"
    # schema value result of running
    # 'jsonencode({ "$ref" : "#/$defs/MyProp", "$defs" : { "MyProp" : { "properties" : { "name" : { "type" : "string", "title" : "the new name", "description" : "The name of a friend", "default" : "alex", "examples" : ["joe", "lucy"] } }, "additionalProperties" : false, "type" : "object", "required" : ["name"] } } })'
    schema = "{\"$defs\":{\"MyProp\":{\"additionalProperties\":false,\"properties\":{\"name\":{\"default\":\"alex\",\"description\":\"The name of a friend\",\"examples\":[\"joe\",\"lucy\"],\"title\":\"the new name\",\"type\":\"string\"}},\"required\":[\"name\"],\"type\":\"object\"}},\"$ref\":\"#/$defs/MyProp\"}"
  }
}

