package nfsv3

import "github.com/dlorch/nfsv3/rpcv2"

// PathConf3Args (struct PATHCONF3args)
type PathConf3Args struct {
	FileHandle []byte
}

// PathConf3ResultOK (struct PATHCONF3resok)
type PathConf3ResultOK struct {
	PathConf3Result
	objattributes   uint32 // TODO
	linkmax         uint32
	namemax         uint32
	notrunc         uint32 // TODO bool
	chownrestricted uint32 // TODO bool
	caseinsensitive uint32 // TODO bool
	casepreserving  uint32 // TODO bool
}

// TODO PATHCONF3resfail

// PathConf3Result (union PATHCONF3res)
type PathConf3Result struct {
	status uint32
}

// ToBytes serializes the PathConf3ResultOK to be sent back to the client
func (reply *PathConf3ResultOK) ToBytes() ([]byte, error) {
	return rpcv2.SerializeFixedSizeStruct(reply)
}

func nfsProcedure3PathConf(procedureArguments []byte) (rpcv2.Serializable, error) {
	// parse request
	// TODO

	// prepare result
	pathConfResult := &PathConf3ResultOK{
		PathConf3Result: PathConf3Result{
			status: NFS3OK,
		},
		objattributes:   0,
		linkmax:         32000,
		namemax:         255,
		notrunc:         0,
		chownrestricted: 1,
		caseinsensitive: 0,
		casepreserving:  1,
	}

	return pathConfResult, nil
}
