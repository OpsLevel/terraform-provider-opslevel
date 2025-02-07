resource "opslevel_component_type" "example" {
  name        = "Mobile Apps"
  alias       = "mobile-app"
  description = "mobile app component type"
  properties = {
    os = {
      name                    = "Operation System"
      description             = "The Operating System this app is deployed on"
      allowed_in_config_files = "true"
      locked_status           = "unlocked"
      schema = jsonencode({
        "enum" : [
          "ios",
          "android",
          "both"
        ],
        "type" : "string"
      })
    }
    version = {
      name = "Release Version"
      schema = jsonencode({
        "type" : "string",
        "pattern" : "^([0-9]+).([0-9]+).([0-9]+)"
      })
    }
    bundle-id = {
      name = "Bundle ID"
      schema = jsonencode({
        "type" : "string"
      })
    }
    update-strategy = {
      name          = "Update Strategy"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "code-push",
          "firebase",
          "liftoff"
        ],
        "type" : "string"
      })
    }
    crash-reporting = {
      name                    = "Crash Reporting"
      allowed_in_config_files = "true"
      locked_status           = "unlocked"
      schema = jsonencode({
        "enum" : [
          "sentry",
          "firebase",
          "crashlytics"
        ],
        "type" : "string"
      })
    }
  }
}
