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
	s := newStructTagState()

	err := e.marshal(v, s)
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)

	return buf, nil
}

func (e *encodeState) marshal(v interface{}, sts *structTagState) error {
	val := reflect.ValueOf(v)

	if !val.IsValid() {
		return &MarshalError{s: "invalid zero value for marshalling"}
	}

	switch val.Kind() {
	case reflect.Ptr:
		return e.marshal(val.Elem().Interface(), sts)
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).CanInterface() { // only consider exported field symbols
				f := reflect.TypeOf(v).Field(i) // to get the struct tag, need to go via reflect.TypeOf(v).Field(i) and not via reflect.ValueOf(v).Field(i)
				x := f.Tag.Get("xdr")
				s := strings.Split(x, "=")

				switch s[0] {
				case "switch":
					u, ok := val.Field(i).Interface().(uint32)
					if !ok {
						return &MarshalError{s: fmt.Sprintf("invalid type for struct field '%s': require uint32 for `xdr:\"switch\"`", f.Name)}
					}
					sts.switchStatement(u)
				case "case":
					if !sts.isSwitch {
						return &MarshalError{s: fmt.Sprintf("invalid `xdr:\"case=%s\" for struct field '%s': no corresponding `xdr:\"switch\"` statement found", s[1], f.Name)}
					}
					cs := strings.Split(s[1], ",")
					for _, c := range cs {
						u, err := strconv.ParseUint(c, 10, 32)
						if err != nil {
							return &MarshalError{s: fmt.Sprintf("invalid value '%s' in `xdr:\"case=%s\"` for struct field '%s': require uint32 value", c, c, f.Name)}
						}
						sts.caseStatement(uint32(u))
						if sts.matched {
							break
						}
					}
				case "default":
					if !sts.isSwitch {
						return &MarshalError{s: fmt.Sprintf("invalid `xdr:\"default\"` for struct field '%s': no corresponding `xdr:\"switch\"` statement found", f.Name)}
					}
					sts.defaultStatement()
				}

				if s[0] == "switch" || sts.caseMatch() {
					err := e.marshal(val.Field(i).Interface(), newStructTagState())
					if err != nil {
						return err
					}
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
		if ok {
			l := uint32(val.Len())
			err := binary.Write(e, binary.BigEndian, &l)
			if err != nil {
				return err
			}
			_, err = e.Write(a)
			if err != nil {
				return err
			}
			if l%4 > 0 {
				pad := int(4 - (l % 4))
				for i := 0; i < pad; i++ {
					err = e.WriteByte(0)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}
		u, ok := v.([]uint32)
		if ok {
			l := uint32(val.Len())
			err := binary.Write(e, binary.BigEndian, &l)
			if err != nil {
				return err
			}
			err = binary.Write(e, binary.BigEndian, &u)
			if err != nil {
				return err
			}
			return nil
		}
		return &MarshalError{s: "error for type " + val.Type().String() + ": type assertion to []byte / []uint32 failed"}
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
		if len(s)%4 > 0 {
			pad := 4 - (len(s) % 4)
			for i := 0; i < pad; i++ {
				err = e.WriteByte(0)
				if err != nil {
					return err
				}
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
