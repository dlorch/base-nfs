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

package rpcv2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// Serializable ...
type Serializable interface {
	ToBytes() ([]byte, error)
}

// RPCMessage describes the request/response RPC header (RFC1057: struct rpc_msg)
type RPCMessage struct {
	XID         uint32
	MessageType uint32
}

// OpaqueAuth describes the type of authentication used (RFC1057: struct opaque_auth)
type OpaqueAuth struct {
	Flavor uint32
	Body   []byte
}

// CallBody describes the body of a CALL request (RFC1057: struct call_body)
type CallBody struct {
	RPCMessage
	RPCVersion     uint32
	Program        uint32
	ProgramVersion uint32
	Procedure      uint32
	Credentials    OpaqueAuth
	Verifier       OpaqueAuth
}

// AcceptedReplySuccess (RFC1057: struct accepted_reply)
type AcceptedReplySuccess struct {
	RPCMessage
	ReplyStatus uint32 // must be MessageAccepted = 0
	Verifier    OpaqueAuth
	AcceptState uint32 // must be Success = 0
}

// AcceptedReplyProgramMismatch (RFC1057: struct accepted_reply)
type AcceptedReplyProgramMismatch struct {
	RPCMessage
	ReplyStatus             uint32 // must be MessageAccepted = 0
	Verifier                OpaqueAuth
	AcceptState             uint32 // must be ProgramMismatch = 2
	LowestVersionSupported  uint32
	HighestVersionSupported uint32
}

// AcceptedReply (RFC1057: struct accepted_reply)
type AcceptedReply struct {
	RPCMessage
	ReplyStatus uint32 // must be MessageAccepted = 0
	Verifier    OpaqueAuth
	AcceptState uint32 // must be ProgramUnavailable = 1, ProcedureUnavailable = 3 or GarbageArguments = 4
}

// RejectedReplyRPCMismatch (RFC1057: rejected_reply)
type RejectedReplyRPCMismatch struct {
	RPCMessage
	ReplyStatus                uint32 // must be MessageDenied = 1
	RejectState                uint32 // must be RPCMismatch = 0
	LowestSupportedRPCVersion  uint32
	HighestSupportedRPCVersion uint32
}

// RejectedReplyAuthenticationError (RFC1057: rejected_reply)
type RejectedReplyAuthenticationError struct {
	RPCMessage
	ReplyStatus         uint32 // must be MessageDenied = 1
	RejectState         uint32 // must be AuthenticationError = 1
	AuthenticationState uint32
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
	RPCVersion              uint32 = 2       // RPC version number
	Call                    uint32 = 0       // Message type CALL
	Reply                   uint32 = 1       // Message type REPLY
	MessageAccepted         uint32 = 0       // Reply to a call message
	MessageDenied           uint32 = 1       // Reply to a call message
	Success                 uint32 = 0       // RPC executed succesfully
	ProgramUnavailable      uint32 = 1       // remote hasn't exported program
	ProgramMismatch         uint32 = 2       // remote can't support version
	ProcedureUnavailable    uint32 = 3       // program can't support procedure
	GarbageArguments        uint32 = 4       // procedure can't decode params
	RPCMismatch             uint32 = 0       // RPC version number != 2
	AuthenticationError     uint32 = 1       // remote can't authenticate caller
	LastFragment            uint32 = 1 << 31 // last fragment delimiter for record marking (RM)
	OpaqueAuthBodyMaxLength uint32 = 400     // maximal length of OpaqueAuth.Body
)

// handleTCPClient handles TCP client connections, reads requests and delimits them into
// individual messages (RFC 1057: 10. Record Marking Standard) for further processing
func handleTCPClient(clientConnection net.Conn, rpcProcedures map[uint32]rpcProcedureHandler) error {
	var responseBytes []byte
	requestBytes, isLastFragment, err := readNextRequestFragment(clientConnection) // TODO not sure this is correct, shouldn't fragments be concatenated?

	for {
		if err != nil {
			if err.Error() == "EOF" { // all good, we reached the end of the TCP request
				return nil
			}
			return err
		}

		responseBytes, err = handleClient(requestBytes, rpcProcedures)

		if err != nil {
			return err
		}

		err = writeResponseFragment(clientConnection, responseBytes, isLastFragment)

		if err != nil {
			return err
		}

		requestBytes, isLastFragment, err = readNextRequestFragment(clientConnection)
	}
}

