// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package portmapv2

import (
	"fmt"

	"github.com/dlorch/base-nfs/mountv3"
	"github.com/dlorch/base-nfs/nfsv3"
	"github.com/dlorch/base-nfs/rpcv2"
)

// PortmapService ...
type PortmapService struct {
	rpcv2.RPCService
}

// TODO register service with portmapper

// NewPortmapService ...
func NewPortmapService() *PortmapService {
	portmapService := &PortmapService{
		RPCService: *rpcv2.NewRPCService("portmap", Program, Version),
	}

	portmapService.RegisterProcedure(PortmapProcedureNull, procedureNull)
	portmapService.RegisterProcedure(PortmapProcedureGetPort, procedureGetPort)

	return portmapService
}

func getPort(mapping Mapping) uint32 {
	fmt.Println("getPort")
	fmt.Println(mapping.Program)
	fmt.Println(mapping.Protocol)
	fmt.Println(mapping.Version)

	// TODO check mapping.Version (1) == mountv3.Version (3)
	if mapping.Program == mountv3.Program && mapping.Protocol == IPProtocolTCP {
		return 892
	} else if mapping.Program == nfsv3.Program && mapping.Protocol == IPProtocolTCP && mapping.Version == nfsv3.Version {
		return 2049
	}

	return ProgramNotAvailable
}
