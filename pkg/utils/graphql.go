package utils

import (
	"math/big"
)

// ParseGraphQLArgument parse GraphQL argument
func ParseGraphQLArgument[T string | bool | int](args map[string]interface{}, name string) *T {
	if arg, exist := args[name]; exist {
		if val, ok := arg.(T); ok {
			return &val
		}
	}
	return nil
}

// ParseGraphQLID parse GraphQL argument as uint
func ParseGraphQLID(args map[string]interface{}, name string) *uint {
	if arg, exist := args[name]; exist {
		return ConvGraphQLID(arg)
	}
	return nil
}

// ParseGraphQLID parse GraphQL argument as uint
func ConvGraphQLID(id interface{}) *uint {
	switch val := id.(type) {
	case int:
		v := uint(val)
		return &v
	case string:
		if f, _, err := big.ParseFloat(val, 10, 0, big.ToNearestEven); err == nil {
			v64, _ := f.Uint64()
			v := uint(v64)
			return &v
		}
	}
	return nil
}
