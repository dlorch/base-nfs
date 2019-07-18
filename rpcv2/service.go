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
}

// RPCService describes and RPC service
type RPCService struct { // TODO lowercase structs to unexport them
	ShortName     string
	Program       uint32
	Version       uint32
	Network       string // "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6"
	TCPListener   net.Listener
	UDPConnection *net.UDPConn
}

// NewRPCService returns a new RPC service
func NewRPCService(ShortName string, Program uint32, Version uint32) *RPCService {
	rpcService := &RPCService{
		ShortName: ShortName,
		Program:   Program,
		Version:   Version,
	}

	return rpcService
}

// Listen announces the local network address
func (rpcService *RPCService) Listen(network string, address string) (err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		rpcService.TCPListener, err = net.Listen(network, address)

		if err != nil {
			return err
		}

		fmt.Printf("[%s] Listening on TCP %s\n", rpcService.ShortName, rpcService.TCPListener.Addr())

		rpcService.Network = network
	case "udp", "udp4", "udp6":
		serverAddress, err := net.ResolveUDPAddr(network, address)

		if err != nil {
			return err
		}

		rpcService.UDPConnection, err = net.ListenUDP(network, serverAddress)

		if err != nil {
			return err
		}

		fmt.Printf("[%s] Listening on UDP %s\n", rpcService.ShortName, serverAddress)

		rpcService.Network = network
	default:
		return errors.New("Invalid network provided. Valid options: tcp, tcp4, tpc6, udp, udp4 or udp6")
	}

	return nil
}

// HandleClients accepts and processes clients
func (rpcService *RPCService) HandleClients() {
	switch rpcService.Network {
	case "tcp", "tcp4", "tcp6":
		defer rpcService.TCPListener.Close()

		for {
			clientConnection, err := rpcService.TCPListener.Accept()

			if err != nil {
				fmt.Printf("[%s] Error: %s\n", rpcService.ShortName, err.Error())
			} else {
				fmt.Printf("[%s] Received TCP request from %s\n", rpcService.ShortName, clientConnection.RemoteAddr())
				HandleTCPClient(clientConnection)
			}
		}
	case "udp", "udp4", "udp6":
		defer rpcService.UDPConnection.Close()
		requestBytes := make([]byte, 1024) // TODO: optimal/maximal UDP size?

		for {
			_, clientAddress, err := rpcService.UDPConnection.ReadFromUDP(requestBytes)

			if err != nil {
				fmt.Printf("[%s] Error: %s\n", rpcService.ShortName, err.Error())
			} else {
				fmt.Printf("[%s] Received UDP request from %s\n", rpcService.ShortName, clientAddress)
				HandleUDPClient(requestBytes, rpcService.UDPConnection, clientAddress)
			}
		}
	}

	fmt.Printf("[%s] Error: service not listening\n", rpcService.ShortName)
}
