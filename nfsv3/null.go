package nfsv3

// VoidReply is an empty reply
type VoidReply struct{}

func nfsProcedure3Null(procedureArguments []byte) (interface{}, error) {
	return &VoidReply{}, nil
}
