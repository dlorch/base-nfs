/*
	NFS Version 3 Protocol Specification (RFC1813)

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

package nfsv3

import "github.com/dlorch/nfsv3/rpcv2"

// NFSService ...
type NFSService struct {
	rpcv2.RPCService
}

// NewNFSv3Service ...
func NewNFSv3Service() *NFSService {
	nfsService := &NFSService{
		RPCService: *rpcv2.NewRPCService("nfsv3", Program, Version),
	}

	nfsService.RegisterProcedure(NFSProcedure3Null, nfsProcedure3Null)
	nfsService.RegisterProcedure(NFSProcedure3GetAttributes, nfsProcedure3GetAttributes)
	nfsService.RegisterProcedure(NFSProcedure3Lookup, nfsProcedure3Lookup)
	nfsService.RegisterProcedure(NFSProcedure3Access, nfsProcedure3Access)
	nfsService.RegisterProcedure(NFSProcedure3FSInfo, nfsProcedure3FSInfo)
	nfsService.RegisterProcedure(NFSProcedure3PathConf, nfsProcedure3PathConf)

	return nfsService
}
