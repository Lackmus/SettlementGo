package helpers

func CopyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func DeepCopyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		switch val := any(v).(type) {
		case map[K]V:
			dst[k] = any(DeepCopyMap(val)).(V) // Recursive map copy
		case []V:
			newSlice := make([]V, len(val))
			copy(newSlice, val) // Copy slice contents
			dst[k] = any(newSlice).(V)
		default:
			dst[k] = v
		}
	}
	return dst
}

func CopyStringSliceMap(source map[string][]string) map[string][]string {
	if len(source) == 0 {
		return map[string][]string{}
	}

	result := make(map[string][]string, len(source))
	for key, values := range source {
		result[key] = append([]string(nil), values...)
	}

	return result
}
