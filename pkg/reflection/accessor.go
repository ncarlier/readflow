package reflection

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	propNameMatcher = regexp.MustCompile(`([^[]+)?(?:\[(\d+)])?`)
	mapType         = reflect.TypeOf(map[string]interface{}{})
	sliceType       = reflect.TypeOf([]interface{}{})
)

func decodeFieldName(name string) (string, int) {
	parts := propNameMatcher.FindStringSubmatch(strings.TrimSpace(name))
	index := -1
	if len(parts) != 3 {
		return name, index
	}
	if parts[2] != "" {
		index, _ = strconv.Atoi(parts[2])
	}
	return parts[1], index
}

func getField(obj interface{}, fieldName string) (interface{}, bool) {
	val := reflect.ValueOf(obj)
	if !val.IsValid() || !val.CanConvert(mapType) {
		return obj, false
	}
	name, index := decodeFieldName(fieldName)
	field := val.Convert(mapType).Interface().(map[string]interface{})[name]

	if index < 0 {
		return field, true
	}

	// manage slice part
	val = reflect.ValueOf(field)
	if val.IsValid() && val.CanConvert(sliceType) {
		slice := val.Convert(sliceType).Interface().([]interface{})
		if index >= 0 && index < len(slice) {
			return slice[index], true
		}
	}
	return field, false
}

// GetField try to get struct field from a given path
func GetField(obj interface{}, path string) *Field {
	fieldNames := strings.Split(path, ".")

	current := obj
	for _, fieldName := range fieldNames {
		var ok bool
		current, ok = getField(current, fieldName)
		if !ok {
			return &Field{}
		}
	}
	return &Field{
		value: current,
	}
}
