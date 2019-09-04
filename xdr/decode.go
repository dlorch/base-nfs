// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	s := newStructTagState()
	d.init(data)
	return d.unmarshal(v, s)
}

func (d *decodeState) unmarshal(v interface{}, sts *structTagState) (bytesRead int, err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return d.off, &UnmarshalError{s: "invalid value for unmarshalling: must be pointer and not nil"}
	}

	val := rv.Elem()
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).CanInterface() { // only consider exported field symbols
				f := reflect.TypeOf(v).Elem().Field(i) // to get the struct tag, need to go via reflect.TypeOf(v).Field(i) and not via reflect.ValueOf(v).Field(i)
				x := f.Tag.Get("xdr")
				s := strings.Split(x, "=")

				switch s[0] {
				case "switch":
					if val.Field(i).Kind() != reflect.Uint32 {
						return d.off, &MarshalError{s: fmt.Sprintf("invalid type for struct field '%s': require uint32 for `xdr:\"switch\"`", f.Name)}
					}
					_, err := d.unmarshal(val.Field(i).Addr().Interface(), sts)
					if err != nil {
						return d.off, err
					}
					u, _ := val.Field(i).Interface().(uint32)
					sts.switchStatement(u)
					continue
				case "case":
					if !sts.isSwitch {
						return d.off, &UnmarshalError{s: fmt.Sprintf("invalid `xdr:\"case=%s\" for struct field '%s': no corresponding `xdr:\"switch\"` statement found", s[1], f.Name)}
					}
					u, err := strconv.ParseUint(s[1], 10, 32)
					if err != nil {
						return d.off, &UnmarshalError{s: fmt.Sprintf("invalid value '%s' in `xdr:\"case=%s\"` for struct field '%s': require uint32 value", s[1], s[1], f.Name)}
					}
					sts.caseStatement(uint32(u))
				case "default":
					if !sts.isSwitch {
						return d.off, &UnmarshalError{s: fmt.Sprintf("invalid `xdr:\"default\"` for struct field '%s': no corresponding `xdr:\"switch\"` statement found", f.Name)}
					}
					sts.defaultStatement()
				}

				if sts.caseMatch() {
					_, err := d.unmarshal(val.Field(i).Addr().Interface(), sts)
					if err != nil {
						return d.off, err
					}
				}
			}
		}
	case reflect.Array:
		var b byte
		for i := 0; i < val.Len(); i++ {
			err := binary.Read(d.data, binary.BigEndian, &b)
			if err != nil {
				return d.off, err
			}
			val.Index(i).SetUint(uint64(b))
		}
	case reflect.Slice:
		var l uint32
		err := binary.Read(d.data, binary.BigEndian, &l)
		if err != nil {
			return d.off, err
		}

		_, ok := v.(*[]byte)
		if ok {
			b := make([]byte, l)
			n, err := d.data.Read(b)
			if err != nil {
				return d.off, err
			}
			if n != int(l) {
				return d.off, &UnmarshalError{s: fmt.Sprintf("slice variable supposed to be length %d, but could only ready %d bytes", l, n)}
			}
			if l%4 > 0 {
				pad := int(4 - (l % 4))
				for i := 0; i < pad; i++ {
					_, err := d.data.ReadByte()
					if err != nil {
						return d.off, err
					}
				}
			}
			val.SetBytes(b)
			return d.off, nil
		}
		_, ok = v.(*[]uint32)
		if ok {
			u := make([]uint32, l)
			err := binary.Read(d.data, binary.BigEndian, &u)
			if err != nil {
				return d.off, err
			}
			val.Set(reflect.ValueOf(u))
			return d.off, nil
		}
		return d.off, &UnmarshalError{s: "error for type " + val.Type().String() + ": type assertion to []byte / []uint32 failed"}
	case reflect.String:
		var len uint32
		err := binary.Read(d.data, binary.BigEndian, &len)
		if err != nil {
			return d.off, err
		}
		b := make([]byte, len)
		n, err := d.data.Read(b)
		if err != nil {
			return d.off, err
		}
		if n != int(len) {
			return d.off, &UnmarshalError{s: fmt.Sprintf("string variable supposed to be length %d, but could only ready %d bytes", len, n)}
		}
		if len%4 > 0 {
			pad := int(4 - (len % 4))
			for i := 0; i < pad; i++ {
				_, err := d.data.ReadByte()
				if err != nil {
					return d.off, err
				}
			}
		}
		val.SetString(string(b))
	case reflect.Uint32:
		var v uint32
		err := binary.Read(d.data, binary.BigEndian, &v)
		if err != nil {
			return d.off, err
		}
		val.SetUint(uint64(v))
	case reflect.Uint64:
		var v uint64
		err := binary.Read(d.data, binary.BigEndian, &v)
		if err != nil {
			return d.off, err
		}
		val.SetUint(v)
	case reflect.Ptr:
		val.Set(reflect.New(val.Type().Elem()))
		_, err := d.unmarshal(val.Elem().Addr().Interface(), newStructTagState())
		if err != nil {
			return d.off, err
		}
	default:
		return d.off, &UnmarshalError{s: "unsupported type: " + val.Type().String() + " of kind " + val.Kind().String()}
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
