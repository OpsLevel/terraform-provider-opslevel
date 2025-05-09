package internal

const providerIssueUrl = "https://github.com/OpsLevel/terraform-provider-opslevel/issues"

func MergeMaps[TKey comparable, TValue any](map1, map2 map[TKey]TValue) map[TKey]TValue {
	merged := make(map[TKey]TValue)

	for key, value := range map1 {
		merged[key] = value
	}

	for key, value := range map2 {
		if _, present := merged[key]; !present {
			merged[key] = value
		}
	}

	return merged
}
