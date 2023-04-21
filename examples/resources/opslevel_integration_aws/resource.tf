resource "opslevel_integration_aws" "example" {
  name = "example"
  role = "arn:aws:iam::XXXXXXXXXXX:role/aws-integration-role"
  external_id = "XXXXXXXXXXXXX"
  ownership_tag_override = true
  ownership_tag_keys = ["owner","team","group"]
}
