// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package portmapv2

import (
	"bytes"
	"encoding/binary"
)

// GetPortResult represents the requested port number
type GetPortResult struct {
	Port uint32
}

// Mapping of (program, version, protocol) to port number (RFC1057: struct_mapping)
type Mapping struct {
	Program  uint32
	Version  uint32
	Protocol uint32
	Port     uint32
}

func procedureGetPort(procedureArguments []byte) (interface{}, error) {
	var requestBody = bytes.NewBuffer(procedureArguments)
	var mapping Mapping

	err := binary.Read(requestBody, binary.BigEndian, &mapping)

	if err != nil {
		return &GetPortResult{Port: ProgramNotAvailable}, err
	}

	// TODO check callBody.Version == portmapv2.Version

	port := getPort(mapping)

	return &GetPortResult{Port: port}, nil
}
