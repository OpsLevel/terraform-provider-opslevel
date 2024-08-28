variables {
  resource_name = "opslevel_infrastructure"

  # required fields
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
  owner  = null
  schema = "Database"

  # optional fields
  aliases = ["foo", "bar", "baz"]
  provider_data = {
    account = "TF test acct"
    name    = "google cloud"
    type    = "BigQuery"
    url     = "https://console.cloud.google.com/..."
  }
}

run "from_team_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./team"
  }
}

run "resource_infrastructure_create_with_all_fields" {

  variables {
    aliases       = var.aliases
    data          = var.data
    owner         = run.from_team_module.first_team.id
    provider_data = var.provider_data
    schema        = var.schema
  }

  module {
    source = "./infrastructure"
  }

  assert {
    condition = alltrue([
      can(opslevel_infrastructure.test.aliases),
      can(opslevel_infrastructure.test.data),
      can(opslevel_infrastructure.test.id),
      can(opslevel_infrastructure.test.owner),
      can(opslevel_infrastructure.test.provider_data),
      can(opslevel_infrastructure.test.schema),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_infrastructure.test.aliases == var.aliases
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_infrastructure.test.aliases,
    )
  }

  assert {
    condition = opslevel_infrastructure.test.data == var.data
    error_message = format(
      "expected '%v' but got '%v'",
      var.data,
      opslevel_infrastructure.test.data,
    )
  }

  assert {
    condition     = startswith(opslevel_infrastructure.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_infrastructure.test.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_infrastructure.test.owner,
    )
  }

  assert {
    condition = opslevel_infrastructure.test.provider_data == var.provider_data
    error_message = format(
      "expected '%v' but got '%v'",
      var.provider_data,
      opslevel_infrastructure.test.provider_data,
    )
  }

  assert {
    condition = opslevel_infrastructure.test.schema == var.schema
    error_message = format(
      "expected '%v' but got '%v'",
      var.schema,
      opslevel_infrastructure.test.schema,
    )
  }

}
