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
			ObjectAttributes: FAttr3{
				Type:  2,
				Mode:  040777,
				Nlink: 4,
				UID:   0,
				GID:   0,
				Size:  4096,
				Used:  8192,
				RDev: SpecData3{
					SpecData1: 0,
					SpecData2: 0,
				},
				FSID:   0x388e4346cfc706a8,
				FileID: 16,
				ATime: NFSTime3{
					Seconds:  1563137262,
					NSeconds: 460002975,
				},
				MTime: NFSTime3{
					Seconds:  1537128120,
					NSeconds: 839607220,
				},
				CTime: NFSTime3{
					Seconds:  1537128120,
					NSeconds: 839607220,
				},
			},
		},
		Access: 0x1f,
	}

	return fsInfoResult, nil
}
