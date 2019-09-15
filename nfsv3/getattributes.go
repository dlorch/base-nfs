// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nfsv3

// GetAttr3Args (struct FSINFOargs)
type GetAttr3Args struct {
	FileHandle []byte
}

// GetAttr3ResultOK (struct GETATTR3resok)
type GetAttr3ResultOK struct {
	GetAttr3Result
	ObjectAttributes FAttr3
}

// GetAttr3Result (union GETATTR3res)
type GetAttr3Result struct {
	Status uint32
}

func nfsProcedure3GetAttributes(procedureArguments []byte) (interface{}, error) {
	// parse request
	// TODO

	// prepare result
	getAttrResult := &GetAttr3ResultOK{
		GetAttr3Result: GetAttr3Result{
			Status: NFS3OK,
		},
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
	}

	return getAttrResult, nil
}
