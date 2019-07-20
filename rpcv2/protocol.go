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

// AcceptedReply (RFC1057: struct accepted_reply)
type AcceptedReply struct {
	Verifier                OpaqueAuth
	AcceptState             uint32
	Results                 []byte
	LowestVersionSupported  uint32
	HighestVersionSupported uint32
}

// RejectedReply (RFC1057: rjected_reply)
type RejectedReply struct {
	RejectState                uint32
	LowestSupportedRPCVersion  uint32
	HighestSupportedRPCVersion uint32
	// TODO auth_stat stat
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
	RPCVersion           uint32 = 2       // RPC version number
	Call                 uint32 = 0       // Message type CALL
	Reply                uint32 = 1       // Message type REPLY
	MessageAccepted      uint32 = 0       // Reply to a call message
	MessageDenied        uint32 = 1       // Reply to a call message
	Success              uint32 = 0       // RPC executed succesfully
	ProgramUnavailable   uint32 = 1       // remote hasn't exported program
	ProgramMismatch      uint32 = 2       // remote can't support version
	ProcedureUnavailable uint32 = 3       // program can't support procedure
	GarbageArguments     uint32 = 4       // procedure can't decode params
	RPCMismatch          uint32 = 0       // RPC version number != 2
	AuthenticationError  uint32 = 1       // remote can't authenticate caller
	LastFragment         uint32 = 1 << 31 // last fragment delimiter for record marking (RM)
)

// handleTCPClient handles TCP client connections, reads requests and delimits them into
// individual messages (RFC 1057: 10. Record Marking Standard) for further processing
func handleTCPClient(clientConnection net.Conn, rpcProcedures map[uint32]procedureHandler) error {
	var responseBytes []byte
	requestBytes, isLastFragment, err := readNextRequestFragment(clientConnection)

	for {
		if err != nil {
			if err.Error() == "EOF" {
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
func handleUDPClient(requestBytes []byte, serverConnection *net.UDPConn, clientAddress *net.UDPAddr, rpcProcedures map[uint32]procedureHandler) error {
	responseBytes, err := handleClient(requestBytes, rpcProcedures)

	if err != nil {
		return err
	}

	serverConnection.WriteToUDP(responseBytes, clientAddress)

	return nil
}

func handleClient(requestBytes []byte, rpcProcedures map[uint32]procedureHandler) (responseBytes []byte, err error) {
	rpcRequest, err := parseRPCRequest(requestBytes)

	if err != nil {
		return responseBytes, errors.New("Malformed RPC request")
	}

	var rpcResponse *RPCResponse

	if rpcRequest.CallBody.RPCVersion != RPCVersion {
		rpcResponse = &RPCResponse{
			RPCMessage: RPCMsg{
				XID:         rpcRequest.RPCMessage.XID,
				MessageType: Reply,
			},
			ReplyBody: ReplyBody{
				ReplyStatus: MessageDenied,
			},
			RejectedReply: RejectedReply{
				RejectState:                RPCMismatch,
				LowestSupportedRPCVersion:  RPCVersion,
				HighestSupportedRPCVersion: RPCVersion,
			},
		}
	} else {
		rpcProcedure, found := rpcProcedures[rpcRequest.CallBody.Procedure]

		fmt.Println("Procedure ", rpcRequest.CallBody.Procedure, " for program ", rpcRequest.CallBody.Program)

		if !found {
			rpcResponse = &RPCResponse{
				RPCMessage: RPCMsg{
					XID:         rpcRequest.RPCMessage.XID,
					MessageType: Reply,
				},
				ReplyBody: ReplyBody{
					ReplyStatus: MessageAccepted,
				},
				AcceptedReply: AcceptedReply{
					Verifier: OpaqueAuth{
						Flavor: AuthenticationNull,
						Length: 0,
					},
					AcceptState: ProcedureUnavailable,
				},
			}
		} else {
			rpcResponse = rpcProcedure(&rpcRequest)
		}
	}

	responseBytes, err = serializeRPCResponse(rpcResponse)

	return responseBytes, err
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

func parseRPCRequest(requestBytes []byte) (rpcRequest RPCRequest, err error) {
	requestBuffer := bytes.NewBuffer(requestBytes)

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

func serializeRPCResponse(rpcResponse *RPCResponse) (responseBytes []byte, err error) {
	responseBuffer := new(bytes.Buffer)

	err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.RPCMessage)

	if err != nil {
		return responseBytes, err
	}

	switch rpcResponse.ReplyBody.ReplyStatus {
	case MessageAccepted:
		err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.ReplyBody)

		if err != nil {
			return responseBytes, err
		}

		err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.AcceptedReply.Verifier)

		if err != nil {
			return responseBytes, err
		}

		err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.AcceptedReply.AcceptState)

		if err != nil {
			return responseBytes, err
		}

		switch rpcResponse.AcceptedReply.AcceptState {
		case Success:
			_, err = responseBuffer.Write(rpcResponse.AcceptedReply.Results)

			if err != nil {
				return responseBytes, err
			}
		case ProgramMismatch:
			err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.AcceptedReply.LowestVersionSupported)

			if err != nil {
				return responseBytes, err
			}

			err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse.AcceptedReply.HighestVersionSupported)

			if err != nil {
				return responseBytes, err
			}
		case ProgramUnavailable, ProcedureUnavailable, GarbageArguments:
			// void (intentionally left blank)
		default:
			return responseBytes, fmt.Errorf("Unrecognized AcceptState value '%d' in AcceptedReply", rpcResponse.AcceptedReply.AcceptState)
		}
	case MessageDenied:
		// TODO
		return responseBytes, fmt.Errorf("Unimplemented ReplyStatus value '%d' in ReplyBody", rpcResponse.ReplyBody.ReplyStatus)
	default:
		return responseBytes, fmt.Errorf("Invalid ReplyStatus value '%d' in ReplyBody", rpcResponse.ReplyBody.ReplyStatus)
	}

	responseBytes = make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())

	return responseBytes, nil
}
