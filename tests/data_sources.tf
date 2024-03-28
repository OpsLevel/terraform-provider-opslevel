# Domain data sources

data "opslevel_domain" "mock_domain" {
  identifier = "example"
}

data "opslevel_domains" "all" {}

# Filter data sources

data "opslevel_filter" "name_filter" {
  filter {
    field = "name"
    value = "name-value"
  }
}

data "opslevel_filter" "id_filter" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }
}

data "opslevel_filter" "mock_filter" {
  filter {
    field = "name"
    value = "stuff"
  }
}

data "opslevel_integration" "name_filter" {
  filter {
    field = "name"
    value = "My Integration I Got By Filtering Name"
  }
}

data "opslevel_integration" "id_filter" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvSW50ZWdyYXRpb25zOjpTbGFja0ludGVncmF0aW9uLzI2"
  }
}

data "opslevel_rubric_category" "name_filter" {
  filter {
    field = "name"
    value = "name-value"
  }
}

data "opslevel_rubric_category" "id_filter" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }
}

# PropertyDefinition data sources

data "opslevel_property_definition" "mock_property_definition" {
  identifier = "mock-property_definition-alias"
}

# Service data sources

data "opslevel_service" "mock_service_with_alias" {
  alias = "mock-service-alias"
}

data "opslevel_service" "mock_service_with_id" {
  id = "Z2lkOi8vmock123"
}

# rubric Level data sources

data "opslevel_rubric_level" "alias_filter" {
  filter {
    field = "alias"
    value = "alias-value"
  }
}

data "opslevel_rubric_level" "id_filter" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }
}

data "opslevel_rubric_level" "index_filter" {
  filter {
    field = "index"
    value = 123
  }
}

data "opslevel_rubric_level" "name_filter" {
  filter {
    field = "name"
    value = "name-value"
  }
}

# lifecycle

data "opslevel_lifecycle" "alias_filter" {
  filter {
    field = "alias"
    value = "generally_available"
  }
}

data "opslevel_lifecycle" "id_filter" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvTGlmZWN5Y2xlLzQ"
  }
}

data "opslevel_lifecycle" "index_filter" {
  filter {
    field = "index"
    value = 123
  }
}

data "opslevel_lifecycle" "name_filter" {
  filter {
    field = "name"
    value = "Generally Available"
  }
}

# Scorecard data sources

data "opslevel_scorecard" "mock_scorecard" {
  identifier = "mock-scorecard-alias"
}

# System data sources

data "opslevel_system" "mock_system" {
  identifier = "my_system"
}

# Tier data sources

data "opslevel_tier" "mock_tier" {
  filter {
    field = "alias"
    value = ""
  }
}

data "opslevel_tier" "alias_filter" {
  filter {
    field = "alias"
    value = "alias-value"
  }
}

data "opslevel_tier" "id_filter" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }
}

data "opslevel_tier" "index_filter" {
  filter {
    field = "index"
    value = 123
  }
}

data "opslevel_tier" "name_filter" {
  filter {
    field = "name"
    value = "name-value"
  }
}

# Webhook Action data sources

data "opslevel_webhook_action" "mock_webhook_action" {
  identifier = "mock-webhook-action-alias"
}

# User data sources

data "opslevel_user" "mock_user" {
  identifier = "mock-user-name@example.com"
}
