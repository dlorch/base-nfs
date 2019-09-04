// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpcv2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/dlorch/base-nfs/xdr"
)

// RPCMessage describes the request/response RPC header (RFC1057: struct rpc_msg)
type RPCMessage struct {
	XID         uint32
	MessageType uint32    `xdr:"switch"`
	CBody       CallBody  `xdr:"case=0"`
	RBody       ReplyBody `xdr:"case=1"`
}

// CallBody describes the body of a CALL request (RFC1057: struct call_body)
type CallBody struct {
	RPCVersion     uint32
	Program        uint32
	ProgramVersion uint32
	Procedure      uint32
	Credentials    OpaqueAuth
	Verifier       OpaqueAuth
}

// ReplyBody (RFC1057: union reply_body)
type ReplyBody struct {
	ReplyStatus uint32        `xdr:"switch"`
	AReply      AcceptedReply `xdr:"case=0"`
	RReply      RejectedReply `xdr:"case=1"`
}

// AcceptedReply (RFC1057: struct accepted_reply)
type AcceptedReply struct {
	Verf         OpaqueAuth
	AcceptState  uint32       `xdr:"switch"`
	Results      interface{}  `xdr:"case=0"`
	MismatchInfo MismatchInfo `xdr:"case=2"`
}

// RejectedReply (RFC1057: struct rejected_reply)
type RejectedReply struct {
	RejectState  uint32       `xdr:"switch"`
	MismatchInfo MismatchInfo `xdr:"case=0"`
	Stat         uint32       `xdr:"case=1"`
}

// MismatchInfo indicates the lowest and highest versions supported
type MismatchInfo struct {
	Low  uint32 // lowest version supported
	High uint32 // highest version supported
}

// OpaqueAuth describes the type of authentication used (RFC1057: struct opaque_auth)
type OpaqueAuth struct {
	Flavor uint32
	Body   []byte
}

// Void is a void reply
type Void struct{}

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

	if rpcRequest.CBody.RPCVersion != RPCVersion {
		rpcMismatch := &RPCMessage{
			XID:         rpcRequest.XID,
			MessageType: Reply,
			RBody: ReplyBody{
				ReplyStatus: MessageDenied,
				RReply: RejectedReply{
					RejectState: RPCMismatch,
					MismatchInfo: MismatchInfo{
						Low:  RPCVersion,
						High: RPCVersion,
					},
				},
			},
		}

		return xdr.Marshal(rpcMismatch)
	}

	rpcProcedure, found := rpcProcedures[rpcRequest.CBody.Procedure]

	// TODO how to check for ProgramMismatch?
	fmt.Println("Procedure ", rpcRequest.CBody.Procedure, " for program ", rpcRequest.CBody.Program)

	if !found {
		procUnavail := &RPCMessage{
			XID:         rpcRequest.XID,
			MessageType: Reply,
			RBody: ReplyBody{
				ReplyStatus: MessageAccepted,
				AReply: AcceptedReply{
					Verf: OpaqueAuth{
						Flavor: AuthenticationNull,
						Body:   []byte{},
					},
					AcceptState: ProcedureUnavailable,
				},
			},
		}

		return xdr.Marshal(procUnavail)
	}

	procedureResponse, err := rpcProcedure(requestBytes[argumentsIndex:])

	if err != nil {
		fmt.Println("Error: ", err.Error())

		garbageArgs := &RPCMessage{
			XID:         rpcRequest.XID,
			MessageType: Reply,
			RBody: ReplyBody{
				ReplyStatus: MessageAccepted,
				AReply: AcceptedReply{
					Verf: OpaqueAuth{
						Flavor: AuthenticationNull,
						Body:   []byte{},
					},
					AcceptState: GarbageArguments,
				},
			},
		}

		return xdr.Marshal(garbageArgs)
	}

	acceptedReply := &RPCMessage{
		XID:         rpcRequest.XID,
		MessageType: Reply,
		RBody: ReplyBody{
			ReplyStatus: MessageAccepted,
			AReply: AcceptedReply{
				Verf: OpaqueAuth{
					Flavor: AuthenticationNull,
					Body:   []byte{},
				},
				AcceptState: Success,
				Results:     procedureResponse,
			},
		},
	}

	return xdr.Marshal(acceptedReply)
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

func parseRPCCallBody(requestBytes []byte) (rpcCallBody RPCMessage, bytesRead int, err error) {
	requestBuffer := bytes.NewBuffer(requestBytes)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.XID)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.MessageType)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.RPCVersion)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.Program)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.ProgramVersion)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.Procedure)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.Credentials.Flavor)

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

	rpcCallBody.CBody.Credentials.Body = make([]byte, credentialsLength)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.Credentials.Body)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.Verifier.Flavor)

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

	rpcCallBody.CBody.Verifier.Body = make([]byte, verifierLength)

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcCallBody.CBody.Verifier.Body)

	if err != nil {
		return rpcCallBody, len(requestBytes) - requestBuffer.Len(), err
	}

	return rpcCallBody, len(requestBytes) - requestBuffer.Len(), nil
}
