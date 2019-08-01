package xdr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// MarshalError is returned by Marshal when an unexpected error
// occurred during the marshalling process
type MarshalError struct {
	s string
}

func (e *MarshalError) Error() string {
	return "xdr: " + e.s
}

type encodeState struct {
	bytes.Buffer // accumulated output
}

// Marshal serializes a value in XDR format to a byte sequence representation
func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()

	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)

	return buf, nil
}

func (e *encodeState) marshal(v interface{}) error {
	val := reflect.ValueOf(v)

	if !val.IsValid() {
		return &MarshalError{s: "invalid zero value for 'v'"}
	}

	switch val.Kind() {
	case reflect.Ptr:
		return e.marshal(val.Elem().Interface())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).CanInterface() { // only consider exported field symbols
				err := e.marshal(val.Field(i).Interface())
				if err != nil {
					return err
				}
			}
		}
	case reflect.Array:
		for i := 0; i < val.Len(); i++ {
			b, ok := val.Index(i).Interface().(byte)
			if !ok {
				return &MarshalError{s: "error for type " + val.Type().String() + ": type assertion to byte failed"}
			}
			err := e.WriteByte(b)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice:
		a, ok := v.([]byte)
		if !ok {
			return &MarshalError{s: "error for type " + val.Type().String() + ": type assertion to []byte failed"}
		}
		l := uint32(len(a))
		err := binary.Write(e, binary.BigEndian, &l)
		if err != nil {
			return err
		}
		_, err = e.Write(a)
		return err
	case reflect.String:
		s := val.String()
		l := uint32(len(s))
		err := binary.Write(e, binary.BigEndian, &l)
		if err != nil {
			return err
		}
		_, err = e.WriteString(s)
		if err != nil {
			return err
		}
		pad := 4 - (len(s) % 4)
		for i := 0; i < pad; i++ {
			err = e.WriteByte(0)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Uint64:
		err := binary.Write(e, binary.BigEndian, val.Uint())
		return err
	case reflect.Uint32:
		err := binary.Write(e, binary.BigEndian, uint32(val.Uint()))
		return err
	default:
		return &MarshalError{s: "unsupported type: " + val.Type().String()}
	}
	return nil
}

func newEncodeState() *encodeState {
	return new(encodeState)
}
