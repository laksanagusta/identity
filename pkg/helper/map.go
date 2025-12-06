package helper

import "reflect"

func IndexBy[T any, K comparable](items []T, keyFn func(T) K) map[K]T {
	m := make(map[K]T, len(items))
	for _, it := range items {
		m[keyFn(it)] = it
	}
	return m
}

func CollectIDs(slice interface{}) []string {
	v := reflect.ValueOf(slice)

	// pastikan input slice
	if v.Kind() != reflect.Slice {
		panic("CollectIDs: input harus slice")
	}

	ids := make([]string, 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)

		// kalau elem pointer, dereference
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		// pastikan struct
		if elem.Kind() != reflect.Struct {
			continue
		}

		// cek ada field "ID"
		field := elem.FieldByName("UUID")
		if field.IsValid() && field.Kind() == reflect.String {
			ids = append(ids, field.String())
		}
	}

	return ids
}
