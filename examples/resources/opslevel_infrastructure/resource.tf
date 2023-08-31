data "opslevel_group" "foo" {
    identifier = "foo"
}

// Minimum example
resource "opslevel_infrastructure" "example_1" {
    schema = "Database"
    provider_data {
        account = "dev"
    }
    data = jsonencode({
        name = "my-database"
    })
}

// Detailed example
resource "opslevel_infrastructure" "example_2" {
    schema = "Database"
    owner = data.opslevel_group.foo.id
    provider_data {
        account = "dev"
        name = "google cloud"
        type = "BigQuery"
        url = "https://console.cloud.google.com/..."
    }
    data = jsonencode({
        name = "big-query"
        zone = "us-east-1"
        engine = "bigquery"
        engine_version = "1.28.0"
        endpoint = "https://console.cloud.google.com/..."
        replica = false
        publicly_accessible = false
        storage_size = {
            unit = "GB"
            value = 700
        }
        storage_type = "gp3"
        storage_iops = {
            unit = "per second"
            value = 12000
        }
    })
}