// handleUDPClient handles UDP connections
func handleUDPClient(requestBytes []byte, serverConnection *net.UDPConn, clientAddress *net.UDPAddr, rpcProcedures map[uint32]rpcProcedureHandler) error {
	responseBytes, err := handleClient(requestBytes, rpcProcedures)

	if err != nil {
		return err
	}

	serverConnection.WriteToUDP(responseBytes, clientAddress)

	return nil
}

func handleClient(requestBytes []byte, rpcProcedures map[uint32]rpcProcedureHandler) (responseBytes []byte, err error) {
	rpcRequest, argumentsIndex, err := parseRPCCallBody(requestBytes)

	if err != nil {
		return responseBytes, errors.New("Malformed RPC request")
	}

	var rpcResponse Serializable

	if rpcRequest.RPCVersion != RPCVersion {
		rpcResponse = &RejectedReplyRPCMismatch{
			RPCMessage: RPCMessage{
				XID:         rpcRequest.RPCMessage.XID,
				MessageType: Reply,
			},
			ReplyStatus:                MessageDenied,
			RejectState:                RPCMismatch,
			LowestSupportedRPCVersion:  RPCVersion,
			HighestSupportedRPCVersion: RPCVersion,
		}

		return rpcResponse.ToBytes()
	}

	rpcProcedure, found := rpcProcedures[rpcRequest.Procedure]

	// TODO how to check for ProgramMismatch?
	fmt.Println("Procedure ", rpcRequest.Procedure, " for program ", rpcRequest.Program)

	if !found {
		rpcResponse = &AcceptedReply{
			RPCMessage: RPCMessage{
				XID:         rpcRequest.RPCMessage.XID,
				MessageType: Reply,
			},
			ReplyStatus: MessageAccepted,
			Verifier: OpaqueAuth{
				Flavor: AuthenticationNull,
				Body:   []byte{},
			},
			AcceptState: ProcedureUnavailable,
		}

		return rpcResponse.ToBytes()
	}

	var procedureResponse Serializable

	procedureResponse, err = rpcProcedure(requestBytes[argumentsIndex:])

	if err != nil {
		fmt.Println("Error: ", err.Error())

		rpcResponse = &AcceptedReply{
			RPCMessage: RPCMessage{
				XID:         rpcRequest.RPCMessage.XID,
				MessageType: Reply,
			},
			ReplyStatus: MessageAccepted,
			Verifier: OpaqueAuth{
				Flavor: AuthenticationNull,
				Body:   []byte{},
			},
			AcceptState: GarbageArguments,
		}

		return rpcResponse.ToBytes()
	}

	rpcResponse = &AcceptedReplySuccess{
		RPCMessage: RPCMessage{
			XID:         rpcRequest.RPCMessage.XID,
			MessageType: Reply,
		},
		ReplyStatus: MessageAccepted,
		Verifier: OpaqueAuth{
			Flavor: AuthenticationNull,
			Body:   []byte{},
		},
		AcceptState: Success,
	}

	var response []byte
	responseBuffer := new(bytes.Buffer)

	response, err = rpcResponse.ToBytes()

	if err != nil {
		return responseBytes, err
	}

	_, err = responseBuffer.Write(response)

	if err != nil {
		return responseBytes, err
	}

	response, err = procedureResponse.ToBytes()

	if err != nil {
		return responseBytes, err
	}

	_, err = responseBuffer.Write(response)

	if err != nil {
		return responseBytes, err
	}

	response = make([]byte, responseBuffer.Len())
	copy(response, responseBuffer.Bytes())

	return response, nil
}

