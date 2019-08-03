// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdr_test

import (
	"reflect"
	"testing"

	"github.com/dlorch/nfsv3/xdr"
)

func TestEncodeNil(t *testing.T) {
	got, err := xdr.Marshal(nil)
	if err == nil {
		t.Fatalf("Expected error, but got %v", got)
	}
}

type Empty struct{}

var empty = &Empty{}

var emptyExpect = []byte{}

func TestEncodeEmpty(t *testing.T) {
	got, err := xdr.Marshal(empty)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(got) != 0 {
		t.Fatalf("Expected %v but got %v", emptyExpect, got)
	}
}

type Simple struct {
	Type   uint32
	Size   uint64
	hidden uint32
}

var simple = &Simple{
	Type: 1,
	Size: 5034543534,
}

var simpleExpect = []byte{0, 0, 0, 1, 0, 0, 0, 1, 44, 21, 9, 174}

func TestEncodeSimple(t *testing.T) {
	got, err := xdr.Marshal(simple)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, simpleExpect) {
		t.Fatalf("Expected %v but got %v", simpleExpect, got)
	}
}

type Nested struct {
	Value uint32
	Inner
}

type Inner struct {
	Value uint32
}

var nested = &Nested{
	Value: 12,
	Inner: Inner{
		Value: 13,
	},
}

var nestedExpect = []byte{0, 0, 0, 12, 0, 0, 0, 13}

func TestEncodeNested(t *testing.T) {
	got, err := xdr.Marshal(nested)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, nestedExpect) {
		t.Fatalf("Expected %v but got %v", nestedExpect, got)
	}
}

type FixedByteArray struct {
	Data [4]byte
}

var fixedByteArray = &FixedByteArray{
	Data: [4]byte{55, 43, 99, 102},
}

var fixedByteArrayExpect = []byte{55, 43, 99, 102}

func TestEncodeFixedByteArray(t *testing.T) {
	got, err := xdr.Marshal(fixedByteArray)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, fixedByteArrayExpect) {
		t.Fatalf("Expected %v but got %v", fixedByteArrayExpect, got)
	}
}

type DynamicallySizedValues struct {
	Data    []byte
	Values  []uint32
	Name    string
	Another string
}

var dynamicallySizedValues = &DynamicallySizedValues{
	Data:    []byte{41, 22, 13, 4, 15}, // encodes as: length + bytes + padding (total length must be multiple of four)
	Values:  []uint32{99, 33},          // encodes as: length + uint32 + no padding (length is already multiple of four)
	Name:    "Gopher",                  // encodes as: length + bytes + padding (total length must be multiple of four)
	Another: "lisp",                    // encodes as: length + bytes + no padding (length is already multiple of four)
}

var dynamicallySizedValuesExpect = []byte{0, 0, 0, 5, 41, 22, 13, 4, 15, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 99, 0, 0, 0, 33, 0, 0, 0, 6, 71, 111, 112, 104, 101, 114, 0, 0, 0, 0, 0, 4, 108, 105, 115, 112}

func TestEncodeDynamicallySizedValues(t *testing.T) {
	got, err := xdr.Marshal(dynamicallySizedValues)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, dynamicallySizedValuesExpect) {
		t.Fatalf("Expected %v but got %v", dynamicallySizedValuesExpect, got)
	}
}

type OptionalAttribute struct {
	AttributeFollows uint32 `xdr:"switch"`
	Attribute        Simple `xdr:"case=1"`
}

var optionalAttributeYes = &OptionalAttribute{
	AttributeFollows: 1,
	Attribute: Simple{
		Type: 12,
		Size: 33,
	},
}

var OptionalAttributeNo = &OptionalAttribute{
	AttributeFollows: 0,
}

var optionalAttributeYesExpect = []byte{0, 0, 0, 1, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0, 0, 33}

var optionalAttributeNoExpect = []byte{0, 0, 0, 0}

func TestEncodeOptionalAttributes(t *testing.T) {
	got, err := xdr.Marshal(optionalAttributeYes)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, optionalAttributeYesExpect) {
		t.Fatalf("Expected %v but got %v", optionalAttributeYesExpect, got)
	}

	got, err = xdr.Marshal(OptionalAttributeNo)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, optionalAttributeNoExpect) {
		t.Fatalf("Expected %v but got %v", optionalAttributeNoExpect, got)
	}
}

type Union struct {
	Status  uint32        `xdr:"switch"`
	Success SuccessResult `xdr:"case=0"`
	Failure FailureResult `xdr:"default"`
}

type SuccessResult struct {
	First  uint32
	Second uint32
}

