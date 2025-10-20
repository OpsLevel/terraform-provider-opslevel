mock_resource "opslevel_component_type" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    name        = "API"
    alias       = "api"
    description = "An API component type"
    icon = {
      color = "#F59E0B"
      name  = "PhCloud"
    }
    properties = {
      "api_version" = {
        name                    = "API Version"
        description             = "The version of the API"
        allowed_in_config_files = true
        display_status          = "visible"
        locked_status           = "unlocked"
        schema                  = "{\"type\":\"string\"}"
      }
    }
    relationships = {}
  }
}

