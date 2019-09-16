// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdr_test

import (
	"reflect"
	"testing"

	"github.com/dlorch/base-nfs/xdr"
)

func TestEncodeNil(t *testing.T) {
	got, err := xdr.Marshal(nil)
	if err == nil {
		t.Fatalf("Expected error, but got %v", got)
	}
}

func TestDecodeNil(t *testing.T) {
	_, err := xdr.Unmarshal([]byte{}, nil)
	if err == nil {
		t.Fatalf("Expected error for nil")
	}
}

func TestDecodeNonPtr(t *testing.T) {
	_, err := xdr.Unmarshal([]byte{}, 42)
	if err == nil {
		t.Fatalf("Expected error for non-pointer value")
	}
}

type Empty struct{}

var empty = &Empty{}

var emptyBytes = []byte{}

func TestEncodeEmpty(t *testing.T) {
	got, err := xdr.Marshal(empty)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(got) != 0 {
		t.Fatalf("Expected %v but got %v", emptyBytes, got)
	}
}

func TestDecodeEmpty(t *testing.T) {
	var got Empty
	_, err := xdr.Unmarshal(emptyBytes, &got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(&got, empty) {
		t.Fatalf("Expected %v but got %v", empty, got)
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

var simpleBytes = []byte{0, 0, 0, 1, 0, 0, 0, 1, 44, 21, 9, 174}

func TestEncodeSimple(t *testing.T) {
	got, err := xdr.Marshal(simple)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, simpleBytes) {
		t.Fatalf("Expected %v but got %v", simpleBytes, got)
	}
}

func TestDecodeSimple(t *testing.T) {
	got := &Simple{}
	_, err := xdr.Unmarshal(simpleBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, simple) {
		t.Fatalf("Expected %v but got %v", simple, got)
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

var nestedBytes = []byte{0, 0, 0, 12, 0, 0, 0, 13}

func TestEncodeNested(t *testing.T) {
	got, err := xdr.Marshal(nested)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, nestedBytes) {
		t.Fatalf("Expected %v but got %v", nestedBytes, got)
	}
}

func TestDecodeNested(t *testing.T) {
	got := &Nested{}
	_, err := xdr.Unmarshal(nestedBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, nested) {
		t.Fatalf("Expected %v but got %v", nested, got)
	}
}

type FixedByteArray struct {
	Data [4]byte
}

var fixedByteArray = &FixedByteArray{
	Data: [4]byte{55, 43, 99, 102},
}

var fixedByteArrayBytes = []byte{55, 43, 99, 102}

func TestEncodeFixedByteArray(t *testing.T) {
	got, err := xdr.Marshal(fixedByteArray)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, fixedByteArrayBytes) {
		t.Fatalf("Expected %v but got %v", fixedByteArrayBytes, got)
	}
}

func TestDecodeFixedByteArray(t *testing.T) {
	got := &FixedByteArray{}
	_, err := xdr.Unmarshal(fixedByteArrayBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, fixedByteArray) {
		t.Fatalf("Expected %v but got %v", fixedByteArray, got)
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

var dynamicallySizedValuesBytes = []byte{0, 0, 0, 5, 41, 22, 13, 4, 15, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 99, 0, 0, 0, 33, 0, 0, 0, 6, 71, 111, 112, 104, 101, 114, 0, 0, 0, 0, 0, 4, 108, 105, 115, 112}

func TestEncodeDynamicallySizedValues(t *testing.T) {
	got, err := xdr.Marshal(dynamicallySizedValues)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, dynamicallySizedValuesBytes) {
		t.Fatalf("Expected %v but got %v", dynamicallySizedValuesBytes, got)
	}
}

func TestDecodeDynamicallySizedValues(t *testing.T) {
	got := &DynamicallySizedValues{}
	_, err := xdr.Unmarshal(dynamicallySizedValuesBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, dynamicallySizedValues) {
		t.Fatalf("Expected %v but got %v", dynamicallySizedValues, got)
	}
}

var invalidSizedValuesBytes = []byte{0, 0, 0, 3, 1, 2}

func TestInvalidSizedValues(t *testing.T) {
	got := &DynamicallySizedValues{}
	_, err := xdr.Unmarshal(invalidSizedValuesBytes, got)
	if err == nil {
		t.Fatalf("Expected error, but got %v", got)
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

var optionalAttributeNo = &OptionalAttribute{
	AttributeFollows: 0,
}

var optionalAttributeYesBytes = []byte{0, 0, 0, 1, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0, 0, 33}

var optionalAttributeNoBytes = []byte{0, 0, 0, 0}

func TestEncodeOptionalAttributes(t *testing.T) {
	got, err := xdr.Marshal(optionalAttributeYes)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, optionalAttributeYesBytes) {
		t.Fatalf("Expected %v but got %v", optionalAttributeYesBytes, got)
	}

	got, err = xdr.Marshal(optionalAttributeNo)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, optionalAttributeNoBytes) {
		t.Fatalf("Expected %v but got %v", optionalAttributeNoBytes, got)
	}
}

func TestDecodeOptionalAttributes(t *testing.T) {
	got := &OptionalAttribute{}
	_, err := xdr.Unmarshal(optionalAttributeYesBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, optionalAttributeYes) {
		t.Fatalf("Expected %v but got %v", optionalAttributeYes, got)
	}

	got = &OptionalAttribute{}
	_, err = xdr.Unmarshal(optionalAttributeNoBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, optionalAttributeNo) {
		t.Fatalf("Expected %v but got %v", optionalAttributeNo, got)
	}
}

type MultipleCase struct {
	Mode      uint32 `xdr:"switch"`
	Attribute uint32 `xdr:"case=0,1"`
	Verifier  uint32 `xdr:"case=2"`
}

var multipleCaseZero = &MultipleCase{
	Mode:      0,
	Attribute: 12,
}

var multipleCaseOne = &MultipleCase{
	Mode:      1,
	Attribute: 44,
}

var multipleCaseTwo = &MultipleCase{
	Mode:     2,
	Verifier: 35,
}

var multipleCaseThree = &MultipleCase{
	Mode: 3,
}

var multipleCaseZeroBytes = []byte{0, 0, 0, 0, 0, 0, 0, 12}

var multipleCaseOneBytes = []byte{0, 0, 0, 1, 0, 0, 0, 44}

var multipleCaseTwoBytes = []byte{0, 0, 0, 2, 0, 0, 0, 35}

var multipleCaseThreeBytes = []byte{0, 0, 0, 3}

func TestEncodeMultipleCase(t *testing.T) {
	got, err := xdr.Marshal(multipleCaseZero)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseZeroBytes) {
		t.Fatalf("Expected %v but got %v", multipleCaseZeroBytes, got)
	}

	got, err = xdr.Marshal(multipleCaseOne)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseOneBytes) {
		t.Fatalf("Expected %v but got %v", multipleCaseOneBytes, got)
	}

	got, err = xdr.Marshal(multipleCaseTwo)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseTwoBytes) {
		t.Fatalf("Expected %v but got %v", multipleCaseTwoBytes, got)
	}

	got, err = xdr.Marshal(multipleCaseThree)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseThreeBytes) {
		t.Fatalf("Expected %v but got %v", multipleCaseThreeBytes, got)
	}
}

