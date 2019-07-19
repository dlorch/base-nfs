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
	"errors"
	"fmt"
	"net"
)

// RPCRequest describes an RPC request
type RPCRequest struct {
	RPCMessage      RPCMsg
	CallBody        CallBody
	Credentials     OpaqueAuth
	CredentialsBody []byte // TODO parse byte array
	Verifier        OpaqueAuth
	VerifierBody    []byte // TODO parse byte array
	RequestBody     []byte
}

// RPCResponse describes an RPC response
type RPCResponse struct {
	RPCMessage           RPCMsg
	ReplyBody            ReplyBody
	AcceptedReplySuccess AcceptedReplySuccess
	ResponseBody         []byte
}

type udpClient struct {
	requestBytes     []byte
	serverConnection *net.UDPConn
	clientAddress    *net.UDPAddr
}

type procedureHandler func(*RPCRequest) *RPCResponse

// RPCService describes and RPC service
type RPCService struct { // TODO lowercase structs to unexport them
	ShortName  string // friendly name (for logging)
	Program    uint32 // RPC program number
	Version    uint32 // RPC program version
	TCPClients chan net.Conn
	UDPClients chan udpClient
	Procedures map[uint32]procedureHandler
}

// NewRPCService returns a new RPC service
func NewRPCService(ShortName string, Program uint32, Version uint32) *RPCService {
	rpcService := &RPCService{
		ShortName:  ShortName,
		Program:    Program,
		Version:    Version,
		TCPClients: make(chan net.Conn),
		UDPClients: make(chan udpClient),
		Procedures: make(map[uint32]procedureHandler),
	}

	return rpcService
}

// AddListener announces the local network address
func (rpcService *RPCService) AddListener(network string, address string) (err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		tcpListener, err := net.Listen(network, address)

		if err != nil {
			return err
		}

		fmt.Printf("[%s] Listening on TCP %s\n", rpcService.ShortName, tcpListener.Addr())

		go func() {
			for {
				clientConnection, err := tcpListener.Accept()

				if err != nil {
					fmt.Printf("[%s] Error: %s\n", rpcService.ShortName, err.Error())
				} else {
					fmt.Printf("[%s] Received TCP request from %s\n", rpcService.ShortName, clientConnection.RemoteAddr())
					rpcService.TCPClients <- clientConnection
				}
			}
		}()
	case "udp", "udp4", "udp6":
		serverAddress, err := net.ResolveUDPAddr(network, address)

		if err != nil {
			return err
		}

		serverConnection, err := net.ListenUDP(network, serverAddress)

		if err != nil {
			return err
		}

		fmt.Printf("[%s] Listening on UDP %s\n", rpcService.ShortName, serverAddress)

		go func() {
			requestBytes := make([]byte, 1024) // TODO: optimal/maximal UDP size?

			for {
				_, clientAddress, err := serverConnection.ReadFromUDP(requestBytes)

				if err != nil {
					fmt.Printf("[%s] Error: %s\n", rpcService.ShortName, err.Error())
				} else {
					fmt.Printf("[%s] Received UDP request from %s\n", rpcService.ShortName, clientAddress)

					udpClient := udpClient{
						requestBytes:     requestBytes,
						serverConnection: serverConnection,
						clientAddress:    clientAddress,
					}

					rpcService.UDPClients <- udpClient
				}
			}
		}()
	default:
		return errors.New("Invalid network provided. Valid options are: tcp, tcp4, tpc6, udp, udp4 or udp6")
	}

	return nil
}

// HandleClients accepts and processes clients
func (rpcService *RPCService) HandleClients() {
	for {
		select {
		case clientConnection := <-rpcService.TCPClients:
			handleTCPClient(clientConnection, rpcService)
		case udpClient := <-rpcService.UDPClients:
			handleUDPClient(udpClient.requestBytes, udpClient.serverConnection, udpClient.clientAddress, rpcService)
		}
	}
}

// RegisterProcedure registers a callback function for a given RPC procedure number
func (rpcService *RPCService) RegisterProcedure(procedure uint32, procedureHandler procedureHandler) {
	rpcService.Procedures[procedure] = procedureHandler
}
