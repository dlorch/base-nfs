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
	"sync"
)

type udpClient struct {
	requestBytes     []byte
	serverConnection *net.UDPConn
	clientAddress    *net.UDPAddr
}

type rpcProcedureHandler func([]byte) (Serializable, error)

// RPCService represents an RPC service
type RPCService struct {
	shortName    string // friendly name (for logging)
	program      uint32 // RPC program number
	version      uint32 // RPC program version
	tcpClients   chan net.Conn
	tcpListeners []net.Listener
	udpClients   chan udpClient
	udpListeners []*net.UDPConn
	procedures   map[uint32]rpcProcedureHandler
	listening    bool
	waitGroup    sync.WaitGroup
}

// NewRPCService returns a new RPC service
func NewRPCService(shortName string, program uint32, version uint32) *RPCService {
	rpcService := &RPCService{
		shortName:  shortName,
		program:    program,
		version:    version,
		tcpClients: make(chan net.Conn),
		udpClients: make(chan udpClient),
		procedures: make(map[uint32]rpcProcedureHandler),
		listening:  false,
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

		fmt.Printf("[%s] Listening on TCP %s\n", rpcService.shortName, tcpListener.Addr())

		rpcService.listening = true
		rpcService.tcpListeners = append(rpcService.tcpListeners, tcpListener)
		rpcService.waitGroup.Add(1)

		go func() {
			for rpcService.listening {
				clientConnection, err := tcpListener.Accept()

				if rpcService.listening { // closing the tcpListener in RemoveAllListeners() will cause an accept error - ignore
					if err != nil {
						fmt.Printf("[%s] Error: %s\n", rpcService.shortName, err.Error())
					} else {
						fmt.Printf("[%s] Received TCP request from %s\n", rpcService.shortName, clientConnection.RemoteAddr())
						rpcService.tcpClients <- clientConnection
					}
				}
			}

			rpcService.waitGroup.Done()
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

		fmt.Printf("[%s] Listening on UDP %s\n", rpcService.shortName, serverAddress)

		rpcService.listening = true
		rpcService.udpListeners = append(rpcService.udpListeners, serverConnection)
		rpcService.waitGroup.Add(1)

		go func() {
			// What is the largest UDP datagram for an RPC request? RFC1057 does not mention a maximum size.
			// The Linux NFS page https://wiki.linux-nfs.org/wiki/index.php?title=NetworkTracing&oldid=2945#RPC_over_UDP_datagrams
			// talks about a maximum size of 64 KB, without giving any sources. This value probably originates
			// from the maximum size of an UDP datagram (not considering the header sizes). Although such a big
			// UDP datagram is unlikely to arrive over the Internet due to fragmentation, let's a assume a
			// safe upper limit.
			b := make([]byte, 65536)

			for rpcService.listening {
				n, clientAddress, err := serverConnection.ReadFromUDP(b)

				if rpcService.listening { // closing the udpListener in RemoveAllListeners() will cause a read error - ignore
					if err != nil {
						fmt.Printf("[%s] Error: %s\n", rpcService.shortName, err.Error())
					} else {
						fmt.Printf("[%s] Received UDP request from %s\n", rpcService.shortName, clientAddress)

						requestBytes := make([]byte, n)
						copy(requestBytes, b)

						udpClient := udpClient{
							requestBytes:     requestBytes,
							serverConnection: serverConnection,
							clientAddress:    clientAddress,
						}

						rpcService.udpClients <- udpClient
					}
				}
			}

			rpcService.waitGroup.Done()
		}()
	default:
		return errors.New("Invalid network provided. Valid options are: tcp, tcp4, tpc6, udp, udp4 or udp6")
	}

	return nil
}

// HandleClients accepts and processes clients
func (rpcService *RPCService) HandleClients() {
	var err error

	for {
		select {
		case clientConnection := <-rpcService.tcpClients:
			err = handleTCPClient(clientConnection, rpcService.procedures)
		case udpClient := <-rpcService.udpClients:
			err = handleUDPClient(udpClient.requestBytes, udpClient.serverConnection, udpClient.clientAddress, rpcService.procedures)
		}

		if err != nil {
			fmt.Printf("[%s] Error: %s\n", rpcService.shortName, err.Error())
		}
	}
}

// RegisterProcedure registers a callback function for a given RPC procedure number
func (rpcService *RPCService) RegisterProcedure(procedure uint32, rpcProcedureHandler rpcProcedureHandler) {
	rpcService.procedures[procedure] = rpcProcedureHandler
}

// RemoveAllListeners stops all UDP and TCP listeners, and removes them
func (rpcService *RPCService) RemoveAllListeners() {
	rpcService.listening = false

	for _, tcpListener := range rpcService.tcpListeners {
		tcpListener.Close()
	}
	rpcService.tcpListeners = make([]net.Listener, 0)

	for _, udpListener := range rpcService.udpListeners {
		udpListener.Close()
	}
	rpcService.udpListeners = make([]*net.UDPConn, 0)
}

// WaitUntilDone is a blocking call that waits until all listeners are stopped
func (rpcService *RPCService) WaitUntilDone() {
	rpcService.waitGroup.Wait()
}
