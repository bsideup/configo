package flatmap

import (
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	cases := []struct {
		Input  map[string]interface{}
		Output map[string]interface{}
	}{
		{
			Input: map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			Output: map[string]interface{}{
				"FOO": "bar",
				"BAR": "baz",
			},
		},

		{
			Input: map[string]interface{}{
				"foo": []string{
					"one",
					"two",
				},
			},
			Output: map[string]interface{}{
				"FOO_0": "one",
				"FOO_1": "two",
			},
		},

		{
			Input: map[string]interface{}{
				"foo": []map[interface{}]interface{}{
					map[interface{}]interface{}{
						"name":    "bar",
						"port":    3000,
						"enabled": true,
					},
				},
			},
			Output: map[string]interface{}{
				"FOO_0_NAME":    "bar",
				"FOO_0_PORT":    3000,
				"FOO_0_ENABLED": true,
			},
		},

		{
			Input: map[string]interface{}{
				"foo": []map[interface{}]interface{}{
					map[interface{}]interface{}{
						"name": "bar",
						"ports": []string{
							"1",
							"2",
						},
					},
				},
			},
			Output: map[string]interface{}{
				"FOO_0_NAME":    "bar",
				"FOO_0_PORTS_0": "1",
				"FOO_0_PORTS_1": "2",
			},
		},
	}

	for _, tc := range cases {
		actual := Flatten(tc.Input, true)
		if !reflect.DeepEqual(actual, map[string]interface{}(tc.Output)) {
			t.Fatalf(
				"Input:\n\n%#v\n\nOutput:\n\n%#v\n\nExpected:\n\n%#v\n",
				tc.Input,
				actual,
				tc.Output)
		}
	}
}
