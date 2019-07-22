/*
	NFS Version 3 Protocol Specification (RFC1813)

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

package nfsv3

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dlorch/nfsv3/rpcv2"
)

// FSInfoArgs (struct FSINFOargs)
type FSInfoArgs struct {
	FileHandle []byte
}

// FSInfo3ResultOK (struct FSINFO3resok)
type FSInfo3ResultOK struct {
	FSInfo3Result
	objattributes        uint32 // TODO
	rtmax                uint32
	rtpref               uint32
	rtmult               uint32
	wtmax                uint32
	wtpref               uint32
	wtmult               uint32
	dtpref               uint32
	maxfilesize          uint64
	timedeltaseconds     uint32
	timedeltananoseconds uint32
	properties           uint32
}

// FSInfo3ResultFail (struct FSINFO3resfail)
type FSInfo3ResultFail struct {
	FSInfo3Result
	// TODO post_op_attr obj_attributes
}

// FSInfo3Result (union FSINFO3res)
type FSInfo3Result struct {
	status uint32
}

// Constants for mount protocol (RFC1813)
const (
	Program                    uint32 = 100003 // Mount service program number
	Version                    uint32 = 3      // Mount service program version
	NFSProcedure3Null          uint32 = 0      // NFSPROC3_NULL
	NFSProcedure3GetAttributes uint32 = 1      // NFSPROC3_GETATTR
	NFSProcedure3SetAttributes uint32 = 2      // NFSPROC3_SETATTR
	NFSProcedure3Lookup        uint32 = 3      // NFSPROC3_LOOKUP
	NFSProcedure3Access        uint32 = 4      // NFSPROC3_ACCESS
	NFSProcedure3Readlink      uint32 = 5      // NFSPROC3_READLINK
	NFSProcedure3Read          uint32 = 6      // NFSPROC3_READ
	NFSProcedure3Write         uint32 = 7      // NFSPROC3_WRITE
	NFSProcedure3Create        uint32 = 8      // NFSPROC3_CREATE
	NFSProcedure3MkDir         uint32 = 9      // NFSPROC3_MKDIR
	NFSProcedure3Symlink       uint32 = 10     // NFSPROC3_SYMLINK
	NFSProcedure3MkNod         uint32 = 11     // NFSPROC3_MKNOD
	NFSProcedure3Remove        uint32 = 12     // NFSPROC3_REMOVE
	NFSProcedure3RmDir         uint32 = 13     // NFSPROC3_RMDIR
	NFSProcedure3Rename        uint32 = 14     // NFSPROC3_RENAME
	NFSProcedure3Link          uint32 = 15     // NFSPROC3_LINK
	NFSProcedure3ReadDir       uint32 = 16     // NFSPROC3_READDIR
	NFSProcedure3ReadDirPlus   uint32 = 17     // NFSPROC3_READDIRPLUS
	NFSProcedure3FSStat        uint32 = 18     // NFSPROC3_FSSTAT
	NFSProcedure3FSInfo        uint32 = 19     // NFSPROC3_FSINFO
	NFSProcedure3PathConf      uint32 = 20     // NFSPROC3_PATHCONF
	NFSProcedure3Commint       uint32 = 21     // NFSPROC3_COMMIT
	NFS3OK                     uint32 = 0      // NFS3_OK
)

func nfsProcedure3Null(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
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

func nfsProcedure3GetAttributes(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	fmt.Println("nfsProcedure3GetAttributes")
	return nil
}

func nfsProcedure3Lookup(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	fmt.Println("nfsProcedure3Access")
	return nil
}

func nfsProcedure3Access(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	fmt.Println("nfsProcedure3Access")
	return nil
}

func nfsProcedure3FSInfo(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	// parse request
	requestBuffer := bytes.NewBuffer(request.RequestBody)

	var fileHandleLength uint32

	err := binary.Read(requestBuffer, binary.BigEndian, &fileHandleLength)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	fsInfoArgs := FSInfoArgs{
		FileHandle: make([]byte, fileHandleLength), // TODO unsafe?
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &fsInfoArgs.FileHandle)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	// prepare result
	fsInfoResult := FSInfo3ResultOK{
		FSInfo3Result: FSInfo3Result{
			status: NFS3OK,
		},
		objattributes:        0,
		rtmax:                131072,
		rtpref:               131072,
		rtmult:               4096,
		wtmax:                131072,
		wtpref:               131072,
		wtmult:               4096,
		dtpref:               4096,
		maxfilesize:          8796093022207,
		timedeltaseconds:     1,
		timedeltananoseconds: 0,
		properties:           0x0000001b,
	}

	rpcResponse := &rpcv2.RPCResponse{
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

	// create response
	var resultBuffer = new(bytes.Buffer)

	err = binary.Write(resultBuffer, binary.BigEndian, &fsInfoResult)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	rpcResponse.AcceptedReply.Results = make([]byte, resultBuffer.Len())
	copy(rpcResponse.AcceptedReply.Results, resultBuffer.Bytes())

	return rpcResponse
}

func nfsProcedure3PathConf(request *rpcv2.RPCRequest) *rpcv2.RPCResponse {
	fmt.Println("nfsProcedure3PathConf")
	return nil
}
