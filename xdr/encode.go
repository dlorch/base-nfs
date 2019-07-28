package xdr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// UnsupportedTypeError is returned by Marshal when attempting to
// encode an unsupported value type
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "xdr: unsupported type: " + e.Type.String()
}

// MarshalError is returned by Marshal when an unexpected error
// occurred during the marshalling process
type MarshalError struct {
	T reflect.Type
	S string
}

func (e *MarshalError) Error() string {
	return "xdr: error for type " + e.T.String() + ": " + e.S
}

// Marshal serializes a value in XDR format to a byte sequence representation
func Marshal(v interface{}) ([]byte, error) {
	var buf []byte

	val := reflect.ValueOf(v)
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
	case reflect.Slice:
		a, ok := v.([]byte)
		if !ok {
			return buf, &MarshalError{T: val.Type(), S: "type assertion to []byte failed"}
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
			return buf, &UnsupportedTypeError{Type: val.Type()}
		}
	}
	return buf, nil
}
