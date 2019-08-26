/*
	Port Mapper Protocol Specification Version 2 (RFC1057)

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
