locals {
  exists_types            = ["does_not_exist", "exists"]
  equals_and_exists_types = concat(["does_not_equal", "equals"], local.exists_types)
  string_predicates = setproduct(
    ["aliases", "framework", "language", "name", "product"],
    concat([
      "contains",
      "does_not_contain",
      "does_not_match_regex",
      "ends_with",
      "matches_regex",
      "starts_with",
    ], local.equals_and_exists_types),
  )
  numeric_predicates = setproduct(
    ["lifecycle_index", "tier_index"],
    concat([
      "less_than_or_equal_to",
      "greater_than_or_equal_to",
    ], local.equals_and_exists_types),
  )
  tags_predicates = setproduct(
    ["tags"],
    concat([
      "contains",
      "does_not_contain",
      "does_not_match_regex",
      "ends_with",
      "matches_regex",
      "satisfies_version_constraint",
      "starts_with",
    ], local.equals_and_exists_types),
  )
  domain_id_predicates      = setproduct(["domain_id"], local.equals_and_exists_types)
  filter_id_predicates      = setproduct(["filter_id"], ["does_not_match", "matches"])
  owner_id_predicates       = setproduct(["owner_id"], local.equals_and_exists_types)
  properties_predicates     = setproduct(["properties"], concat(local.exists_types, ["satisfies_jq_expression"]))
  system_id_predicates      = setproduct(["system_id"], local.equals_and_exists_types)
  repository_ids_predicates = setproduct(["repository_ids"], local.exists_types)
  all_predicates = merge(
    tomap({ for a, pair in local.numeric_predicates : "test_numeric_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = null, value = contains(local.exists_types, pair[1]) ? null : 1
      }
    }),
    tomap({ for a, pair in local.string_predicates : "test_string_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = null, value = contains(local.exists_types, pair[1]) ? null : "value_${a}"
      }
    }),
    tomap({ for a, pair in local.domain_id_predicates : "test_domain_id_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = null, value = contains(local.exists_types, pair[1]) ? null : data.opslevel_domains.all.domains[0].id
      }
    }),
    tomap({ for a, pair in local.filter_id_predicates : "test_filter_id_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = null, value = data.opslevel_filters.test.filters[0].id
      }
    }),
    tomap({ for a, pair in local.owner_id_predicates : "test_owner_id_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = null, value = contains(local.exists_types, pair[1]) ? null : data.opslevel_teams.all.teams[0].id
      }
    }),
    tomap({ for a, pair in local.properties_predicates : "test_tags_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = "test_tag_key_${a}", value = contains(local.exists_types, pair[1]) ? null : "test jq expression"
      }
    }),
    tomap({ for a, pair in local.repository_ids_predicates : "test_repository_ids_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = null, value = null
      }
    }),
    tomap({ for a, pair in local.system_id_predicates : "test_ids_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = null, value = contains(local.exists_types, pair[1]) ? null : data.opslevel_systems.all.systems[0].id
      }
    }),
    tomap({ for a, pair in local.tags_predicates : "test_tags_predicate_${a}" => {
      key = pair[0], type = pair[1], key_data = "test_tag_key_${a}", value = contains(local.exists_types, pair[1]) ? null : "test_tag_${a}"
      }
    }),
    { "test_predicate_owner_ids_0" = { key = "owner_ids", type = "equals", key_data = null, value = data.opslevel_teams.all.teams[0].id } },
  )
}
