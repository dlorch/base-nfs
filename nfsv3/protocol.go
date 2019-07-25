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

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dlorch/nfsv3/rpcv2"
)

// VoidReply is an empty reply
type VoidReply struct{}

// FileAttr3 (struct fattr3)
type FileAttr3 struct {
	typ              uint32
	mode             uint32
	nlink            uint32
	uid              uint32
	gid              uint32
	size             uint64
	used             uint64
	specdata1        uint32
	specdata2        uint32
	fsid             uint64
	fileid           uint64
	atimeseconds     uint32
	atimenanoseconds uint32
	mtimeseconds     uint32
	mtimenanoseconds uint32
	ctimeseconds     uint32
	ctimenanoseconds uint32
}

// GetAttr3Args (struct FSINFOargs)
type GetAttr3Args struct {
	FileHandle []byte
}

// GetAttr3ResultOK (struct GETATTR3resok)
type GetAttr3ResultOK struct {
	GetAttr3Result
	ObjectAttributes FileAttr3
}

// GetAttr3Result (union GETATTR3res)
type GetAttr3Result struct {
	status uint32
}

// FSInfo3Args (struct FSINFOargs)
type FSInfo3Args struct {
	FileHandle []byte
}

// FSInfo3ResultOK (struct FSINFO3resok)
type FSInfo3ResultOK struct {
	FSInfo3Result
	objattributes        uint32 // TODO
	rtmax                uint32
	rtpref               uint32
	rtmult               uint32
	wtmax                uint32
	wtpref               uint32
	wtmult               uint32
	dtpref               uint32
	maxfilesize          uint64
	timedeltaseconds     uint32
	timedeltananoseconds uint32
	properties           uint32
}

// FSInfo3ResultFail (struct FSINFO3resfail)
type FSInfo3ResultFail struct {
	FSInfo3Result
	// TODO post_op_attr obj_attributes
}

// FSInfo3Result (union FSINFO3res)
type FSInfo3Result struct {
	status uint32
}

// PathConf3Args (struct PATHCONF3args)
type PathConf3Args struct {
	FileHandle []byte
}

// PathConf3ResultOK (struct PATHCONF3resok)
type PathConf3ResultOK struct {
	PathConf3Result
	objattributes   uint32 // TODO
	linkmax         uint32
	namemax         uint32
	notrunc         uint32 // TODO bool
	chownrestricted uint32 // TODO bool
	caseinsensitive uint32 // TODO bool
	casepreserving  uint32 // TODO bool
}

// TODO PATHCONF3resfail

// PathConf3Result (union PATHCONF3res)
type PathConf3Result struct {
	status uint32
}

// Constants for mount protocol (RFC1813)
const (
	Program                    uint32 = 100003 // Mount service program number
	Version                    uint32 = 3      // Mount service program version
	NFSProcedure3Null          uint32 = 0      // NFSPROC3_NULL
	NFSProcedure3GetAttributes uint32 = 1      // NFSPROC3_GETATTR
	NFSProcedure3SetAttributes uint32 = 2      // NFSPROC3_SETATTR
	NFSProcedure3Lookup        uint32 = 3      // NFSPROC3_LOOKUP
	NFSProcedure3Access        uint32 = 4      // NFSPROC3_ACCESS
	NFSProcedure3Readlink      uint32 = 5      // NFSPROC3_READLINK
	NFSProcedure3Read          uint32 = 6      // NFSPROC3_READ
	NFSProcedure3Write         uint32 = 7      // NFSPROC3_WRITE
	NFSProcedure3Create        uint32 = 8      // NFSPROC3_CREATE
	NFSProcedure3MkDir         uint32 = 9      // NFSPROC3_MKDIR
	NFSProcedure3Symlink       uint32 = 10     // NFSPROC3_SYMLINK
	NFSProcedure3MkNod         uint32 = 11     // NFSPROC3_MKNOD
	NFSProcedure3Remove        uint32 = 12     // NFSPROC3_REMOVE
	NFSProcedure3RmDir         uint32 = 13     // NFSPROC3_RMDIR
	NFSProcedure3Rename        uint32 = 14     // NFSPROC3_RENAME
	NFSProcedure3Link          uint32 = 15     // NFSPROC3_LINK
	NFSProcedure3ReadDir       uint32 = 16     // NFSPROC3_READDIR
	NFSProcedure3ReadDirPlus   uint32 = 17     // NFSPROC3_READDIRPLUS
	NFSProcedure3FSStat        uint32 = 18     // NFSPROC3_FSSTAT
	NFSProcedure3FSInfo        uint32 = 19     // NFSPROC3_FSINFO
	NFSProcedure3PathConf      uint32 = 20     // NFSPROC3_PATHCONF
	NFSProcedure3Commint       uint32 = 21     // NFSPROC3_COMMIT
	NFS3OK                     uint32 = 0      // NFS3_OK
)

