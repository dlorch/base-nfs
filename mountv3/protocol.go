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

package mountv3

import (
	"bytes"
	"encoding/binary"

	"github.com/dlorch/nfsv3/rpcv2"
)

// MountVoidReply is an empty reply
type MountVoidReply struct{}

// GroupNode (RFC1813: struct groupnode)
type GroupNode struct {
}

// ExportNode (RFC1813: struct exportnode)
type ExportNode struct {
}

// MountResult3 (RFC1813: struct mountres3)
type MountResult3 struct {
	Status         uint32
	MountResult3OK MountResult3OK
}

// MountResult3OK (RFC1813: struct mountres3_ok)
type MountResult3OK struct {
	FileHandle3 []byte
	AuthFlavors []uint32
}

// Constants for mount protocol (RFC1813)
const (
	Program                    uint32 = 100005 // Mount service program number
	Version                    uint32 = 3      // Mount service version
	MountProcedure3Null        uint32 = 0      // MOUNTPROC3_NULL
	MountProcedure3Mount       uint32 = 1      // MOUNTPROC3_MNT
	MountProcedure3Dump        uint32 = 2      // MOUNTPROC3_DUMP
	MountProcedure3Unmount     uint32 = 3      // MOUNTPROC3_UMNT
	MountProcedure3UnmountAll  uint32 = 4      // MOUNTPROC3_UMNTALL
	MountProcedure3Export      uint32 = 5      // MOUNTPROC3_EXPORT
	Mount3OK                   uint32 = 0      // MNT3_OK: no error
	Mount3ErrorPermissions     uint32 = 1      // MNT3ERR_PERM: Not owner
	Mount3ErrorNoEntry         uint32 = 2      // MNT3ERR_NOENT: No such file or directory
	Mount3ErrorIO              uint32 = 5      // MNT3ERR_IO: I/O error
	Mount3ErrorAccess          uint32 = 13     // MNT3ERR_ACCES: Permission denied
	Mount3ErrorNotDirectory    uint32 = 20     // MNT3ERR_NOTDIR: Not a directory
	Mount3ErrorInvalidArgument uint32 = 22     // MNT3ERR_INVAL: Invalid argument
	Mount3ErrorNameTooLong     uint32 = 63     // MNT3ERR_NAMETOOLONG: Filename too long
	Mount3ErrorNotSupported    uint32 = 10004  // MNT3ERR_NOTSUPP: Operation not supported
	Mount3ErrorServerFault     uint32 = 10006  // MNT3ERR_SERVERFAULT: A failure on the server
)

// ----- MountProcedure3Null

// ToBytes serializes the VoidReply to be sent back to the client
func (reply *MountVoidReply) ToBytes() ([]byte, error) {
	return []byte{}, nil
}

func mountProcedure3Null(procedureArguments []byte) (rpcv2.Serializable, error) {
	return &MountVoidReply{}, nil
}

// ----- MountProcedure3Export

// ToBytes serializes the ExportNode to be sent back to the client
func (reply *ExportNode) ToBytes() ([]byte, error) {
	var responseBuffer = new(bytes.Buffer)

	// --- mount service body

	var valueFollowsYes uint32
	valueFollowsYes = 1

	var valueFollowsNo uint32
	valueFollowsNo = 0

	directoryContents := "/volume1/Public"
	directoryLength := uint32(len(directoryContents))

	groupContents := "*"
	groupLength := uint32(len(groupContents))

	fillBytes := uint8(0)

	err := binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)
	err = binary.Write(responseBuffer, binary.BigEndian, &directoryLength)
	_, err = responseBuffer.Write([]byte(directoryContents))
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)

	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)
	err = binary.Write(responseBuffer, binary.BigEndian, &groupLength)
	_, err = responseBuffer.Write([]byte(groupContents))
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsNo)

	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsNo)

	responseBytes := make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())
	return responseBytes, err
}

func mountProcedure3export(procedureArguments []byte) (rpcv2.Serializable, error) {
	return &ExportNode{}, nil
}

// ----- MountProcedure3Mount

// ToBytes serializes the MountResult3 to be sent back to the client
func (reply *MountResult3) ToBytes() ([]byte, error) {
	var responseBuffer = new(bytes.Buffer)

	err := binary.Write(responseBuffer, binary.BigEndian, Mount3OK)

	err = binary.Write(responseBuffer, binary.BigEndian, uint32(4))  // length of file handle
	err = binary.Write(responseBuffer, binary.BigEndian, uint32(42)) // file handle

	err = binary.Write(responseBuffer, binary.BigEndian, uint32(1))                // number of auth flavors
	err = binary.Write(responseBuffer, binary.BigEndian, rpcv2.AuthenticationUNIX) // allowed flavors

	responseBytes := make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())
	return responseBytes, err
}

func mountProcedure3mount(procedureArguments []byte) (rpcv2.Serializable, error) {
	// parse request
	requestBuffer := bytes.NewBuffer(procedureArguments)

	var dirPathLength uint32

	err := binary.Read(requestBuffer, binary.BigEndian, &dirPathLength)

	if err != nil {
		return nil, err
	}

	dirPathName := make([]byte, dirPathLength) // TODO check MNTPATHLEN

	err = binary.Read(requestBuffer, binary.BigEndian, &dirPathName)

	if err != nil {
		return nil, err
	}

	return &MountResult3{}, nil
}
