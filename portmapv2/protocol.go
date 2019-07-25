/*
	Port Mapper Protocol Specification Version 2 (RFC1057)

	BSD 2-Clause License

	Copyright (c) 2019, Daniel Lorch
	All rights reserved.

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are met:

	1. Redistributions of source code must retain the above copyright notice, this
	   list of conditions and the following disclaimer.

	2. Redistributions in binary form must reproduce the above copyright notice,
       this list of conditions and the following disclaimer in the documentation
       and/or other materials provided with the distribution.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package portmapv2

import (
	"bytes"
	"encoding/binary"

	"github.com/dlorch/nfsv3/rpcv2"
)

// PortmapVoidReply is an empty reply
type PortmapVoidReply struct{}

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

// Constants for port mapper
const (
	Program                 uint32 = 100000 // Portmap service program number (PMAP_PROG)
	Version                 uint32 = 2      // Portmap service version number
	PortmapProcedureNull    uint32 = 0      // PMAPPROC_NULL
	PortmapProcedureSet     uint32 = 1      // PMAPPROC_SET
	PortmapProcedureUnset   uint32 = 2      // PMAPPROC_UNSET
	PortmapProcedureGetPort uint32 = 3      // PMAPPROC_GETPORT
	PortmapProcedureDump    uint32 = 4      // PMAPPROC_DUMP
	PortmapProcedureCallIt  uint32 = 5      // PMAPPROC_CALLIT
	IPProtocolTCP           uint32 = 6      // protocol number for TCP/IP
	IPProtocolUDP           uint32 = 17     // protocol number for UCP/IP
	ProgramNotAvailable     uint32 = 0      // Port value of zero means the program has not been registered
)

// ----- PortmapProcedureNull

// ToBytes serializes the VoidReply to be sent back to the client
func (reply *PortmapVoidReply) ToBytes() ([]byte, error) {
	return []byte{}, nil
}

func procedureNull(procedureArguments []byte) (rpcv2.Serializable, error) {
	return &PortmapVoidReply{}, nil
}

// ----- PortmapProcedureGetPort

// ToBytes serializes the uint32 reply to be sent back to the client
func (reply *GetPortResult) ToBytes() ([]byte, error) {
	var responseBytes []byte
	responseBuffer := new(bytes.Buffer)
	err := binary.Write(responseBuffer, binary.BigEndian, &reply.Port)

	if err != nil {
		return responseBytes, err
	}

	responseBytes = make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())

	return responseBytes, err
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
