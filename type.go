package marker

import (
	"fmt"
	"reflect"
)

type Type int

const (
	InvalidType Type = iota
	RawType
	AnyType
	BoolType
	IntegerType
	StringType
	SliceType
	MapType
)

var (
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	rawType       = reflect.TypeOf((*[]byte)(nil)).Elem()
)

func GetType(typ reflect.Type) (Type, error) {
	if typ == rawType {
		return RawType, nil
	}

	if typ == interfaceType {
		return AnyType, nil
	}

	if typ.Kind() == reflect.Ptr {
		rawType = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		return StringType, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return IntegerType, nil
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		return IntegerType, nil
	case reflect.Bool:
		return BoolType, nil
	case reflect.Slice:
		return SliceType, nil
	case reflect.Map:
		if typ.Key().Kind() != reflect.String {
			return InvalidType, fmt.Errorf("bad map key type: map key must be string")
		}

		return MapType, nil
	default:
		return InvalidType, fmt.Errorf("type has unsupported kind %s", rawType.Kind())
	}
}
