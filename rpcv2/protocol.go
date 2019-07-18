package rpcv2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/dlorch/nfsv3/mountv3"
	"github.com/dlorch/nfsv3/portmapv2"
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
func HandleTCPClient(clientConnection net.Conn) {
	moreRequests := true

	requestBytes, isLastRequest, err := readNextRequest(clientConnection)

	for moreRequests {
		if err != nil {
			fmt.Println("Error: ", err.Error())
			// TODO send error back to client
		}

		// responseBytes, err := handleClient(requestBytes)
		// TODO
		request, err := parseRPCMessage(requestBytes)

		if err != nil {
			fmt.Println("Error: " + err.Error())
		}

		if request.CallBody.Program == portmapv2.Program && request.CallBody.Procedure == portmapv2.ProcedureGetPort {

		} else if request.CallBody.Program == mountv3.Program && request.CallBody.Procedure == mountv3.MountProcedure3Export {
			/*
			 * Response
			 */
			var responseBuffer = new(bytes.Buffer)

			rpcResponse := RPCMsg{
				XID:         request.RPCMessage.XID,
				MessageType: Reply,
			}

			err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse)

			replyBody := ReplyBody{
				ReplyStatus: MessageAccepted,
			}

			err = binary.Write(responseBuffer, binary.BigEndian, &replyBody)

			verifierReply := OpaqueAuth{
				Flavor: AuthenticationNull,
				Length: 0,
			}

			successReply := AcceptedReplySuccess{
				Verifier:    verifierReply,
				AcceptState: Success,
			}

			err = binary.Write(responseBuffer, binary.BigEndian, &successReply)

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

			err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)
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

			// ---- fragments

			var fragmentBuffer = new(bytes.Buffer)
			lastFragment := uint32(1 << 31)
			fragmentLength := uint32(responseBuffer.Len())

			err = binary.Write(fragmentBuffer, binary.BigEndian, lastFragment|fragmentLength)

			// ---- end

			fragmentBuffer.Write(responseBuffer.Bytes())

			_, err = clientConnection.Write(fragmentBuffer.Bytes())

			if err != nil {
				fmt.Println("[mount] Error sending response: ", err.Error())
			}
		} else {
			fmt.Println("Error: unrecognized program ", request.CallBody.Program, " with procedure ", request.CallBody.Procedure)
		}

		if isLastRequest {
			moreRequests = false
		}
	}

	clientConnection.Close()
}

// handleUDPClient handles UDP connections
func HandleUDPClient(requestBytes []byte, serverConnection *net.UDPConn, clientAddress *net.UDPAddr) {
	request, err := parseRPCMessage(requestBytes)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		// TODO send error message back to client
	}

	var mapping portmapv2.Mapping
	var requestBody = bytes.NewBuffer(request.RequestBody)

	err = binary.Read(requestBody, binary.BigEndian, &mapping)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO: send error message back to client
	}

	var responseBuffer = new(bytes.Buffer)

	rpcResponse := RPCMsg{
		XID:         request.RPCMessage.XID,
		MessageType: Reply,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse)

	replyBody := ReplyBody{
		ReplyStatus: MessageAccepted,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &replyBody)

	verifierReply := OpaqueAuth{
		Flavor: AuthenticationNull,
		Length: 0,
	}

	successReply := AcceptedReplySuccess{
		Verifier:    verifierReply,
		AcceptState: Success,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &successReply)

	// TODO check callBody.Program == portmapv2.Program
	// TODO check callBody.Program == portmapv2.Version

	if request.CallBody.Procedure == portmapv2.ProcedureGetPort {
		// TODO check mapping.Version (1) == mountv3.Version (3)
		if mapping.Program == mountv3.Program && mapping.Protocol == portmapv2.IPProtocolTCP {
			var result uint32
			result = 892
			err = binary.Write(responseBuffer, binary.BigEndian, &result)
		}
	}

	serverConnection.WriteToUDP(responseBuffer.Bytes(), clientAddress)
}

func handleRequest(requestBytes []byte) (responseBytes []byte, err error) {
	return responseBytes, err
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

func parseRPCMessage(rpcMessage []byte) (rpcRequest RPCRequest, err error) {
	var requestBuffer = bytes.NewBuffer(rpcMessage)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.RPCMessage)

	if err != nil {
		return rpcRequest, err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcRequest.CallBody)

	if err != nil {
		return rpcRequest, err
	}

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