func readNextRequestFragment(clientConnection net.Conn) (requestBytes []byte, isLastFragment bool, err error) {
	var fragmentHeader uint32
	fragmentHeaderBytes := make([]byte, 4)

	_, err = clientConnection.Read(fragmentHeaderBytes)

	if err != nil {
		return requestBytes, isLastFragment, err
	}

	fragmentHeaderBytesBuffer := bytes.NewBuffer(fragmentHeaderBytes)

	err = binary.Read(fragmentHeaderBytesBuffer, binary.BigEndian, &fragmentHeader)

	if err != nil {
		return requestBytes, isLastFragment, err
	}

	isLastFragment = (fragmentHeader & LastFragment) != 0
	remainingFragmentLength := uint32(fragmentHeader & ^LastFragment)

	messageBytesBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, 1024)

	for remainingFragmentLength > 0 { // TODO rewrite this
		bytesRead, err := clientConnection.Read(readBuffer)

		if err != nil {
			return requestBytes, isLastFragment, err
		}

		_, err = messageBytesBuffer.Write(readBuffer)

		if err != nil {
			return requestBytes, isLastFragment, err
		}

		remainingFragmentLength -= uint32(bytesRead)
	}

	requestBytes = make([]byte, messageBytesBuffer.Len())
	copy(requestBytes, messageBytesBuffer.Bytes())

	return requestBytes, isLastFragment, nil
}

func writeResponseFragment(clientConnection net.Conn, responseBytes []byte, lastFragment bool) error {
	var err error
	fragmentBuffer := new(bytes.Buffer)
	fragmentLength := uint32(len(responseBytes))

	if lastFragment {
		err = binary.Write(fragmentBuffer, binary.BigEndian, LastFragment|fragmentLength)
	} else {
		err = binary.Write(fragmentBuffer, binary.BigEndian, fragmentLength)
	}

	if err != nil {
		return err
	}

	_, err = fragmentBuffer.Write(responseBytes)

	if err != nil {
		return err
	}

	_, err = clientConnection.Write(fragmentBuffer.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func parseRPCCallBody(requestBytes []byte) (rpcCallBody CallBody, bytesRead int, err error) {
	requestBuffer := bytes.NewBuffer(requestBytes)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.RPCMessage)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.RPCVersion)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.Program)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.ProgramVersion)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.Procedure)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.Credentials.Flavor)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	var credentialsLength uint32

	err = binary.Read(requestBuffer, binary.BigEndian, &credentialsLength)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	if credentialsLength > OpaqueAuthBodyMaxLength {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(),
			fmt.Errorf("Invalid length '%d' for Credentials in CallBody. Maximum value of '%d' allowed", credentialsLength, OpaqueAuthBodyMaxLength)
	}

	rpcCallBody.Credentials.Body = make([]byte, credentialsLength)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.Credentials.Body)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.Verifier.Flavor)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	var verifierLength uint32

	err = binary.Read(requestBuffer, binary.BigEndian, &verifierLength)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	if verifierLength > OpaqueAuthBodyMaxLength {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(),
			fmt.Errorf("Invalid length '%d' for Verifier in CallBody. Maximum value of '%d' allowed", verifierLength, OpaqueAuthBodyMaxLength)
	}

	rpcCallBody.Verifier.Body = make([]byte, verifierLength)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.Verifier.Body)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	return rpcCallBody, len(requestBytes) - requestBuffer.Len(), nil
}

