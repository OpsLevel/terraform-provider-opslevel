variables {
  resource_name = "opslevel_integration_endpoint"

  # required fields
  name = "TF Test endpoint Integration"
  type = "deploy"

  # optional fields
}

run "resource_integration_endpoint_create_api_doc_type" {

  variables {
    type = "apiDoc"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_update_api_doc_type" {

  variables {
    name = "apiDoc Updated"
    type = "apiDoc"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_aqua_security_type" {

  variables {
    type = "aquaSecurity"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_argocd_type" {

  variables {
    type = "argocd"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_aws_ecr_type" {

  variables {
    type = "awsEcr"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_bugsnag_type" {

  variables {
    type = "bugsnag"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_circleci_type" {

  variables {
    type = "circleci"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_codacy_type" {

  variables {
    type = "codacy"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_coveralls_type" {

  variables {
    type = "coveralls"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_custom_event_type" {

  variables {
    type = "customEvent"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_datadog_check_type" {

  variables {
    type = "datadogCheck"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_deploy_type" {

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = alltrue([
      can(opslevel_integration_endpoint.this.id),
      can(opslevel_integration_endpoint.this.name),
      can(opslevel_integration_endpoint.this.type),
      can(opslevel_integration_endpoint.this.webhook_url),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_integration_endpoint.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_integration_endpoint.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_integration_endpoint.this.name,
    )
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_dynatrace_type" {

  variables {
    type = "dynatrace"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_flux_type" {

  variables {
    type = "flux"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_github_actions_type" {

  variables {
    type = "githubActions"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_gitlab_ci_type" {

  variables {
    type = "gitlabCi"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_grafana_type" {

  variables {
    type = "grafana"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_grype_type" {

  variables {
    type = "grype"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_jenkins_type" {

  variables {
    type = "jenkins"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_jfrog_xray_type" {

  variables {
    type = "jfrogXray"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_lacework_type" {

  variables {
    type = "lacework"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_new_relic_check_type" {

  variables {
    type = "newRelicCheck"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_octopus_type" {

  variables {
    type = "octopus"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_prisma_cloud_type" {

  variables {
    type = "prismaCloud"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_prometheus_type" {

  variables {
    type = "prometheus"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_rollbar_type" {

  variables {
    type = "rollbar"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_sentry_type" {

  variables {
    type = "sentry"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_snyk_type" {

  variables {
    type = "snyk"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_sonarqube_type" {

  variables {
    type = "sonarqube"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_stackhawk_type" {

  variables {
    type = "stackhawk"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_sumo_logic_type" {

  variables {
    type = "sumoLogic"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}

run "resource_integration_endpoint_create_veracode_type" {

  variables {
    type = "veracode"
  }

  module {
    source = "./opslevel_modules/modules/integration/endpoint"
  }

  assert {
    condition = opslevel_integration_endpoint.this.type == var.type
    error_message = format(
      "expected '%v' but got '%v'",
      var.type,
      opslevel_integration_endpoint.this.type,
    )
  }
}
