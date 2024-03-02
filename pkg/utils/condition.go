package utils

// If return value regarding the condition
func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// PtrValueOr return value of the pointer or a fallback if nil
func PtrValueOr[T any](obj *T, fallback T) T {
	if obj == nil {
		return fallback
	}
	return *obj
}
