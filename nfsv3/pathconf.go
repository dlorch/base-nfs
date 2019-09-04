// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nfsv3

// PathConf3Args (struct PATHCONF3args)
type PathConf3Args struct {
	FileHandle []byte
}

// PathConf3ResultOK (struct PATHCONF3resok)
type PathConf3ResultOK struct {
	PathConf3Result
	Objattributes   uint32 // TODO
	Linkmax         uint32
	Namemax         uint32
	Notrunc         uint32 // TODO bool
	Chownrestricted uint32 // TODO bool
	Caseinsensitive uint32 // TODO bool
	Casepreserving  uint32 // TODO bool
}

// TODO PATHCONF3resfail

// PathConf3Result (union PATHCONF3res)
type PathConf3Result struct {
	Status uint32
}

func nfsProcedure3PathConf(procedureArguments []byte) (interface{}, error) {
	// parse request
	// TODO

	// prepare result
	pathConfResult := &PathConf3ResultOK{
		PathConf3Result: PathConf3Result{
			Status: NFS3OK,
		},
		Objattributes:   0,
		Linkmax:         32000,
		Namemax:         255,
		Notrunc:         0,
		Chownrestricted: 1,
		Caseinsensitive: 0,
		Casepreserving:  1,
	}

	return pathConfResult, nil
}
