# Domain resources

resource "opslevel_domain" "fancy" {
  name        = "Example"
  description = "The whole app in one monolith"
  owner       = "Developers"
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
    key  = var.predicate_key_enum
    type = var.predicate_type_enum
  }
  predicate {
    case_insensitive = false
    case_sensitive   = true
    key              = "lifecycle_index"
    key_data         = "big_predicate"
    type             = "ends_with"
    value            = "1"
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