func TestDecodeMultipleCase(t *testing.T) {
	got := &MultipleCase{}
	_, err := xdr.Unmarshal(multipleCaseZeroBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseZero) {
		t.Fatalf("Expected %v but got %v", multipleCaseZero, got)
	}

	got = &MultipleCase{}
	_, err = xdr.Unmarshal(multipleCaseOneBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseOne) {
		t.Fatalf("Expected %v but got %v", multipleCaseOne, got)
	}

	got = &MultipleCase{}
	_, err = xdr.Unmarshal(multipleCaseTwoBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseTwo) {
		t.Fatalf("Expected %v but got %v", multipleCaseTwo, got)
	}

	got = &MultipleCase{}
	_, err = xdr.Unmarshal(multipleCaseThreeBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, multipleCaseThree) {
		t.Fatalf("Expected %v but got %v", multipleCaseThree, got)
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

var unionSuccessBytes = []byte{0, 0, 0, 0, 0, 0, 0, 44, 0, 0, 0, 36}

var unionFailureBytes = []byte{0, 0, 0, 1, 0, 0, 0, 99}

func TestEncodeUnion(t *testing.T) {
	got, err := xdr.Marshal(unionSuccess)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, unionSuccessBytes) {
		t.Fatalf("Expected %v but got %v", unionSuccessBytes, got)
	}

	got, err = xdr.Marshal(unionFailure)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, unionFailureBytes) {
		t.Fatalf("Expected %v but got %v", unionFailureBytes, got)
	}
}

