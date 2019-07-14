package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/dlorch/nfsv3/portmapv2"
	"github.com/dlorch/nfsv3/rpcv2"
)

/*
	NFSv3 Server

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

func handlePortmapRequest(requestBytes []byte, serverConn *net.UDPConn, clientAddr *net.UDPAddr) {
	var requestBuffer = bytes.NewBuffer(requestBytes)
	var rpcMsg rpcv2.RPCMsg
	var callBody rpcv2.CallBody
	var credentials rpcv2.OpaqueAuth
	var verifier rpcv2.OpaqueAuth
	var mapping portmapv2.Mapping

	/*
	 * Request
	 */

	err := binary.Read(requestBuffer, binary.BigEndian, &rpcMsg) // network byte order is big endian

	if err != nil {
		fmt.Println("Error decoding RPC message: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &callBody)

	if err != nil {
		fmt.Println("Error decoding call body: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &credentials)

	if err != nil {
		fmt.Println("Error decoding credentials: ", err.Error())
	}

	credentialsBody := make([]byte, credentials.Length)
	err = binary.Read(requestBuffer, binary.BigEndian, &credentialsBody)

	if err != nil {
		fmt.Println("[mount] Error decoding credentials body: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &verifier)

	if err != nil {
		fmt.Println("Error decoding verifier: ", err.Error())
	}

	verifierBody := make([]byte, verifier.Length)
	err = binary.Read(requestBuffer, binary.BigEndian, &verifierBody)

	if err != nil {
		fmt.Println("[mount] Error decoding verifier body: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &mapping)

	if err != nil {
		fmt.Println("Error decoding mapping: ", err.Error())
	}

	/*
	 * Response
	 */

	var responseBuffer = new(bytes.Buffer)

	rpcResponse := rpcv2.RPCMsg{
		XID:         rpcMsg.XID,
		MessageType: rpcv2.Reply,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse)

	replyBody := rpcv2.ReplyBody{
		ReplyStatus: rpcv2.MessageAccepted,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &replyBody)

	verifierReply := rpcv2.OpaqueAuth{
		Flavor: rpcv2.AuthenticationNull,
		Length: 0,
	}

	successReply := rpcv2.AcceptedReplySuccess{
		Verifier:    verifierReply,
		AcceptState: rpcv2.Success,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &successReply)

	var result uint32
	result = 892

	err = binary.Write(responseBuffer, binary.BigEndian, &result)

	serverConn.WriteToUDP(responseBuffer.Bytes(), clientAddr)
}

func runPortmapperService() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":111")

	if err != nil {
		fmt.Println("[portmapper] Error resolving UDP address: ", err.Error())
		os.Exit(1)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)

	if err != nil {
		fmt.Println("[portmap] Error listening: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("[portmap] Listening at: ", serverAddr)

	defer serverConn.Close()

	requestBytes := make([]byte, 1024)

	for {
		_, clientAddr, err := serverConn.ReadFromUDP(requestBytes)

		if err != nil {
			fmt.Println("[portmap] Error receiving: ", err.Error())
		} else {
			fmt.Println("[portmap] Received request from ", clientAddr)
			go handlePortmapRequest(requestBytes, serverConn, clientAddr)
		}
	}
}

func handleMountRequest(clientConnection net.Conn) {
	requestBytes := make([]byte, 1024)

	_, err := clientConnection.Read(requestBytes)

	if err != nil {
		fmt.Println("[mount] Error reading: ", err.Error())
	}

	/*
	 * Request
	 */
	var requestBuffer = bytes.NewBuffer(requestBytes)
	var fragmentHeader uint32
	var rpcMsg rpcv2.RPCMsg
	var callBody rpcv2.CallBody
	var credentials rpcv2.OpaqueAuth
	var verifier rpcv2.OpaqueAuth

	err = binary.Read(requestBuffer, binary.BigEndian, &fragmentHeader)

	if err != nil {
		fmt.Println("[mount] Error reading fragment header: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &rpcMsg)

	if err != nil {
		fmt.Println("[mount] Error decoding RPC message: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &callBody)

	if err != nil {
		fmt.Println("[mount] Error decoding call body: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &credentials)

	if err != nil {
		fmt.Println("[mount] Error decoding credentials: ", err.Error())
	}

	credentialsBody := make([]byte, credentials.Length)
	err = binary.Read(requestBuffer, binary.BigEndian, &credentialsBody)

	if err != nil {
		fmt.Println("[mount] Error decoding credentials body: ", err.Error())
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &verifier)

	if err != nil {
		fmt.Println("[mount] Error decoding verifier: ", err.Error())
	}

	verifierBody := make([]byte, verifier.Length)
	err = binary.Read(requestBuffer, binary.BigEndian, &verifierBody)

	if err != nil {
		fmt.Println("[mount] Error decoding verifier body: ", err.Error())
	}

	/*
	 * Response
	 */
	var responseBuffer = new(bytes.Buffer)

	rpcResponse := rpcv2.RPCMsg{
		XID:         rpcMsg.XID,
		MessageType: rpcv2.Reply,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &rpcResponse)

	replyBody := rpcv2.ReplyBody{
		ReplyStatus: rpcv2.MessageAccepted,
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &replyBody)

	verifierReply := rpcv2.OpaqueAuth{
		Flavor: rpcv2.AuthenticationNull,
		Length: 0,
	}

	successReply := rpcv2.AcceptedReplySuccess{
		Verifier:    verifierReply,
		AcceptState: rpcv2.Success,
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

	clientConnection.Close()
}

func runMountService() {
	serverConn, err := net.Listen("tcp", ":892")

	if err != nil {
		fmt.Println("[mount] Error listening: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("[mount] Listening at: ", serverConn.Addr())

	defer serverConn.Close()

	for {
		clientConnection, err := serverConn.Accept()

		if err != nil {
			fmt.Println("[mount] Error receiving: ", err.Error())
		} else {
			fmt.Println("[mount] Received request from ", clientConnection.RemoteAddr())
			go handleMountRequest(clientConnection)
		}
	}
}

func main() {
	go runPortmapperService()
	runMountService()
}
