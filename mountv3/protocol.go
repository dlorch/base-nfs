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
