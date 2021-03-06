// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mountv3

import (
	"bytes"
	"encoding/binary"

	"github.com/dlorch/base-nfs/rpcv2"
)

// MountProcedure3Mnt is the number for this RPC procedure (MOUNTPROC3_MNT)
const MountProcedure3Mnt uint32 = 1

// MountRes3OK (struct mountres3_ok)
type MountRes3OK struct {
	FHandle     []byte
	AuthFlavors []uint32
}

// MountRes3 (struct mountres3)
type MountRes3 struct {
	FhsStatus uint32      `xdr:"switch"`
	MountInfo MountRes3OK `xdr:"case=0"`
}

// Mnt maps a pathname on the server to a file handle.
// https://tools.ietf.org/html/rfc1813#page-109
func Mnt(procedureArguments []byte) (interface{}, error) {
	// parse request
	requestBuffer := bytes.NewBuffer(procedureArguments)

	var dirPathLength uint32

	err := binary.Read(requestBuffer, binary.BigEndian, &dirPathLength)

	if err != nil {
		return nil, err
	}

	dirPathName := make([]byte, dirPathLength) // TODO check MNTPATHLEN

	err = binary.Read(requestBuffer, binary.BigEndian, &dirPathName)

	if err != nil {
		return nil, err
	}

	mountOk := &MountRes3{
		FhsStatus: Mount3OK,
		MountInfo: MountRes3OK{
			FHandle:     []byte{0, 0, 0, 42},
			AuthFlavors: []uint32{rpcv2.AuthenticationUNIX},
		},
	}

	return mountOk, nil
}
