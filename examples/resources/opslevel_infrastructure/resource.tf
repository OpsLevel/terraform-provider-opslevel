data "opslevel_group" "foo" {
    identifier = "foo"
}

resource "opslevel_infrastructure" "example" {
    type = "BigQuery"
    schema = "Database"
    owner = data.opslevel_group.foo.id
    provider_data {
        account_name = "dev"
        external_url = "https://console.cloud.google.com/..."
        provider_name = "google cloud"
    }
    data = jsonencode({
        region = "us-east-1"
        engine = "bigquery"
        engine_version = "1.28.0"
        endpoint = "https://console.cloud.google.com/..."
        replica = false
        storage_size = "300 GB"
        publicly_accessible = false
    })
}