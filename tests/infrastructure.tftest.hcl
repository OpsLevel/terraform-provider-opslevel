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
  owner  = null # sourced from module
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

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_infrastructure_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    owner = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/infrastructure"
  }

  assert {
    condition = alltrue([
      can(opslevel_infrastructure.this.aliases),
      can(opslevel_infrastructure.this.data),
      can(opslevel_infrastructure.this.id),
      can(opslevel_infrastructure.this.owner),
      can(opslevel_infrastructure.this.provider_data),
      can(opslevel_infrastructure.this.schema),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_infrastructure.this.aliases == var.aliases
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_infrastructure.this.aliases,
    )
  }

  assert {
    condition = opslevel_infrastructure.this.data == var.data
    error_message = format(
      "expected '%v' but got '%v'",
      var.data,
      opslevel_infrastructure.this.data,
    )
  }

  assert {
    condition     = startswith(opslevel_infrastructure.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_infrastructure.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_infrastructure.this.owner,
    )
  }

  assert {
    condition = opslevel_infrastructure.this.provider_data == var.provider_data
    error_message = format(
      "expected '%v' but got '%v'",
      var.provider_data,
      opslevel_infrastructure.this.provider_data,
    )
  }

  assert {
    condition = opslevel_infrastructure.this.schema == var.schema
    error_message = format(
      "expected '%v' but got '%v'",
      var.schema,
      opslevel_infrastructure.this.schema,
    )
  }

}

run "resource_infrastructure_create_unset_optional_fields" {

  variables {
    # other fields from file scoped variables block
    aliases       = null
    owner         = run.from_data_module.first_team.id
    provider_data = null
  }

  module {
    source = "./opslevel_modules/modules/infrastructure"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.aliases == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.provider_data == null
    error_message = var.error_expected_null_field
  }

}
