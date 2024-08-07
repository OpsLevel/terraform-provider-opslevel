Testing the OpsLevel Terraform Provider
===============================

# Integration Tests

Tests in [./tests/remote/](./tests/remote/) create actual OpsLevel resources in
the account corresponding to the set `OPSLEVEL_API_TOKEN`. These integration tests
use Terraform's [test framework](https://developer.hashicorp.com/terraform/language/tests)
which uses Terraform config files and tears down resources created by tests.

## Running Integration Tests

Integration tests can be run remotely by a
[manually triggered Github Action](https://github.com/OpsLevel/terraform-provider-opslevel/actions/workflows/tests-integration.yml)
or locally via `task test-integration`. Specific tests we can be run. For example `task test-integration -- -filter=system.tftest.hcl`
and `task test-integration -- -filter=team.tftest.hcl -filter=user.tftest.hcl`, etc.

**Make sure `OPSLEVEL_API_TOKEN` is set to the appropriate test account**. While Terraform's
test framework does its best to clean up test resources, assume there may be cases where
resources are not cleaned up.

### Resource test files

Each resource_type in [./tests/remote/](./tests/remote/) has a `<resource_type>.tftest.hcl`
[test file](https://developer.hashicorp.com/terraform/language/tests#syntax). These test files
have multiple [run blocks](https://developer.hashicorp.com/terraform/language/tests#run-blocks)
where [assert blocks](https://developer.hashicorp.com/terraform/language/tests#assertions)
verify each OpsLevel resource was correctly applied. Resources are defined in modules and are
sourced into each test file.

### Resource modules

Each resource_type in [./tests/remote/](./tests/remote/) has a `./tests/remote/<resource_type>/`
directory that is sourced into a test file. [Modules](https://developer.hashicorp.com/terraform/language/modules)
are basically directories. Each module has the following:
- `main.tf` - resources under test are defined here in a standard Terraform config file.
- `variables.tf` - input variables are defined here. See `variables` blocks in `tftest.hcl` files.
- `outputs.tf` - module resources returned from implicit `terraform apply` during test.

# Local Tests (terraform validate)

Tests in [./tests/local/](./tests/local/) use a mocked OpsLevel Terraform provider
so no authentication is required. These tests are meant to catch config validation
errors.

## Running Local Tests

Local tests can be run via `task test` and make up part of `task ci`.

### TODO: Focus Local Tests on config validation, remove others

Local tests should cover all the uses cases a customer may encounter when using
`terraform validate` on their Terraform config files with OpsLevel resources.

`terraform validate` catches pre-plan misconfigurations and is baked into
`terraform plan` itself.