// ----- NFSProcedure3Null

// ToBytes serializes the VoidReply to be sent back to the client
func (reply *VoidReply) ToBytes() ([]byte, error) {
	return rpcv2.SerializeFixedSizeStruct(reply)
}

func nfsProcedure3Null(procedureArguments []byte) (rpcv2.Serializable, error) {
	return &VoidReply{}, nil
}

// ----- NFSProcedure3GetAttributes

// ToBytes serializes the GetAttr3ResultOK to be sent back to the client
func (reply *GetAttr3ResultOK) ToBytes() ([]byte, error) {
	return rpcv2.SerializeFixedSizeStruct(reply)
}

func nfsProcedure3GetAttributes(procedureArguments []byte) (rpcv2.Serializable, error) {
	// parse request
	// TODO

	// prepare result
	getAttrResult := &GetAttr3ResultOK{
		GetAttr3Result: GetAttr3Result{
			status: NFS3OK,
		},
		ObjectAttributes: FileAttr3{
			typ:              2,
			mode:             040777,
			nlink:            4,
			uid:              0,
			gid:              0,
			size:             4096,
			used:             8192,
			specdata1:        0,
			specdata2:        0,
			fsid:             0x388e4346cfc706a8,
			fileid:           16,
			atimeseconds:     1563137262,
			atimenanoseconds: 460002975,
			mtimeseconds:     1537128120,
			mtimenanoseconds: 839607220,
			ctimeseconds:     1537128120,
			ctimenanoseconds: 839607220,
		},
	}

	return getAttrResult, nil
}

// ----- NFSProcedure3Lookup

func nfsProcedure3Lookup(procedureArguments []byte) (rpcv2.Serializable, error) {
	fmt.Println("nfsProcedure3Access")
	return nil, nil
}

// ----- NFSProcedure3Access

func nfsProcedure3Access(procedureArguments []byte) (rpcv2.Serializable, error) {
	fmt.Println("nfsProcedure3Access")
	return nil, nil
}

// ----- NFSProcedure3FSInfo

// ToBytes serializes the FSInfo3ResultOK to be sent back to the client
func (reply *FSInfo3ResultOK) ToBytes() ([]byte, error) {
	return rpcv2.SerializeFixedSizeStruct(reply)
}

func nfsProcedure3FSInfo(procedureArguments []byte) (rpcv2.Serializable, error) {
	// parse request
	requestBuffer := bytes.NewBuffer(procedureArguments)

	var fileHandleLength uint32

	err := binary.Read(requestBuffer, binary.BigEndian, &fileHandleLength)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	fsInfoArgs := FSInfo3Args{
		FileHandle: make([]byte, fileHandleLength), // TODO unsafe?
	}

	err = binary.Read(requestBuffer, binary.BigEndian, &fsInfoArgs.FileHandle)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		// TODO
	}

	// prepare result
	fsInfoResult := &FSInfo3ResultOK{
		FSInfo3Result: FSInfo3Result{
			status: NFS3OK,
		},
		objattributes:        0,
		rtmax:                131072,
		rtpref:               131072,
		rtmult:               4096,
		wtmax:                131072,
		wtpref:               131072,
		wtmult:               4096,
		dtpref:               4096,
		maxfilesize:          8796093022207,
		timedeltaseconds:     1,
		timedeltananoseconds: 0,
		properties:           0x0000001b,
	}

	return fsInfoResult, nil
}

// ----- NFSProcedure3PathConf

// ToBytes serializes the PathConf3ResultOK to be sent back to the client
func (reply *PathConf3ResultOK) ToBytes() ([]byte, error) {
	return rpcv2.SerializeFixedSizeStruct(reply)
}

func nfsProcedure3PathConf(procedureArguments []byte) (rpcv2.Serializable, error) {
	// parse request
	// TODO

	// prepare result
	pathConfResult := &PathConf3ResultOK{
		PathConf3Result: PathConf3Result{
			status: NFS3OK,
		},
		objattributes:   0,
		linkmax:         32000,
		namemax:         255,
		notrunc:         0,
		chownrestricted: 1,
		caseinsensitive: 0,
		casepreserving:  1,
	}

	return pathConfResult, nil
}
