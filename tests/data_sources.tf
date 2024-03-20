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

# PropertyDefinition data sources
data "opslevel_property_definition" "mock_property_definition" {
  identifier = "mock-property_definition-alias"
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
