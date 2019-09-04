// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nfsv3

import "github.com/dlorch/base-nfs/rpcv2"

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
	nfsService.RegisterProcedure(NFSProcedure3Lookup, Lookup3)
	nfsService.RegisterProcedure(NFSProcedure3Access, nfsProcedure3Access)
	nfsService.RegisterProcedure(NFSProcedure3FSInfo, nfsProcedure3FSInfo)
	nfsService.RegisterProcedure(NFSProcedure3PathConf, nfsProcedure3PathConf)
	nfsService.RegisterProcedure(NFSProcedure3ReadDirPlus, nfsProcedure3ReadDirPlus)

	return nfsService
}