func TestDecodeUnion(t *testing.T) {
	got := &Union{}
	_, err := xdr.Unmarshal(unionSuccessBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, unionSuccess) {
		t.Fatalf("Expected %v but got %v", unionSuccess, got)
	}

	got = &Union{}
	_, err = xdr.Unmarshal(unionFailureBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, unionFailure) {
		t.Fatalf("Expected %v but got %v", unionFailure, got)
	}
}

type InvalidCaseNoSwitch struct {
	First uint `xdr:"case=0"`
}

var invalidCaseNoSwitch = &InvalidCaseNoSwitch{
	First: 12,
}

func TestEncodeInvalidCaseNoSwitch(t *testing.T) {
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

func TestEncodeInvalidDefaultNoSwitch(t *testing.T) {
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

var switchDefaultBytes = []byte{0, 0, 0, 12, 0, 0, 0, 44}

func TestEncodeSwitchDefault(t *testing.T) {
	got, err := xdr.Marshal(switchDefault)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, switchDefaultBytes) {
		t.Fatalf("Expected %v but got %v", switchDefaultBytes, got)
	}
}

func TestDecodeSwitchDefault(t *testing.T) {
	got := &SwitchDefault{}
	_, err := xdr.Unmarshal(switchDefaultBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, switchDefault) {
		t.Fatalf("Expected %v but got %v", switchDefault, got)
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

var switchSequenceBytes = []byte{0, 0, 0, 12, 0, 0, 0, 44, 0, 0, 0, 5, 0, 0, 0, 122, 0, 0, 0, 93}

var switchSequenceDecoded = &SwitchSequence{
	First:   12,
	Second:  44,
	Third:   5,
	Sixth:   122,
	Seventh: 93,
}

// TestSwitchSequence verifies that two subsequent switch statements are executed correctly. Note that there is no
// nesting support for switch statements: a new switch statement overwrites the previous one. And also, there is
// no explicit "end switch" statement - a new switch statement followed by a default statement has to be used instead
func TestEncodeSwitchSequence(t *testing.T) {
	got, err := xdr.Marshal(switchSequence)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, switchSequenceBytes) {
		t.Fatalf("Expected %v but got %v", switchSequenceBytes, got)
	}
}

func TestDecodeSwitchSequence(t *testing.T) {
	got := &SwitchSequence{}
	_, err := xdr.Unmarshal(switchSequenceBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, switchSequenceDecoded) {
		t.Fatalf("Expected %v but got %v", switchSequenceDecoded, got)
	}
}

type UserLinkedList struct {
	ValueFollows uint32           `xdr:"switch"`
	Groups       *GroupLinkedList `xdr:"case=1"`
	Next         *UserLinkedList
}

type GroupLinkedList struct {
	ValueFollows uint32 `xdr:"switch"`
	GroupID      uint32 `xdr:"case=1"`
	Next         *GroupLinkedList
}

var userLinkedList = &UserLinkedList{
	ValueFollows: 1,
	Groups: &GroupLinkedList{
		ValueFollows: 1,
		GroupID:      12,
		Next: &GroupLinkedList{
			ValueFollows: 0,
		},
	},
	Next: &UserLinkedList{
		ValueFollows: 0,
	},
}

var userLinkedListBytes = []byte{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0, 0, 0}

func TestUserLinkedList(t *testing.T) {
	got, err := xdr.Marshal(userLinkedList)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, userLinkedListBytes) {
		t.Fatalf("Expected %v but got %v", userLinkedListBytes, got)
	}
}

func TestDecodeUserLinkedList(t *testing.T) {
	got := &UserLinkedList{}
	_, err := xdr.Unmarshal(userLinkedListBytes, got)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, userLinkedList) {
		t.Fatalf("Expected %v but got %v", userLinkedList, got)
	}
}
