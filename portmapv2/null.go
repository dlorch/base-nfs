package portmapv2

// VoidReply is an empty reply
type VoidReply struct{}

func procedureNull(procedureArguments []byte) (interface{}, error) {
	return &VoidReply{}, nil
}
