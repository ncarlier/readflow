package tooling

import "strconv"

// ConvGQLStringToUint converts ID GraphQL interface object to unsigned integer
func ConvGQLStringToUint(id interface{}) (uint, bool) {
	str, ok := id.(string)
	if !ok {
		return 0, ok
	}
	ui64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(ui64), true
}

// ConvGQLIntToUint converts GraphQL int to unsigned integer
func ConvGQLIntToUint(val interface{}) (uint, bool) {
	if iVal, ok := val.(int); ok {
		return uint(iVal), true
	}
	return 0, false
}
