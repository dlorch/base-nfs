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
	Type uint32
	Size uint64
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

var dynamicallySizedValuesExpect = []byte{5, 41, 22, 13, 4, 15, 6, 0, 0, 0, 6, 71, 111, 112, 104, 101, 114, 0, 0}

func TestEncodeDynamicallySizedValues(t *testing.T) {
	got, err := xdr.Marshal(dynamicallySizedValues)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(got, dynamicallySizedValuesExpect) {
		t.Fatalf("Expected %v but got %v", dynamicallySizedValuesExpect, got)
	}
}
