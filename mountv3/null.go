package mountv3

// VoidReply is an empty reply
type VoidReply struct{}

func mountProcedure3Null(procedureArguments []byte) (interface{}, error) {
	return &VoidReply{}, nil
}
