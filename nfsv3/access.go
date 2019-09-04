// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nfsv3

// Access3ResultOK (union ACCESS3resok)
type Access3ResultOK struct {
	Access3Result
	PostOperationAttributes
	Access uint32
}

// Access3Result (union ACCESS3res)
type Access3Result struct {
	Status uint32
}

func nfsProcedure3Access(procedureArguments []byte) (interface{}, error) {
	// prepare result
	fsInfoResult := &Access3ResultOK{
		Access3Result: Access3Result{
			Status: NFS3OK,
		},
		PostOperationAttributes: PostOperationAttributes{
			AttributesFollow: 1,
			ObjectAttributes: FileAttr3{
				Typ:              2,
				Mode:             040777,
				Nlink:            4,
				UID:              0,
				GID:              0,
				Size:             4096,
				Used:             8192,
				Specdata1:        0,
				Specdata2:        0,
				Fsid:             0x388e4346cfc706a8,
				Fileid:           16,
				Atimeseconds:     1563137262,
				Atimenanoseconds: 460002975,
				Mtimeseconds:     1537128120,
				Mtimenanoseconds: 839607220,
				Ctimeseconds:     1537128120,
				Ctimenanoseconds: 839607220,
			},
		},
		Access: 0x1f,
	}

	return fsInfoResult, nil
}
