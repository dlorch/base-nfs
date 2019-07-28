package xdr_test

import (
	"reflect"
	"testing"

	"github.com/dlorch/nfsv3/xdr"
)

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
