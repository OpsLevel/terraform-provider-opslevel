mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_infra_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = !contains([null, ""], opslevel_infrastructure.small_infra.id)
    error_message = "opslevel_infrastructure.small_infra id should not be empty"
  }

  assert {
    condition     = !contains([null, ""], opslevel_infrastructure.small_infra.last_updated)
    error_message = "opslevel_infrastructure.small_infra last_updated should not be empty"
  }

  assert {
    condition = opslevel_infrastructure.small_infra.data == jsonencode({
      name = "small-query"
    })
    error_message = "wrong data in opslevel_infrastructure.small_infra"
  }

  assert {
    condition     = opslevel_infrastructure.small_infra.owner == var.test_id
    error_message = "wrong owner for opslevel_infrastructure.small_infra"
  }

  assert {
    condition     = opslevel_infrastructure.small_infra.schema == "Small Database"
    error_message = "wrong schema for opslevel_infrastructure.small_infra"
  }

}

run "resource_infra_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_infrastructure.big_infra.aliases == tolist(["big-infra"])
    error_message = "wrong aliases in opslevel_infrastructure.big_infra"
  }

  assert {
    condition = opslevel_infrastructure.big_infra.data == jsonencode({
      name                = "big-query"
      external_id         = 1234
      replica             = true
      publicly_accessible = false
      storage_size = {
        unit  = "GB"
        value = 700
      }
    })
    error_message = "wrong data in opslevel_infrastructure.big_infra"
  }

  assert {
    condition     = startswith(opslevel_infrastructure.big_infra.owner, "Z2lkOi8v")
    error_message = "wrong owner in opslevel_infrastructure.big_infra"
  }

  assert {
    condition = opslevel_infrastructure.big_infra.provider_data == {
      account = "dev"
      name    = "google cloud"
      type    = "BigQuery"
      url     = "https://console.cloud.google.com/"
    }
    error_message = "wrong provider_data in opslevel_infrastructure.big_infra"
  }

  assert {
    condition     = opslevel_infrastructure.big_infra.schema == "Big Database"
    error_message = "wrong schema for opslevel_infrastructure.big_infra"
  }

}
