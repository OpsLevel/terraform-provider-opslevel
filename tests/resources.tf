# Domain resources

resource "opslevel_domain" "fancy" {
  name        = "Example"
  description = "The whole app in one monolith"
  owner       = "Developers"
  note        = "This is an example"
}

# Infrastructure resources

resource "opslevel_infrastructure" "small_infra" {
  data = jsonencode({
    name = "small-query"
  })
  owner  = var.test_id
  schema = "Small Database"
}


resource "opslevel_infrastructure" "big_infra" {
  aliases = ["big-infra"]
  data = jsonencode({
    name                = "big-query"
    external_id         = 1234
    replica             = true
    publicly_accessible = false
    storage_size = {
      unit  = "GB"
      value = 700
    }
  })
  owner = "Z2lkOi8vukla122kljf"
  provider_data = {
    account = "dev"
    name    = "google cloud"
    type    = "BigQuery"
    url     = "https://console.cloud.google.com/"
  }
  schema = "Big Database"
}

# Rubric Category resources

resource "opslevel_rubric_category" "mock_category" {
  name = "Mock Category"
}

# Rubric Level resources

resource "opslevel_rubric_level" "big" {
  description = "big rubric description"
  index       = 5
  name        = "big rubric level"
}

resource "opslevel_rubric_level" "small" {
  name = "small rubric level"
}

# Secret resources

resource "opslevel_secret" "mock_secret" {
  alias = "secret-alias"
  value = "too_many_passwords"
  owner = "Developers"
}

# Service resources

resource "opslevel_service" "big" {
  aliases                       = ["service-1", "service-2"]
  api_document_path             = "api/doc/path.yaml"
  description                   = "Service Description"
  framework                     = "Service Framework"
  language                      = "Service Language"
  lifecycle_alias               = "alpha"
  name                          = "Big Service"
  owner                         = "team-alias"
  preferred_api_document_source = "PULL"
  product                       = "Mock Product"
  tags                          = ["key1:value1", "key2:value2"]
  tier_alias                    = "Service Tier"
}

resource "opslevel_service" "small" {
  name = "Small Service"
}

# Scorecard resources

resource "opslevel_scorecard" "big" {
  affects_overall_service_levels = false
  description                    = "This is a big scorecard"
  filter_id                      = var.test_id
  name                           = "Big Scorecard"
  owner_id                       = var.test_id
}

resource "opslevel_scorecard" "small" {
  affects_overall_service_levels = true
  name                           = "Small Scorecard"
  owner_id                       = var.test_id
}

# User resources

resource "opslevel_user" "mock_user" {
  name  = "Mock User"
  email = "mock_user@mock.com"
  role  = "user"
}

resource "opslevel_user" "mock_user_no_role" {
  name  = "Mock User"
  email = "mock_user@mock.com"
}

resource "opslevel_user" "mock_user_admin" {
  name  = "Mock User"
  email = "mock_user@mock.com"
  role  = "admin"
}

# Webhook Action resources

resource "opslevel_webhook_action" "mock" {
  description = "Pages the On Call"
  headers = {
    accept        = "application/vnd.pagerduty+json;version=2"
    authorization = "Token token=XXXXXXXXXXXXXX"
    content-type  = "application/json"
    from          = "foo@opslevel.com"
  }
  method  = "POST"
  name    = "Small Webhook Action"
  payload = <<EOT
{
    "incident":
    {
        "type": "incident",
        "title": "{{manualInputs.IncidentTitle}}",
        "service": {
        "id": "{{ service | tag_value: 'pd_id' }}",
        "type": "service_reference"
        },
        "body": {
        "type": "incident_body",
        "details": "Incident triggered from OpsLevel by {{user.name}} with the email {{user.email}}. {{manualInputs.IncidentDescription}}"
        }
    }
}
  EOT
  url     = "https://api.pagerduty.com/incidents"
}

# Checks

# Check Manual

resource "opslevel_check_manual" "example" {
  name      = "foo"
  enable_on = "2022-05-23T14:14:18.782000Z"
  category  = var.test_id
  level     = var.test_id
  owner     = var.test_id
  filter    = var.test_id
  update_frequency = {
    starting_date = "2020-02-12T06:36:13Z"
    time_scale    = "week"
    value         = 1
  }
  update_requires_comment = false
  notes                   = "Optional additional info on why this check is run or how to fix it"
}

