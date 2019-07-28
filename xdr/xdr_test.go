package xdr_test

import (
	"reflect"
	"testing"

	"github.com/dlorch/nfsv3/xdr"
)

func TestEncodeVoid(t *testing.T) {
	got, err := xdr.Marshal(nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(got) != 0 {
		t.Fatalf("Expected %v but got %v", []byte{}, got)
	}
}

type Empty struct{}

var empty = &Empty{}

func TestEncodeEmpty(t *testing.T) {
	got, err := xdr.Marshal(empty)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(got) != 0 {
		t.Fatalf("Expected %v but got %v", []byte{}, got)
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
	Data []byte
	Name string
}

var dynamicallySizedValues = &DynamicallySizedValues{
	Data: []byte{41, 22, 13, 4, 15}, // encodes as: length + bytes
	Name: "Gopher",                  // encodes as: length + bytes + padding (total length must be multiple of four)
}

var dynamicallySizedValuesExpect = []byte{0, 0, 0, 5, 41, 22, 13, 4, 15, 6, 0, 0, 0, 6, 71, 111, 112, 104, 101, 114, 0, 0}

func TestEncodeDynamicallySizedValues(t *testing.T) {
	got, err := xdr.Marshal(dynamicallySizedValues)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, dynamicallySizedValuesExpect) {
		t.Fatalf("Expected %v but got %v", dynamicallySizedValuesExpect, got)
	}
}

type SizeLimit struct {
	Data []byte `xdr:"maxsize=5"`
}

var sizeLimitInputBytes = []byte{1, 2, 3, 4, 5, 6}

var sizeLimitExpect = &SizeLimit{
	Data: []byte{1, 2, 3, 4, 5},
}

func TestDecodeSizeLimit(t *testing.T) {
	t.Error("Unimplemented")
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
