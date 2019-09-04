// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mountv3

// Constants for mount protocol (RFC1813)
const (
	Program                    uint32 = 100005 // Mount service program number
	Version                    uint32 = 3      // Mount service version
	MountProcedure3Null        uint32 = 0      // MOUNTPROC3_NULL
	MountProcedure3Dump        uint32 = 2      // MOUNTPROC3_DUMP
	MountProcedure3Unmount     uint32 = 3      // MOUNTPROC3_UMNT
	MountProcedure3UnmountAll  uint32 = 4      // MOUNTPROC3_UMNTALL
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