# Repo Search

resource "opslevel_check_git_branch_protection" "example" {
  name      = "foo"
  enable_on = "2022-05-23T14:14:18.782000Z"
  category  = var.test_id
  level     = var.test_id
  owner     = var.test_id
  filter    = var.test_id
}

# Repo Integrated

resource "opslevel_check_repository_integrated" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
}

# Repo Grep

resource "opslevel_check_repository_grep" "example" {
  name             = "foo"
  enabled          = true
  category         = var.test_id
  level            = var.test_id
  owner            = var.test_id
  filter           = var.test_id
  directory_search = false
  filepaths        = ["/src", "/tests"]
  file_contents_predicate = {
    type  = "contains"
    value = "**/hello.go"
  }
}

# Has Documentation

resource "opslevel_check_has_documentation" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id

  document_type    = "api"
  document_subtype = "openapi"
}

# Check Alert Source

resource "opslevel_check_alert_source_usage" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id

  alert_type = "pagerduty" # one of: "pagerduty", "datadog", "opsgenie"
  alert_name_predicate = {
    type  = "contains"
    value = "dev"
  }
}

# Check Recent Deploy

resource "opslevel_check_has_recent_deploy" "example" {
  name     = "foo"
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
  days     = 14
}

# Repo File

resource "opslevel_check_repository_file" "example" {
  name             = "foo"
  enabled          = true
  category         = var.test_id
  level            = var.test_id
  owner            = var.test_id
  filter           = var.test_id
  directory_search = false
  filepaths        = ["/src", "/tests"]
  file_contents_predicate = {
    type  = "equals"
    value = "import shim"
  }
  use_absolute_root = false
}

# Check Repo Search

resource "opslevel_check_repository_search" "example" {
  name            = "foo"
  enabled         = true
  category        = var.test_id
  level           = var.test_id
  owner           = var.test_id
  filter          = var.test_id
  file_extensions = ["sbt", "py"]
  file_contents_predicate = {
    type  = "contains"
    value = "postgres"
  }
}

# Check Service Configuration

resource "opslevel_check_service_configuration" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
}

# Check Service Dependency

resource "opslevel_check_service_dependency" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
}

# Check Service Ownership

resource "opslevel_check_service_ownership" "example" {
  name                   = "foo"
  enabled                = true
  category               = var.test_id
  level                  = var.test_id
  owner                  = var.test_id
  filter                 = var.test_id
  notes                  = "Optional additional info on why this check is run or how to fix it"
  require_contact_method = true
  contact_method         = "ANY"
  tag_key                = "team"
  tag_predicate = {
    type  = "equals"
    value = "frontend"
  }
}

# Check Tag Defined

resource "opslevel_check_tag_defined" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
  tag_key  = "environment"
  tag_predicate = {
    type  = "contains"
    value = "dev"
  }
}

# Check Tool Usage

resource "opslevel_check_tool_usage" "example" {
  name          = "foo"
  enabled       = true
  category      = var.test_id
  level         = var.test_id
  owner         = var.test_id
  filter        = var.test_id
  tool_category = "metrics"
  tool_name_predicate = {
    type  = "equals"
    value = "datadog"
  }
  environment_predicate = {
    type  = "equals"
    value = "production"
  }
}

# Check Custom Event

resource "opslevel_check_custom_event" "example" {
  name              = "foo"
  enabled           = true
  category          = var.test_id
  level             = var.test_id
  owner             = var.test_id
  filter            = var.test_id
  integration       = var.test_id
  pass_pending      = true
  service_selector  = ".messages[] | .incident.service.id"
  success_condition = ".messages[] |   select(.incident.service.id == $ctx.alias) | .incident.status == \"resolved\""
  message           = <<-EOT
  {% if check.passed %}
    ### Check passed
  {% else %}
    ### Check failed
    service **{{ data.messages[ctx.index].incident.service.id }}** has an unresolved incident.
  {% endif %}
  OpsLevel note: here you can fill in more details about this check. You can even include `data` from the payload, `params` specified in the URL and context `ctx` such as the service alias for the current evaluation.
  EOT
  notes             = "Optional additional info on why this check is run or how to fix it"
}
# Check Service Property

resource "opslevel_check_service_property" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
  property = "language"
  predicate = {
    type  = "equals"
    value = "python"
  }
}
