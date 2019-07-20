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

	"github.com/dlorch/nfsv3/nfsv3"

	"github.com/dlorch/nfsv3/mountv3"
	"github.com/dlorch/nfsv3/rpcv2"
)

// TODO register service with portmapper

// NewPortmapService ...
func NewPortmapService() rpcv2.RPCService {
	rpcService := rpcv2.NewRPCService("portmap", Program, Version)

	rpcService.RegisterProcedure(ProcedureNull, procedureNull)
	rpcService.RegisterProcedure(ProcedureGetPort, procedureGetPort)

	return rpcService
}

func getPort(mapping Mapping) (port uint32, err error) {
	// TODO check mapping.Version (1) == mountv3.Version (3)
	if mapping.Program == mountv3.Program && mapping.Protocol == IPProtocolTCP {
		return 892, nil
	} else if mapping.Program == nfsv3.Program && mapping.Protocol == IPProtocolTCP && mapping.Version == nfsv3.Version {
		return 2049, nil
	}

	return port, fmt.Errorf("Unregistered program '%d' with protocol '%d'", mapping.Program, mapping.Protocol)
}
