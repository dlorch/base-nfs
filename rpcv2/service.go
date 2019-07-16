package rpcv2

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

/*
func NewRPCService() {

}
*/