// ToBytes serializes the AcceptedReplySuccess to be sent back to the client
func (reply *AcceptedReplySuccess) ToBytes() ([]byte, error) {
	responseBuffer := new(bytes.Buffer)
	responseBytes := []byte{}

	err := binary.Write(responseBuffer, binary.BigEndian, &reply.RPCMessage)

	if err != nil {
		return responseBytes, err
	}

	if reply.ReplyStatus != MessageAccepted {
		return responseBytes, fmt.Errorf("Invalid ReplyStatus '%d' in AcceptedReplySuccess. Expecting '%d'", reply.ReplyStatus, MessageAccepted)
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.ReplyStatus)

	if err != nil {
		return responseBytes, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.Verifier.Flavor)

	if err != nil {
		return responseBytes, err
	}

	verifierLength := uint32(len(reply.Verifier.Body))

	err = binary.Write(responseBuffer, binary.BigEndian, &verifierLength)

	if err != nil {
		return responseBytes, err
	}

	_, err = responseBuffer.Write(reply.Verifier.Body)

	if err != nil {
		return responseBytes, err
	}

	if reply.AcceptState != Success {
		return responseBytes, fmt.Errorf("Invalid AcceptState '%d' in AcceptedReplySuccess. Expecting '%d'", reply.AcceptState, Success)
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.AcceptState)

	if err != nil {
		return responseBytes, err
	}

	responseBytes = make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())

	return responseBytes, nil
}

/*
// ToBytes ..
func (reply *AcceptedReplyProgramMismatch) ToBytes() ([]byte, error) {

}
*/

// ToBytes serializes the AcceptedReply to be sent back to the client
func (reply *AcceptedReply) ToBytes() ([]byte, error) {
	responseBuffer := new(bytes.Buffer)
	responseBytes := []byte{}

	err := binary.Write(responseBuffer, binary.BigEndian, &reply.RPCMessage)

	if err != nil {
		return responseBytes, err
	}

	if reply.ReplyStatus != MessageAccepted {
		return responseBytes, fmt.Errorf("Invalid ReplyStatus '%d' in AcceptedReplySuccess. Expecting '%d'", reply.ReplyStatus, MessageAccepted)
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.ReplyStatus)

	if err != nil {
		return responseBytes, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.Verifier.Flavor)

	if err != nil {
		return responseBytes, err
	}

	verifierLength := uint32(len(reply.Verifier.Body))

	err = binary.Write(responseBuffer, binary.BigEndian, &verifierLength)

	if err != nil {
		return responseBytes, err
	}

	_, err = responseBuffer.Write(reply.Verifier.Body)

	if err != nil {
		return responseBytes, err
	}

	if reply.AcceptState != ProgramUnavailable && reply.AcceptState != ProcedureUnavailable && reply.AcceptState != GarbageArguments {
		return responseBytes, fmt.Errorf("Invalid AcceptState '%d' in AcceptedReplySuccess. Expecting '%d'", reply.AcceptState, Success)
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.AcceptState)

	if err != nil {
		return responseBytes, err
	}

	responseBytes = make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())

	return responseBytes, nil
}

// ToBytes serializes the RejectedReplyRPCMismatch to be sent back to the client
func (reply *RejectedReplyRPCMismatch) ToBytes() ([]byte, error) {
	responseBuffer := new(bytes.Buffer)
	responseBytes := []byte{}

	err := binary.Write(responseBuffer, binary.BigEndian, &reply.RPCMessage)

	if err != nil {
		return responseBytes, err
	}

	if reply.ReplyStatus != MessageDenied {
		return responseBytes, fmt.Errorf("Invalid ReplyStatus '%d' in RejectedReplyRPCMismatch. Expecting '%d'", reply.RejectState, MessageDenied)
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.ReplyStatus)

	if err != nil {
		return responseBytes, err
	}

	if reply.RejectState != RPCMismatch {
		return responseBytes, fmt.Errorf("Invalid RejectState '%d' in RejectedReplyRPCMismatch. Expecting '%d'", reply.RejectState, RPCMismatch)
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.LowestSupportedRPCVersion)

	if err != nil {
		return responseBytes, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.HighestSupportedRPCVersion)

	if err != nil {
		return responseBytes, err
	}

	responseBytes = make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())

	return responseBytes, nil
}

/*
// ToBytes ..
func (reply *RejectedReplyAuthenticationError) ToBytes() ([]byte, error) {

}
*/
