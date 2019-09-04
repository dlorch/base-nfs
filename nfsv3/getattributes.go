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
	ObjectAttributes FileAttr3
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
	}

	return getAttrResult, nil
}
