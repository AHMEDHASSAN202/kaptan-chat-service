package helper

func GenerateLocationId(id string, typ string) (string, string) {
	locations := map[string]string{
		"App\\Models\\Hotel":   "hotel",
		"App\\Models\\Station": "station",
	}

	if id != "" && typ != "" {
		if prefix, ok := locations[typ]; ok {
			return prefix + "-" + id, prefix
		}
	}
	return id, typ
}

func GenerateClientId(id string, typ string) (string, string) {
	clients := map[string]string{
		"App\\Models\\Hotel":        "hotel",
		"App\\Models\\TravelAgency": "company",
	}

	if id != "" && typ != "" {
		if prefix, ok := clients[typ]; ok {
			return prefix + "-" + id, prefix
		}
	}
	return id, typ
}
