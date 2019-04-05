package ruleengine

func toBool(i1 interface{}) bool {
	if i1 == nil {
		return false
	}
	switch i2 := i1.(type) {
	default:
		return false
	case bool:
		return i2
	case string:
		return i2 == "true"
	case int:
		return i2 != 0
	case *bool:
		if i2 == nil {
			return false
		}
		return *i2
	case *string:
		if i2 == nil {
			return false
		}
		return *i2 == "true"
	case *int:
		if i2 == nil {
			return false
		}
		return *i2 != 0
	}
}
