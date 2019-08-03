// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mountv3

import (
	"github.com/dlorch/nfsv3/rpcv2"
	"github.com/dlorch/nfsv3/xdr"
)

// MountProcedure3Export  is the number for this RPC procedure (MOUNTPROC3_EXPORT)
const MountProcedure3Export uint32 = 5

// Groups describes a linked-list of groups (struct groupnode)
type Groups struct {
	ValueFollows uint32 `xdr:"switch"`
	GrName       string `xdr:"case=1"`
	GrNext       interface{}
}

// Exports describes a linked-list of exports (struct exportnode)
type Exports struct {
	ValueFollows uint32 `xdr:"switch"`
	ExDir        string `xdr:"case=1"`
	ExGroups     Groups
	ExNext       interface{}
}

// ExportNode ...
type ExportNode struct{}

// ToBytes serializes the ExportNode to be sent back to the client
func (reply *ExportNode) ToBytes() ([]byte, error) {
	exports := &Exports{
		ValueFollows: 1,
		ExDir:        "/volume1/Public",
		ExGroups: Groups{
			ValueFollows: 1,
			GrName:       "*",
			GrNext: Groups{
				ValueFollows: 0,
			},
		},
		ExNext: Exports{
			ValueFollows: 0,
		},
	}

	return xdr.Marshal(exports)
}

// Export returns a list of all the exported file systems and which
// clients are allowed to mount each one.
// https://tools.ietf.org/html/rfc1813#page-113
func Export(procedureArguments []byte) (rpcv2.Serializable, error) {
	return &ExportNode{}, nil
}
