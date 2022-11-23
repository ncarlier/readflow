package helper

// If return value regarding the condition
func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// IfNil return value if obj is nil
func IfNil[T any](obj *T, fallback T) T {
	if obj == nil {
		return fallback
	}
	return *obj
}
