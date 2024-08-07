data "opslevel_team" "foo" {
  alias = "foo"
}

// Minimum example
resource "opslevel_infrastructure" "example_1" {
  schema = "Database"
  owner  = data.opslevel_team.foo.id
  provider_data = {
    account = "dev"
  }
  data = jsonencode({
    name = "my-database"
  })
}

// Detailed example
resource "opslevel_infrastructure" "example_2" {
  aliases = ["foo", "bar", "baz"]
  schema  = "Database"
  owner   = data.opslevel_team.foo.id
  provider_data = {
    account = "dev"
    name    = "google cloud"
    type    = "BigQuery"
    url     = "https://console.cloud.google.com/..."
  }
  data = jsonencode({
    name                = "big-query"
    external_id         = "example_2_1234"
    zone                = "us-east-1"
    engine              = "bigquery"
    engine_version      = "1.28.0"
    endpoint            = "https://console.cloud.google.com/..."
    replica             = false
    publicly_accessible = false
    storage_size = {
      unit  = "GB"
      value = 700
    }
    storage_type = "gp3"
    storage_iops = {
      unit  = "per second"
      value = 12000
    }
  })
}
