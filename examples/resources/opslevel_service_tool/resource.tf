data "opslevel_service" "foo" {
  alias = "foo"
}

resource "opslevel_service_tool" "foo_datadog" {
  service = data.opslevel_service.foo.id

  name = "Datadog"
  category = "metrics"
  url = "https://datadoghq.com"
  environment = "Production"
}

resource "opslevel_service_tool" "bar_datadog" {
  service_alias = "bar"

  name = "Datadog"
  category = "metrics"
  url = "https://datadoghq.com"
  environment = "Production"
}