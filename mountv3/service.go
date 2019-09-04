// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mountv3

import (
	"github.com/dlorch/base-nfs/rpcv2"
)

// MountService ...
type MountService struct {
	rpcv2.RPCService
}

// NewMountService ...
func NewMountService() *MountService {
	mountService := &MountService{
		RPCService: *rpcv2.NewRPCService("mount", Program, Version),
	}

	mountService.RegisterProcedure(MountProcedure3Null, mountProcedure3Null)
	mountService.RegisterProcedure(MountProcedure3Export, Export)
	mountService.RegisterProcedure(MountProcedure3Mnt, Mnt)

	return mountService
}
