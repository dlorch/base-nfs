package xdr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

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
	case reflect.Uint64:
		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.BigEndian, val.Uint())
		return buffer.Bytes(), nil
	case reflect.Uint32:
		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.BigEndian, uint32(val.Uint()))
		return buffer.Bytes(), nil
	default:
		fmt.Println("Unrecognized type:", val.Kind())
	}
	return buf, nil
}