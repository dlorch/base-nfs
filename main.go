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

	fmt.Println(rpcMsg.XID)
	fmt.Println(rpcMsg.MessageType)

	err = binary.Read(requestBuffer, binary.BigEndian, &callBody)

	if err != nil {
		fmt.Println("Error decoding call body: ", err.Error())
	}

	fmt.Println(callBody.RPCVersion)
	fmt.Println(callBody.Program)
	fmt.Println(callBody.ProgramVersion)
	fmt.Println(callBody.Procedure)

	err = binary.Read(requestBuffer, binary.BigEndian, &credentials)

	if err != nil {
		fmt.Println("Error decoding credentials: ", err.Error())
	}

	fmt.Println(credentials.Flavor)
	fmt.Println(credentials.Length)

	err = binary.Read(requestBuffer, binary.BigEndian, &verifier)

	if err != nil {
		fmt.Println("Error decoding verifier: ", err.Error())
	}

	fmt.Println(verifier.Flavor)
	fmt.Println(verifier.Length)

	err = binary.Read(requestBuffer, binary.BigEndian, &mapping)

	if err != nil {
		fmt.Println("Error decoding verifier: ", err.Error())
	}

	fmt.Println(mapping.Program)
	fmt.Println(mapping.Version)
	fmt.Println(mapping.Protocol)
	fmt.Println(mapping.Port)

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

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":111")

	if err != nil {
		fmt.Println("Error resolving UDP address: ", err.Error())
		os.Exit(1)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)

	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Listening at: ", serverAddr)

	defer serverConn.Close()

	requestBytes := make([]byte, 1024)

	for {
		_, clientAddr, err := serverConn.ReadFromUDP(requestBytes)

		if err != nil {
			fmt.Println("Error receiving: ", err.Error())
		} else {
			fmt.Println("Received request from ", clientAddr)
			go handlePortmapRequest(requestBytes, serverConn, clientAddr)
		}
	}
}
