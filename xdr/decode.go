// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// UnmarshalError is returned by Unmarshal when an unexpected error
// occurred during the unmarshalling process
type UnmarshalError struct {
	s string
}

func (e *UnmarshalError) Error() string {
	return "xdr: " + e.s
}

type decodeState struct {
	data *bytes.Buffer
	off  int // next read offset in data
}

// Unmarshal deserializes a byte array to an XDR format
func Unmarshal(data []byte, v interface{}) (bytesRead int, err error) {
	d := newDecodeState()
	d.init(data)
	return d.unmarshal(v)
}

func (d *decodeState) unmarshal(v interface{}) (bytesRead int, err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return d.off, &UnmarshalError{s: "invalid value for unmarshalling: must be pointer and not nil"}
	}

	val := rv.Elem()
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).CanInterface() { // only consider exported field symbols
				_, err := d.unmarshal(val.Field(i).Addr().Interface())
				if err != nil {
					return d.off, err
				}
			}
		}
	case reflect.Uint32:
		var v uint32
		err := binary.Read(d.data, binary.BigEndian, &v)
		if err != nil {
			return d.off, err
		}
		fmt.Println(v)
		val.SetUint(uint64(v))
	case reflect.Uint64:
		var v uint64
		err := binary.Read(d.data, binary.BigEndian, &v)
		if err != nil {
			return d.off, err
		}
		fmt.Println(v)
		val.SetUint(v)
	default:
		return d.off, &UnmarshalError{s: "unsupported type: " + val.Type().String()}
	}

	return d.off, nil
}

func (d *decodeState) init(data []byte) {
	d.data = bytes.NewBuffer(data)
	d.off = 0
}

func newDecodeState() *decodeState {
	return new(decodeState)
}
