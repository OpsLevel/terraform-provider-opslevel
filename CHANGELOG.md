<a name="unreleased"></a>
## [Unreleased]


<a name="v0.3.0-beta"></a>
## [v0.3.0-beta] - 2021-09-18
### Bugfix
- importing resources service_repository, service_tag and service_tool now works properly
- opslevel_filter resources where not properly importing connective value
- fix a logic error when looking up resources by alias or id

### Docs
- add badges to github readme

### Feature
- upgrade opslevel-go libraryy to 0.3.3
- add ability to override the OpsLevel API Url

### Refactor
- change name of provider field from `apitoken` to `api_token` to standardize with CLI and Kubectl
- fix name of `opslevel_check_service_owner` to be `opslevel_check_service_ownership`


<a name="v0.2.0-beta"></a>
## [v0.2.0-beta] - 2021-08-29
### Feature
- initial terraform provider implementation (teams, services, filters, rubrics, checks)


<a name="v0.1.0"></a>
## v0.1.0 - 2021-08-22

[Unreleased]: https://github.com/OpsLevel/terraform-provider-opslevel/compare/v0.3.0-beta...HEAD
[v0.3.0-beta]: https://github.com/OpsLevel/terraform-provider-opslevel/compare/v0.2.0-beta...v0.3.0-beta
[v0.2.0-beta]: https://github.com/OpsLevel/terraform-provider-opslevel/compare/v0.1.0...v0.2.0-beta