type FailureResult struct {
	Error uint32
}

var unionSuccess = &Union{
	Status: 0,
	Success: SuccessResult{
		First:  44,
		Second: 36,
	},
}

var unionFailure = &Union{
	Status: 1,
	Failure: FailureResult{
		Error: 99,
	},
}

var unionSuccessExpect = []byte{0, 0, 0, 0, 0, 0, 0, 44, 0, 0, 0, 36}

var unionFailureExpect = []byte{0, 0, 0, 1, 0, 0, 0, 99}

func TestEncodeUnion(t *testing.T) {
	got, err := xdr.Marshal(unionSuccess)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, unionSuccessExpect) {
		t.Fatalf("Expected %v but got %v", unionSuccessExpect, got)
	}

	got, err = xdr.Marshal(unionFailure)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, unionFailureExpect) {
		t.Fatalf("Expected %v but got %v", unionFailureExpect, got)
	}
}

type InvalidCaseNoSwitch struct {
	First uint `xdr:"case=0"`
}

var invalidCaseNoSwitch = &InvalidCaseNoSwitch{
	First: 12,
}

func TestInvalidCaseNoSwitch(t *testing.T) {
	got, err := xdr.Marshal(invalidCaseNoSwitch)
	if err == nil {
		t.Fatalf("Expected error, but got %v", got)
	}
}

type InvalidDefaultNoSwitch struct {
	First uint `xdr:"default"`
}

var invalidDefaultNoSwitch = &InvalidDefaultNoSwitch{
	First: 12,
}

func TestInvalidDefaultNoSwitch(t *testing.T) {
	got, err := xdr.Marshal(invalidDefaultNoSwitch)
	if err == nil {
		t.Fatalf("Expected error, but got %v", got)
	}
}

type SwitchDefault struct {
	First  uint32 `xdr:"switch"`
	Second uint32 `xdr:"default"`
}

var switchDefault = &SwitchDefault{
	First:  12,
	Second: 44,
}

var switchDefaultExpect = []byte{0, 0, 0, 12, 0, 0, 0, 44}

func TestSwitchDefault(t *testing.T) {
	got, err := xdr.Marshal(switchDefault)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, switchDefaultExpect) {
		t.Fatalf("Expected %v but got %v", switchDefaultExpect, got)
	}
}

type SwitchSequence struct {
	First   uint32 `xdr:"switch"`
	Second  uint32 `xdr:"case=12"`
	Third   uint32 `xdr:"switch"`
	Fourth  uint32 `xdr:"case=3"`
	Fifth   uint32
	Sixth   uint32 `xdr:"case=5"`
	Seventh uint32
	Eight   uint32 `xdr:"case=12"`
	Ninth   uint32 `xdr:"default"`
}

var switchSequence = &SwitchSequence{
	First:   12,
	Second:  44,
	Third:   5,
	Fourth:  52,
	Fifth:   82,
	Sixth:   122,
	Seventh: 93,
	Eight:   22,
	Ninth:   11,
}

var switchSequenceExpect = []byte{0, 0, 0, 12, 0, 0, 0, 44, 0, 0, 0, 5, 0, 0, 0, 122, 0, 0, 0, 93}

// TestSwitchSequence verifies that two subsequent switch statements are executed correctly. Note that there is no
// nesting support for switch statements: a new switch statement overwrites the previous one. And also, there is
// no explicit "end switch" statement - a new switch statement followed by a default statement has to be used instead
func TestSwitchSequence(t *testing.T) {
	got, err := xdr.Marshal(switchSequence)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, switchSequenceExpect) {
		t.Fatalf("Expected %v but got %v", switchSequenceExpect, got)
	}
}

type UserLinkedList struct {
	ValueFollows uint32          `xdr:"switch"`
	Groups       GroupLinkedList `xdr:"case=1"`
	Next         interface{}
}

type GroupLinkedList struct {
	ValueFollows uint32 `xdr:"switch"`
	GroupID      uint32 `xdr:"case=1"`
	Next         interface{}
}

var userLinkedList = &UserLinkedList{
	ValueFollows: 1,
	Groups: GroupLinkedList{
		ValueFollows: 1,
		GroupID:      12,
		Next: GroupLinkedList{
			ValueFollows: 0,
		},
	},
	Next: UserLinkedList{
		ValueFollows: 0,
	},
}

var userLinkedListExpect = []byte{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0, 0, 0}

func TestUserLinkedList(t *testing.T) {
	got, err := xdr.Marshal(userLinkedList)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, userLinkedListExpect) {
		t.Fatalf("Expected %v but got %v", userLinkedListExpect, got)
	}
}
