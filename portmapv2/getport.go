package portmapv2

import (
	"bytes"
	"encoding/binary"

	"github.com/dlorch/nfsv3/rpcv2"
)

// GetPortResult represents the requested port number
type GetPortResult struct {
	Port uint32
}

// Mapping of (program, version, protocol) to port number (RFC1057: struct_mapping)
type Mapping struct {
	Program  uint32
	Version  uint32
	Protocol uint32
	Port     uint32
}

// ToBytes serializes the uint32 reply to be sent back to the client
func (reply *GetPortResult) ToBytes() ([]byte, error) {
	return rpcv2.SerializeFixedSizeStruct(reply)
}

func procedureGetPort(procedureArguments []byte) (rpcv2.Serializable, error) {
	var requestBody = bytes.NewBuffer(procedureArguments)
	var mapping Mapping

	err := binary.Read(requestBody, binary.BigEndian, &mapping)

	if err != nil {
		return &GetPortResult{Port: ProgramNotAvailable}, err
	}

	// TODO check callBody.Version == portmapv2.Version

	port := getPort(mapping)

	return &GetPortResult{Port: port}, nil
}
