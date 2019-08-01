package xdr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

}

// MarshalError is returned by Marshal when an unexpected error
// occurred during the marshalling process
type MarshalError struct {
	s string
}

func (e *MarshalError) Error() string {
	return "xdr: " + e.s
}

// Marshal serializes a value in XDR format to a byte sequence representation
func Marshal(v interface{}) ([]byte, error) {
	var buf []byte
	val := reflect.ValueOf(v)

	if !val.IsValid() {
		return buf, &MarshalError{s: "invalid zero value for 'v'"}
	}

	switch val.Kind() {
	case reflect.Ptr:
		return Marshal(val.Elem().Interface())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).CanInterface() { // only consider exported field symbols
				b, err := Marshal(val.Field(i).Interface())
				if err != nil {
					return buf, err
				}
				buf = append(buf, b...)
			}
		}
	case reflect.Array:
		b := make([]byte, val.Len())
		for i := 0; i < val.Len(); i++ {
			b[i] = val.Index(i).Interface().(byte)
		}
		return b, nil
	case reflect.Slice:
		a, ok := v.([]byte)
		if !ok {
			return buf, &MarshalError{s: "error for type " + val.Type().String() + ": type assertion to []byte failed"}
		}
		l := uint32(len(a))
		b := new(bytes.Buffer)
		err := binary.Write(b, binary.BigEndian, &l)
		if err != nil {
			return buf, err
		}
		b.Write(a)
		return b.Bytes(), nil
	case reflect.String:
		b := new(bytes.Buffer)
		s := val.String()
		l := uint32(len(s))
		pad := 4 - (len(s) % 4)
		err := binary.Write(b, binary.BigEndian, &l)
		if err != nil {
			return buf, err
		}
		_, err = b.WriteString(s)
		if err != nil {
			return buf, err
		}
		for i := 0; i < pad; i++ {
			err = b.WriteByte(0)
			if err != nil {
				return buf, err
			}
		}
		return b.Bytes(), nil
	case reflect.Uint64:
		b := new(bytes.Buffer)
		err := binary.Write(b, binary.BigEndian, val.Uint())
		return b.Bytes(), err
	case reflect.Uint32:
		b := new(bytes.Buffer)
		err := binary.Write(b, binary.BigEndian, uint32(val.Uint()))
		return b.Bytes(), err
	default:
		if val.IsValid() {
			return buf, &MarshalError{s: "unsupported type: " + val.Type().String()}
		}
	}
	return buf, nil
}
