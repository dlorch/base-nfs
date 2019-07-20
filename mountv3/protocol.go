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

package mountv3

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dlorch/nfsv3/rpcv2"
)

// GroupNode (RFC1813: struct groupnode)
type GroupNode struct {
}

// ExportNode (RFC1813: struct exportnode)
type ExportNode struct {
}

// MountResult3 (RFC1813: struct mountres3)
type MountResult3 struct {
	Status         uint32
	MountResult3OK MountResult3OK
}

// MountResult3OK (RFC1813: struct mountres3_ok)
type MountResult3OK struct {
	FileHandle3 []byte
	AuthFlavors []uint32
}

// Constants for mount protocol (RFC1813)
const (
	Program                    uint32 = 100005 // Mount service program number
	Version                    uint32 = 3      // Mount service version
	MountProcedure3Null        uint32 = 0      // MOUNTPROC3_NULL
	MountProcedure3Mount       uint32 = 1      // MOUNTPROC3_MNT
	MountProcedure3Dump        uint32 = 2      // MOUNTPROC3_DUMP
	MountProcedure3Unmount     uint32 = 3      // MOUNTPROC3_UMNT
	MountProcedure3UnmountAll  uint32 = 4      // MOUNTPROC3_UMNTALL
	MountProcedure3Export      uint32 = 5      // MOUNTPROC3_EXPORT
	Mount3OK                   uint32 = 0      // MNT3_OK: no error
	Mount3ErrorPermissions     uint32 = 1      // MNT3ERR_PERM: Not owner
	Mount3ErrorNoEntry         uint32 = 2      // MNT3ERR_NOENT: No such file or directory
	Mount3ErrorIO              uint32 = 5      // MNT3ERR_IO: I/O error
	Mount3ErrorAccess          uint32 = 13     // MNT3ERR_ACCES: Permission denied
	Mount3ErrorNotDirectory    uint32 = 20     // MNT3ERR_NOTDIR: Not a directory
	Mount3ErrorInvalidArgument uint32 = 22     // MNT3ERR_INVAL: Invalid argument
	Mount3ErrorNameTooLong     uint32 = 63     // MNT3ERR_NAMETOOLONG: Filename too long
	Mount3ErrorNotSupported    uint32 = 10004  // MNT3ERR_NOTSUPP: Operation not supported
	Mount3ErrorServerFault     uint32 = 10006  // MNT3ERR_SERVERFAULT: A failure on the server
)

func mountProcedure3Null(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	response := &rpcv2.RPCResponse{
		RPCMessage: rpcv2.RPCMsg{
			XID:         request.RPCMessage.XID,
			MessageType: rpcv2.Reply,
		},
		ReplyBody: rpcv2.ReplyBody{
			ReplyStatus: rpcv2.MessageAccepted,
		},
	}

	verifier := rpcv2.OpaqueAuth{
		Flavor: rpcv2.AuthenticationNull,
		Length: 0,
	}

	if request.CallBody.ProgramVersion == Version {
		response.AcceptedReply = rpcv2.AcceptedReply{
			Verifier:    verifier,
			AcceptState: rpcv2.Success,
		}
	} else {
		response.AcceptedReply = rpcv2.AcceptedReply{
			Verifier:                verifier,
			AcceptState:             rpcv2.ProgramMismatch,
			LowestVersionSupported:  Version,
			HighestVersionSupported: Version,
		}
	}

	return response
}

func mountProcedure3export(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	var responseBuffer = new(bytes.Buffer)

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

	successReply := rpcv2.AcceptedReply{
		Verifier:    verifierReply,
		AcceptState: rpcv2.Success,
	}

	// --- mount service body

	var valueFollowsYes uint32
	valueFollowsYes = 1

	var valueFollowsNo uint32
	valueFollowsNo = 0

	directoryContents := "/volume1/Public"
	directoryLength := uint32(len(directoryContents))

	groupContents := "*"
	groupLength := uint32(len(groupContents))

	fillBytes := uint8(0)

	err := binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)
	err = binary.Write(responseBuffer, binary.BigEndian, &directoryLength)
	_, err = responseBuffer.Write([]byte(directoryContents))
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)

	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)
	err = binary.Write(responseBuffer, binary.BigEndian, &groupLength)
	_, err = responseBuffer.Write([]byte(groupContents))
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsNo)

	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsNo)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	// --- create response

	response := &rpcv2.RPCResponse{
		RPCMessage:    rpcMessage,
		ReplyBody:     replyBody,
		AcceptedReply: successReply,
	}

	response.AcceptedReply.Results = make([]byte, responseBuffer.Len())
	copy(response.AcceptedReply.Results, responseBuffer.Bytes())

	return response
}

func mountProcedure3mount(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	// parse request
	requestBuffer := bytes.NewBuffer(request.RequestBody)

	var dirPathLength uint32

	err := binary.Read(requestBuffer, binary.BigEndian, &dirPathLength)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	dirPathName := make([]byte, dirPathLength) // TODO unsafe?

	err = binary.Read(requestBuffer, binary.BigEndian, &dirPathName)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	// prepare result
	var resultBuffer = new(bytes.Buffer)

	err = binary.Write(resultBuffer, binary.BigEndian, Mount3OK)

	err = binary.Write(resultBuffer, binary.BigEndian, uint32(4))  // length of file handle
	err = binary.Write(resultBuffer, binary.BigEndian, uint32(42)) // file handle

	err = binary.Write(resultBuffer, binary.BigEndian, uint32(1))                // number of auth flavors
	err = binary.Write(resultBuffer, binary.BigEndian, rpcv2.AuthenticationUNIX) // allowed flavors

	// create response
	response := &rpcv2.RPCResponse{
		RPCMessage: rpcv2.RPCMsg{
			XID:         request.RPCMessage.XID,
			MessageType: rpcv2.Reply,
		},
		ReplyBody: rpcv2.ReplyBody{
			ReplyStatus: rpcv2.MessageAccepted,
		},
		AcceptedReply: rpcv2.AcceptedReply{
			Verifier: rpcv2.OpaqueAuth{
				Flavor: rpcv2.AuthenticationNull,
				Length: 0,
			},
			AcceptState: rpcv2.Success,
		},
	}

	response.AcceptedReply.Results = make([]byte, resultBuffer.Len())
	copy(response.AcceptedReply.Results, resultBuffer.Bytes())

	return response
}
