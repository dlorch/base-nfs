// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdr

import "reflect"

// UnmarshalError is returned by Unmarshal when an unexpected error
// occurred during the unmarshalling process
type UnmarshalError struct {
	s string
}

func (e *UnmarshalError) Error() string {
	return "xdr: " + e.s
}

type decodeState struct {
	data []byte
	off  int // next read offset in data
}

// Unmarshal deserializes a byte array to an XDR format
func Unmarshal(data []byte, v interface{}) (bytesRead int, err error) {
	d := newDecodeState()
	d.init(data)
	return d.unmarshal(v)
}

func (d *decodeState) unmarshal(v interface{}) (bytesRead int, err error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return d.off, &UnmarshalError{s: "invalid value for unmarshalling: must be pointer and not nil"}
	}

	return 0, nil
}

func (d *decodeState) init(data []byte) {
	d.data = data
	d.off = 0
}

func newDecodeState() *decodeState {
	return new(decodeState)
}
