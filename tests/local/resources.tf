# Domain resources

resource "opslevel_domain" "fancy" {
  name        = "Example"
  description = "The whole app in one monolith"
  owner       = var.test_id
  note        = "This is an example"
}

# Filter resources

resource "opslevel_filter" "small" {
  name = "Blank Filter"
}

resource "opslevel_filter" "big" {
  connective = var.connective_enum
  name       = "Big Filter"
  predicate {
    key   = var.predicate_key_enum
    type  = var.predicate_type_enum
    value = "1"
  }
  predicate {
    key   = "lifecycle_index"
    type  = "less_than_or_equal_to"
    value = "1"
  }
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

# Integration AWS resources

resource "opslevel_integration_aws" "example" {
  external_id             = "mock-external-id"
  iam_role                = "arn:aws:ecr:us-east-1:mock-iam-role"
  name                    = "dev"
  ownership_tag_overrides = true
  ownership_tag_keys      = ["owner", "team", "group"]
}

# Integration Azure Resources resources

resource "opslevel_integration_azure_resources" "example" {
  client_id               = "XXX_CLIENT_ID_XXX"
  client_secret           = "XXX_CLIENT_SECRET_XXX"
  name                    = "dev"
  ownership_tag_keys      = ["owner", "team", "group"]
  ownership_tag_overrides = true
  subscription_id         = "01234567-0123-0123-0123-012345678901"
  tenant_id               = "98765432-9876-9876-9876-987654321098"
}

# Property Assignment

resource "opslevel_property_assignment" "color_picker_using_aliases" {
  owner      = "some_service"
  definition = "some_definition"
  value      = jsonencode("green")
}

resource "opslevel_property_assignment" "color_picker_using_id" {
  owner      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MzMyOQ"
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8zNzU"
  value      = jsonencode({ "hello" : "world", "key" : null })
}

# Property Definition

resource "opslevel_property_definition" "color_picker" {
  name = "Color Picker"
  schema = jsonencode({
    "type" : "string",
    "enum" : [
      "red",
      "green",
      "blue",
    ]
  })
  allowed_in_config_files = false
  property_display_status = "visible"
}

# Repository resources

resource "opslevel_repository" "with_alias" {
  identifier = "github.com:rocktavious/autopilot"
}

resource "opslevel_repository" "with_id" {
  identifier = var.test_id
  owner      = var.test_id
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

# Service Repository resources

resource "opslevel_service_repository" "with_alias" {
  base_directory   = "home/opslevel"
  name             = "Service Repo Name"
  repository_alias = "github.com:OpsLevel/terraform-provider-opslevel"
  service_alias    = "service-1"
}

resource "opslevel_service_repository" "with_id" {
  repository = var.test_id
  service    = var.test_id
}

# Service Dependency resources

resource "opslevel_service_dependency" "with_alias" {
  depends_upon = var.test_id
  service      = var.test_id
}

resource "opslevel_service_dependency" "with_id" {
  depends_upon = var.test_id
  note         = <<-EOT
    This is an example of notes on a service dependency
  EOT
  service      = var.test_id
}

# Service Tool resources

resource "opslevel_service_tool" "with_id" {
  category = "metrics"
  name     = "Datadog test"
  service  = var.test_id
  url      = "https://datadoghq.com"
}

resource "opslevel_service_tool" "with_alias" {
  category      = "logs"
  environment   = "Production"
  name          = "Datadog test"
  service_alias = "mock service"
  url           = "https://datadoghq.com"
}

# Service Tag resources

resource "opslevel_service_tag" "using_service_id" {
  key     = "hello_with_id"
  value   = "world_with_id"
  service = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQxMg"
}

resource "opslevel_service_tag" "using_service_alias" {
  key           = "hello_with_alias"
  value         = "world_with_alias"
  service_alias = "cart"
}

# System resources

resource "opslevel_system" "big" {
  description = "It's a big system"
  domain      = var.test_id
  name        = "Big System"
  note        = "Note on System"
  owner       = var.test_id
}

resource "opslevel_system" "small" {
  name = "Small System"
}

# Team resources

resource "opslevel_team" "big" {
  aliases          = ["the_big_team", "big_team"]
  name             = "The Big Team"
  parent           = "small_team"
  responsibilities = "This is a big team"

  member {
    email = "alice@opslevel.com"
    role  = "manager"
  }

  member {
    email = "bob@opslevel.com"
    role  = "contributor"
  }
}

resource "opslevel_team" "small" {
  name = "Small Team"
}

# Team Contact

resource "opslevel_team_contact" "tc_1" {
  name  = "Internal Slack Channel"
  team  = "team_platform_3"
  type  = "slack"
  value = "#platform-3"
}

resource "opslevel_team_contact" "tc_2" {
  name  = "Team Email Internal"
  team  = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQxMg"
  type  = "email"
  value = "team-platform-3-3-3@opslevel.com"
}

# Team Tag resources

resource "opslevel_team_tag" "using_team_id" {
  key   = "hello_with_id"
  value = "world_with_id"
  team  = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQxMg"
}

resource "opslevel_team_tag" "using_team_alias" {
  key        = "hello_with_alias"
  value      = "world_with_alias"
  team_alias = "team_platform_3"
}

# Trigger Definition resources

resource "opslevel_trigger_definition" "big" {
  access_control           = "everyone"
  action                   = var.test_id
  description              = "Pages the On Call"
  entity_type              = "SERVICE"
  extended_team_access     = ["team_1", "team_2"]
  filter                   = var.test_id
  manual_inputs_definition = <<EOT
---
version: 1
inputs:
  - identifier: IncidentTitle
    displayName: Title
    description: Title of the incident to trigger
    type: text_input
    required: true
    maxLength: 60
    defaultValue: Service Incident Manual Trigger
  - identifier: IncidentDescription
    displayName: Incident Description
    description: The description of the incident
    type: text_area
    required: true
  EOT
  response_template        = <<EOT
{% if response.status >= 200 and response.status < 300 %}
## Congratulations!
Your request for {{ service.name }} has succeeded. See the incident here: {{response.body.incident.html_url}}
{% else %}
## Oops something went wrong!
Please contact [{{ action_owner.name }}]({{ action_owner.href }}) for more help.
{% endif %}
  EOT
  name                     = "Big Trigger Definition"
  owner                    = var.test_id
  published                = false
}

resource "opslevel_trigger_definition" "small" {
  access_control = "everyone"
  action         = var.test_id
  name           = "Small Trigger Definition"
  owner          = var.test_id
  published      = true
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


# Repo Search

resource "opslevel_check_git_branch_protection" "example" {
  name      = "foo"
  enable_on = "2022-05-23T14:14:18.782000Z"
  category  = var.test_id
  level     = var.test_id
  owner     = var.test_id
  filter    = var.test_id
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

# Check Recent Deploy

resource "opslevel_check_has_recent_deploy" "example" {
  name     = "foo"
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
  days     = 14
}

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

# Check Repo File

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

# Repo Integrated

resource "opslevel_check_repository_integrated" "example" {
  name     = "foo"
  enabled  = true
  category = var.test_id
  level    = var.test_id
  owner    = var.test_id
  filter   = var.test_id
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
  contact_method         = "any"
  tag_key                = "team"
  tag_predicate = {
    type  = "equals"
    value = "frontend"
  }
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

# Tag Resource

resource "opslevel_tag" "example" {
  resource_type       = "Service"
  resource_identifier = "test-service"

  key   = "yacht"
  value = "racing"
}

