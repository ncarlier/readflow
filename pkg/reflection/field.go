package reflection

import (
	"math"
	"reflect"
)

type Field struct {
	value interface{}
}

func (p *Field) Exists() bool {
	return p.value != nil
}

func (p *Field) Bool() (bool, bool) {
	if p.value != nil {
		v, ok := p.value.(bool)
		return v, ok
	}
	return false, false
}

func (p *Field) String() (string, bool) {
	if p.value != nil {
		v, ok := p.value.(string)
		return v, ok
	}
	return "", false
}

func (p *Field) Int() (int64, bool) {
	switch vt := p.value.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(vt).Int(), true
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(vt).Uint()
		if val <= math.MaxInt64 {
			return int64(val), true
		}
	case float32, float64:
		val := reflect.ValueOf(vt).Float()
		if val >= math.MinInt64 && val <= math.MaxInt64 {
			return int64(val), true
		}
	}
	return 0, false
}

func (p *Field) Map() (map[string]interface{}, bool) {
	val := reflect.ValueOf(p.value)
	if val.IsValid() && val.CanConvert(mapType) {
		return val.Convert(mapType).Interface().(map[string]interface{}), true
	}
	return nil, false
}

func (p *Field) Slice() ([]interface{}, bool) {
	val := reflect.ValueOf(p.value)
	if val.IsValid() && val.CanConvert(sliceType) {
		return val.Convert(sliceType).Interface().([]interface{}), true
	}
	return nil, false
}

func (p *Field) IsValidObject() bool {
	val := reflect.ValueOf(p.value)
	return val.IsValid() && (val.CanConvert(mapType) || val.CanConvert(sliceType))
}
