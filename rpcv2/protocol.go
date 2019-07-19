package rpcv2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

/*
	RPC: Remote Procedure Call Protocol Version 2 (RFC1057)

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

// RPCMsg describes the request/response RPC header (RFC1057: struct rpc_msg)
type RPCMsg struct {
	XID         uint32
	MessageType uint32
}

// CallBody describes the body of a CALL request (RFC1057: struct call_body)
type CallBody struct {
	RPCVersion     uint32
	Program        uint32
	ProgramVersion uint32
	Procedure      uint32
}

// ReplyBody describes the body of a REPLY to an RPC call (RFC1057: union reply_body)
type ReplyBody struct {
	ReplyStatus uint32
}

// AcceptedReplySuccess (RFC1057: struct accepted_reply)
type AcceptedReplySuccess struct {
	Verifier    OpaqueAuth
	AcceptState uint32
}

// OpaqueAuth describes the type of authentication used (RFC1057: struct opaque_auth)
type OpaqueAuth struct {
	Flavor uint32
	Length uint32
}

// Authentication (RFC1057: enum auth_flavor)
const (
	AuthenticationNull  uint32 = 0 // AUTH_NULL
	AuthenticationUNIX  uint32 = 1 // AUTH_UNIX
	AuthenticationShort uint32 = 2 // AUTH_SHORT
	AuthenticationDES   uint32 = 3 // AUTH_DES
)

// Constants for RPCv2 (RFC 1057)
const (
	Call                 uint32 = 0       // Message type CALL
	Reply                uint32 = 1       // Message type REPLY
	MessageAccepted      uint32 = 0       // Reply to a call message
	MessageDenied        uint32 = 1       // Reply to a call message
	Success              uint32 = 0       // RPC executed succesfully
	ProgramUnavailable   uint32 = 1       // remote hasn't exported program
	ProgramMismatch      uint32 = 2       // remote can't support version
	ProcedureUnavailable uint32 = 3       // program can't support procedure
	GarbageArguments     uint32 = 4       // procedure can't decode params
	LastFragment         uint32 = 1 << 31 // last fragment delimiter for record marking (RM)
)

// handleTCPClient handles TCP client connections, reads requests and delimits them into
// individual messages (RFC 1057: 10. Record Marking Standard) for further processing
func handleTCPClient(clientConnection net.Conn, rpcService *RPCService) {
	moreRequests := true

	requestBytes, isLastRequest, err := readNextRequest(clientConnection)

	for moreRequests {
		if err != nil {
			fmt.Println("Error: ", err.Error())
			// TODO send error back to client
		}

		request, err := parseRPCRequest(requestBytes)

		if err != nil {
			fmt.Println("Error: " + err.Error())
		}

		procedure := rpcService.Procedures[request.CallBody.Procedure] // TODO check for not existing procedures

		rpcResponse := procedure(&request)

		responseBytes, err := serializeRPCResponse(rpcResponse)

		if err != nil {
			fmt.Println("Error: " + err.Error())
			// TODO send error message back to client
		}

		// ---- fragments

		var fragmentBuffer = new(bytes.Buffer)
		lastFragment := uint32(1 << 31)
		fragmentLength := uint32(len(responseBytes))

		err = binary.Write(fragmentBuffer, binary.BigEndian, lastFragment|fragmentLength)

		// ---- end

		fragmentBuffer.Write(responseBytes)

		_, err = clientConnection.Write(fragmentBuffer.Bytes())

		if err != nil {
			fmt.Println("[mount] Error sending response: ", err.Error())
		}

		if isLastRequest {
			moreRequests = false
		}
	}

	clientConnection.Close()
}

// handleUDPClient handles UDP connections
func handleUDPClient(requestBytes []byte, serverConnection *net.UDPConn, clientAddress *net.UDPAddr, rpcService *RPCService) {
	request, err := parseRPCRequest(requestBytes)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		// TODO send error message back to client
	}

	procedure := rpcService.Procedures[request.CallBody.Procedure] // TODO check for not existing procedures
	rpcResponse := procedure(&request)

	responseBytes, err := serializeRPCResponse(rpcResponse)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		// TODO send error message back to client
	}

	serverConnection.WriteToUDP(responseBytes, clientAddress)
}

func serializeRPCResponse(rpcResponse *RPCResponse) (response []byte, err error) {
	var responseBuffer = new(bytes.Buffer)

	err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.RPCMessage)

	if err != nil {
		return response, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.ReplyBody)

	if err != nil {
		return response, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.AcceptedReplySuccess)

	if err != nil {
		return response, err
	}

	_, err = responseBuffer.Write(rpcResponse.ResponseBody)

	response = make([]byte, responseBuffer.Len())
	copy(response, responseBuffer.Bytes())
	return response, nil
}

func readNextRequest(clientConnection net.Conn) (requestBytes []byte, isLastRequest bool, err error) {
	var fragmentHeader uint32
	fragmentHeaderBytes := make([]byte, 4)

	_, err = clientConnection.Read(fragmentHeaderBytes)

	if err != nil {
		return requestBytes, isLastRequest, err
	}

	fragmentHeaderBytesBuffer := bytes.NewBuffer(fragmentHeaderBytes)

	err = binary.Read(fragmentHeaderBytesBuffer, binary.BigEndian, &fragmentHeader)

	if err != nil {
		return requestBytes, isLastRequest, err
	}

	isLastRequest = (fragmentHeader & LastFragment) != 0
	remainingFragmentLength := uint32(fragmentHeader & ^LastFragment)

	messageBytesBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)

	for remainingFragmentLength > 0 {
		bytesRead, err := clientConnection.Read(readBuffer)

		if err != nil {
			return requestBytes, isLastRequest, err
		}

		_, err = messageBytesBuffer.Write(readBuffer)

		if err != nil {
			return requestBytes, isLastRequest, err
		}

		remainingFragmentLength -= uint32(bytesRead)
	}

	requestBytes = make([]byte, messageBytesBuffer.Len())
	copy(requestBytes, messageBytesBuffer.Bytes())

	return requestBytes, isLastRequest, nil
}

func writeFragmentedReply(clientConnection net.Conn, messageBytes []byte, lastMessage bool) {
	//
}

func parseRPCRequest(rpcMessage []byte) (rpcRequest RPCRequest, err error) {
	var requestBuffer = bytes.NewBuffer(rpcMessage)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.RPCMessage)

	if err != nil {
		return rpcRequest, err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.CallBody)

	if err != nil {
		return rpcRequest, err
	}

	// TODO verify rpcRequest.CallBody.RPCVersion == 2

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.Credentials)

	if err != nil {
		return rpcRequest, err
	}

	rpcRequest.CredentialsBody = make([]byte, rpcRequest.Credentials.Length) // TODO insecure

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.CredentialsBody)

	if err != nil {
		return rpcRequest, err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.Verifier)

	if err != nil {
		return rpcRequest, err
	}

	rpcRequest.VerifierBody = make([]byte, rpcRequest.Verifier.Length) // TODO insecure

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.VerifierBody)

	if err != nil {
		return rpcRequest, err
	}

	rpcRequest.RequestBody = make([]byte, requestBuffer.Len())
	copy(rpcRequest.RequestBody, requestBuffer.Bytes())

	return rpcRequest, nil
}
