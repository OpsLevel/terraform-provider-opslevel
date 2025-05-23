resource "opslevel_component_type" "mobile" {
  name        = "Mobile App"
  alias       = "mobile-app"
  description = "A mobile app component type."
  icon = {
    color = "#f759ab"
    name  = "PhDeviceMobile"
  }
  properties = {
    os = {
      name                    = "Operation System"
      description             = "The operating system this app is deployed to."
      allowed_in_config_files = "true"
      locked_status           = "ui_locked"
      display_status          = "hidden"
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
      name          = "Release Version"
      locked_status = "unlocked"
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
  relationships = {
    services = {
      name                    = "Services Needed"
      description             = "The services this Mobile App depends on."
      allowed_types           = ["service"]
    }
  }
}

resource "opslevel_component_type" "lambda" {
  name        = "Lambda"
  alias       = "lambda"
  description = "A cloud function like AWS Lambda or Azure Functions."
  icon = {
    color = "#ffa940"
    name  = "PhFunction"
  }
  properties = {
    cloud = {
      name          = "Cloud Provider"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "aws",
          "azure",
          "gcp"
        ],
        "type" : "string"
      })
    }
    runtime = {
      name          = "Runtime"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "nodejs",
          "python",
          "java",
          "dotnet",
          "go",
          "ruby",
          "php"
        ],
        "type" : "string"
      })
    }
    release-strategy = {
      name          = "Release Strategy"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "manual",
          "automated",
        ],
        "type" : "string"
      })
    }
    concurrency = {
      name          = "Concurrency Limit"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "number"
      })
    }
    timeout = {
      name          = "Timeout"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "number"
      })
    }
    cold-start = {
      name          = "Cold Start Mitigation"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "boolean"
      })
    }
    dead-letter = {
      name          = "Dead Letter Queue"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "boolean"
      })
    }
  }
  relationships = {
    services = {
      name                    = "Services Needed"
      description             = "The services this Lambda depends on."
      allowed_types           = ["service"]
    }
    libraries = {
      name                    = "Libraries Needed"
      description             = "The libraries this Lambda depends on."
      allowed_types           = ["service"]
    }
  }
}

resource "opslevel_component_type" "library" {
  name        = "Library"
  alias       = "library"
  description = "A non-runtime based component type that has a release strategy instead of continuous deployment."
  icon = {
    color = "#36cfc9"
    name  = "PhBooks"
  }
  properties = {
    release-strategy = {
      name          = "Release Strategy"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "manual",
          "automated",
        ],
        "type" : "string"
      })
    }
    artifact-repo = {
      name          = "Artifact Repository"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "array",
        "items" : {
          "enum" : [
            "maven",
            "npm",
            "nuget",
            "pypi",
            "gem",
            "docker",
            "homebrew",
            "github",
            "gitlab",
            "bitbucket",
            "aws",
            "azure",
            "gcp",
            "other"
          ],
          "type" : "string"
        },
        "minItems" : 1,
        "uniqueItems" : true
      })
    }
    version-strategy = {
      name          = "Version Strategy"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "semver",
          "calver",
          "custom"
        ],
        "type" : "string"
      })
    }
    cadence = {
      name          = "Cadence"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "weekly",
          "monthly",
          "quarterly",
          "yearly"
        ],
        "type" : "string"
      })
    }
    license = {
      name = "License"
      schema = jsonencode({
        "enum" : [
          "mit",
          "bsl",
          "apache",
          "other"
        ],
        "type" : "string"
      })
    }
    coverage = {
      name = "Test Coverage"
      schema = jsonencode({
        "type" : "number"
      })
    }
  }
}

resource "opslevel_component_type" "vendor" {
  name        = "Vendor Tool"
  alias       = "vendor"
  description = "A third party vendor tool."
  icon = {
    color = "#bae637"
    name  = "PhTruck"
  }
  properties = {
    vendor-contact = {
      name          = "Vendor Contact Name"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "string"
      })
    }
    vendor-contact-email = {
      name          = "Vendor Contact Email"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "string",
        "format" : "email"
      })
    }
    privacy-policy = {
      name          = "Privacy Policy"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "string"
      })
    }
    terms-of-service = {
      name          = "Terms of Service"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "string"
      })
    }
    status = {
      name          = "Status"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "Evaluating",
          "Purchasing",
          "Active",
          "Decommissioning",
          "Inactive"
        ],
        "type" : "string"
      })
    }
    business-unit = {
      name          = "Business Unit"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "Admin",
          "HR",
          "Customer Success",
          "Engineering",
          "Marketing",
          "Sales",
          "Legal",
          "Security",
          "Other"
        ]
        "type" : "string"
      })
    }
    risk = {
      name          = "Risk"
      locked_status = "unlocked"
      schema = jsonencode({
        "enum" : [
          "none",
          "low",
          "medium",
          "high"
        ],
        "type" : "string"
      })
    }
    pii = {
      name          = "Stores PII"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "boolean"
      })
    }
    subprocessor = {
      name          = "Subprocessor"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "boolean"
      })
    }
    contract-value = {
      name          = "Contract Value"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "number"
      })
    }
  }
}

resource "opslevel_component_type" "ml-ai" {
  name        = "ML / AI"
  alias       = "ml-ai"
  description = "A machine learning or ai model."
  icon = {
    color = "#737373"
    name  = "PhAtom"
  }
  properties = {
    contact = {
      name          = "Contact Name"
      locked_status = "unlocked"
      schema = jsonencode({
        "type" : "string"
      })
    }
  }
}
