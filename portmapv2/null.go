package portmapv2

import "github.com/dlorch/nfsv3/rpcv2"

// VoidReply is an empty reply
type VoidReply struct{}

// ToBytes serializes the VoidReply to be sent back to the client
func (reply *VoidReply) ToBytes() ([]byte, error) {
	return rpcv2.SerializeFixedSizeStruct(reply)
}

func procedureNull(procedureArguments []byte) (rpcv2.Serializable, error) {
	return &VoidReply{}, nil
}
