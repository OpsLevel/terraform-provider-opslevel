enum_alert_source_status_type = [
  "alert",
  "warn",
  "ok",
  "no_data",
  "fetching_data",
]
enum_alert_source_type = [
  "pagerduty",
  "datadog",
  "opsgenie",
  "new_relic",
]
enum_predicate_type = [
  "contains",
  "does_not_contain",
  "does_not_equal",
  "does_not_exist",
  "ends_with",
  "equals",
  "exists",
  "greater_than_or_equal_to",
  "less_than_or_equal_to",
  "starts_with",
  "satisfies_version_constraint",
  "matches_regex",
  "does_not_match_regex",
  "belongs_to",
  "matches",
  "does_not_match",
  "satisfies_jq_expression",
]
enum_tool_category = [
  "admin",
  "api_documentation",
  "architecture_diagram",
  "backlog",
  "code",
  "continuous_integration",
  "deployment",
  "design_documentation",
  "errors",
  "feature_flag",
  "health_checks",
  "incidents",
  "issue_tracking",
  "logs",
  "metrics",
  "observability",
  "orchestrator",
  "other",
  "resiliency",
  "runbooks",
  "security_scans",
  "status_page",
  "wiki",
]
error_empty_datasource             = "zero 'TYPE' found in 'TYPE' datasource"
error_expected_empty_string        = "expected field to be an empty string"
error_expected_null_field          = "expected field to be null"
error_unexpected_datasource_fields = "cannot reference all expected 'TYPE' datasource fields"
error_unexpected_resource_fields   = "cannot reference all expected 'TYPE' resource fields"
error_wrong_alias                  = "wrong alias for 'TYPE'"
error_wrong_description            = "wrong description for 'TYPE'"
error_wrong_index                  = "wrong index for 'TYPE'"
error_wrong_id                     = "wrong id for 'TYPE'"
error_wrong_name                   = "wrong name for 'TYPE'"
error_wrong_owner                  = "wrong owner for 'TYPE'"
id_prefix                          = "Z2lkOi8v"
predicate_types_equals             = ["does_not_equal", "equals"]
predicate_types_exists             = ["does_not_exist", "exists"]
