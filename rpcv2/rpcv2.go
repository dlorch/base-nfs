package rpcv2

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
	Call                 uint32 = 0 // Message type CALL
	Reply                uint32 = 1 // Message type REPLY
	MessageAccepted      uint32 = 0 // Reply to a call message
	MessageDenied        uint32 = 1 // Reply to a call message
	Success              uint32 = 0 // RPC executed succesfully
	ProgramUnavailable   uint32 = 1 // remote hasn't exported program
	ProgramMismatch      uint32 = 2 // remote can't support version
	ProcedureUnavailable uint32 = 3 // program can't support procedure
	GarbageArguments     uint32 = 4 // procedure can't decode params
)
