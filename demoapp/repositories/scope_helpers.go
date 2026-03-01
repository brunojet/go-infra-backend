package repositories

func RequireScopeInt64(scopes map[string]any, key string, invalidErr error) (int64, error) {
	rawValue, ok := scopes[key]
	if !ok {
		return 0, invalidErr
	}

	value, ok := rawValue.(int64)
	if !ok || value <= 0 {
		return 0, invalidErr
	}

	return value, nil
}
