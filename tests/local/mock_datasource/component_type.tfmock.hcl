mock_data "opslevel_component_type" {
  defaults = {
    name        = "Service"
    alias       = "service"
    description = "A service component type"
    icon = {
      color = "#3B82F6"
      name  = "PhCube"
    }
    properties = {
      "deployment_platform" = {
        name        = "Deployment Platform"
        description = "The platform used to deploy this service"
        schema = {
          type = "string"
        }
      }
    }
  }
}
