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
	"fmt"

	"github.com/dlorch/nfsv3/rpcv2"
)

// Mapping of (program, version, protocol) to port number (RFC1057: struct_mapping)
type Mapping struct {
	Program  uint32
	Version  uint32
	Protocol uint32
	Port     uint32
}

// Constants for port mapper
const (
	Program          uint32 = 100000 // Portmap service program number (PMAP_PROG)
	Version          uint32 = 2      // Portmap service version number
	ProcedureNull    uint32 = 0      // PMAPPROC_NULL
	ProcedureGetPort uint32 = 3      // PMAPPROC_GETPORT
	IPProtocolTCP    uint32 = 6      // protocol number for TCP/IP
	IPProtocolUDP    uint32 = 17     // protocol number for UCP/IP
)

func procedureNull(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	return nil
}

func procedureGetPort(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	var requestBody = bytes.NewBuffer(request.RequestBody)
	var mapping Mapping

	err := binary.Read(requestBody, binary.BigEndian, &mapping)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO: send error message back to client
	}

	rpcMessage := rpcv2.RPCMsg{
		XID:         request.RPCMessage.XID,
		MessageType: rpcv2.Reply,
	}

	replyBody := rpcv2.ReplyBody{
		ReplyStatus: rpcv2.MessageAccepted,
	}

	verifierReply := rpcv2.OpaqueAuth{
		Flavor: rpcv2.AuthenticationNull,
		Length: 0,
	}

	successReply := rpcv2.AcceptedReplySuccess{
		Verifier:    verifierReply,
		AcceptState: rpcv2.Success,
	}

	// TODO check callBody.Version == portmapv2.Version

	result, err := getPort(mapping)

	var responseBuffer = new(bytes.Buffer)

	err = binary.Write(responseBuffer, binary.BigEndian, &result) // TODO check err

	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	response := &rpcv2.RPCResponse{
		RPCMessage:           rpcMessage,
		ReplyBody:            replyBody,
		AcceptedReplySuccess: successReply,
	}

	response.ResponseBody = make([]byte, responseBuffer.Len())
	copy(response.ResponseBody, responseBuffer.Bytes())

	return response
}
