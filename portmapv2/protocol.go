// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package portmapv2

// Constants for port mapper
const (
	Program                 uint32 = 100000 // Portmap service program number (PMAP_PROG)
	Version                 uint32 = 2      // Portmap service version number
	PortmapProcedureNull    uint32 = 0      // PMAPPROC_NULL
	PortmapProcedureSet     uint32 = 1      // PMAPPROC_SET
	PortmapProcedureUnset   uint32 = 2      // PMAPPROC_UNSET
	PortmapProcedureGetPort uint32 = 3      // PMAPPROC_GETPORT
	PortmapProcedureDump    uint32 = 4      // PMAPPROC_DUMP
	PortmapProcedureCallIt  uint32 = 5      // PMAPPROC_CALLIT
	IPProtocolTCP           uint32 = 6      // protocol number for TCP/IP
	IPProtocolUDP           uint32 = 17     // protocol number for UCP/IP
	ProgramNotAvailable     uint32 = 0      // Port value of zero means the program has not been registered
)
