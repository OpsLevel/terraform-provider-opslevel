resource "opslevel_integration_google_cloud" "dev" {
  client_email            = "hello-world-tf@powerful-surf-427415-v1.iam.gserviceaccount.com"
  name                    = "GCP Dev Integration"
  ownership_tag_keys      = ["opslevel_team", "team", "owner"]
  ownership_tag_overrides = true
  private_key             = <<EOH
-----BEGIN PRIVATE KEY-----
<private key body here>
-----END PRIVATE KEY-----
EOH
}
