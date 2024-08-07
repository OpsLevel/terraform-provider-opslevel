mock_data "opslevel_service" {
  defaults = {
    aliases                       = ["alias-one", "alias-two"]
    api_document_path             = "mock-api-document-path"
    description                   = "mock-service-description"
    framework                     = "mock-framework"
    language                      = "mock-language"
    lifecycle_alias               = "alpha"
    name                          = "mock-service-name"
    owner                         = "mock-team"
    owner_id                      = "Z2lkOi8vmockowner123"
    preferred_api_document_source = "PULL"
    product                       = "mock-product"
    properties = {
      "definition" = {
        aliases = [
          "mock-one",
          "mock-two",
          "mock-three",
        ]
        id = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8yODk"
      }
      "value" = "mock-property-definition",
    }
    repositories = ["repo-one", "repo-two"]
    tags         = ["key1:value2", "key2:value2"]
    tier_alias   = "mock-tier-alias"
  }
}
