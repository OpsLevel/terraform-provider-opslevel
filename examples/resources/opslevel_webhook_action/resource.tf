resource "opslevel_webhook_action" "example" {
  name        = "Page The On Call"
  description = "Pages the On Call"
  url         = "https://api.pagerduty.com/incidents"
  method      = "POST"
  headers = {
    content-type  = "application/json"
    accept        = "application/vnd.pagerduty+json;version=2"
    authorization = "Token token=XXXXXXXXXXXXXX"
    from          = "john@opslevel.com"
  }
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
}