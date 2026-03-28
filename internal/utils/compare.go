package utils

import "github.com/google/go-cmp/cmp"

func IsSubsetMap(m1, m2 map[string]any) bool {
	for key, v1 := range m1 {
		v2, ok := m2[key]
		if !ok {
			return false
		}
		if m1, ok := v1.(map[string]any); ok {
			if m2, ok := v2.(map[string]any); ok {
				if !IsSubsetMap(m1, m2) {
					return false
				}
				continue
			}
			return false
		} else if _, ok := v1.([]any); ok {
			continue // skip slice comparison for simplicity, as it's not needed for our use case
		} else if !cmp.Equal(v1, v2) {
			return false
		}
	}
	return true
}

func IsSubSetInterface(a, b any) bool {
	m1, err := ToMap(a)
	if err != nil {
		return false
	}
	m2, err := ToMap(b)
	if err != nil {
		return false
	}
	return IsSubsetMap(m1, m2)
}
