resource "opslevel_service_tag" "service_tag_1" {
  key     = "hello"
  value   = "world"
  service = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MzMyOQ"
}

resource "opslevel_service_tag" "service_tag_2" {
  key           = "hello_with_alias"
  value         = "world_with_alias"
  service_alias = "cart"
}
