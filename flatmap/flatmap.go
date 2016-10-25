package flatmap

import (
	"fmt"
	"reflect"
	"strings"
)

// Flatten takes a structure and turns into a flat map[string]interface{}.
//
// Within the "thing" parameter, only primitive values are allowed. Structs are
// not supported. Therefore, it can only be slices, maps, primitives, and
// any combination of those together.
func Flatten(thing map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, raw := range thing {
		flatten(result, k, reflect.ValueOf(raw))
	}

	return result
}

func flatten(result map[string]interface{}, prefix string, v reflect.Value) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Invalid:
		return
	case reflect.Map:
		flattenMap(result, prefix, v)
	case reflect.Slice:
		flattenSlice(result, prefix, v)
	default:
		result[strings.ToUpper(prefix)] = v.Interface()
	}
}

func flattenMap(result map[string]interface{}, prefix string, v reflect.Value) {
	for _, k := range v.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}

		if k.Kind() != reflect.String {
			panic(fmt.Sprintf("%s: map key is not string: %s", prefix, k))
		}

		flatten(result, fmt.Sprintf("%s_%s", prefix, k.String()), v.MapIndex(k))
	}
}

func flattenSlice(result map[string]interface{}, prefix string, v reflect.Value) {
	prefix = prefix + "_"

	for i := 0; i < v.Len(); i++ {
		flatten(result, fmt.Sprintf("%s%d", prefix, i), v.Index(i))
	}
}
