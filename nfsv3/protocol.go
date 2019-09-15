// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nfsv3

// Cookie3 (cookie3)
type Cookie3 uint64

// CookieVerifier3 (cookieverf3)
type CookieVerifier3 [NFS3CookieVerifierSize]byte

// FAttr3 (struct fattr3)
type FAttr3 struct {
	Typ              uint32
	Mode             uint32
	Nlink            uint32
	UID              uint32
	GID              uint32
	Size             uint64
	Used             uint64
	Specdata1        uint32
	Specdata2        uint32
	Fsid             uint64
	Fileid           uint64
	Atimeseconds     uint32
	Atimenanoseconds uint32
	Mtimeseconds     uint32
	Mtimenanoseconds uint32
	Ctimeseconds     uint32
	Ctimenanoseconds uint32
}

// PostOperationAttributes (union post_op_attr)
type PostOperationAttributes struct {
	AttributesFollow uint32 `xdr:"switch"` // TODO bool
	ObjectAttributes FAttr3 `xdr:"case=1"`
}

// PostOperationFileHandle3 (union post_op_fh3)
type PostOperationFileHandle3 struct {
	HandleFollows uint32 // bool
	Handle        []byte // TODO struct nfs_fh3
}

// RPC Constants for NFS3 Protocol
const (
	Program uint32 = 100003 // Mount service program number
	Version uint32 = 3      // Mount service program version
)

// Sizes, given in decimal bytes, of various XDR structures
const (
	NFS3CookieVerifierSize uint32 = 8 // The size in bytes of the opaque cookie verifier passed by READDIR and READDIRPLUS (NFS3_COOKIEVERFSIZE)
)

// Returned with every procedure's results except for the NULL procedure (enum nfsstat3)
const (
	NFS3OK             uint32 = 0     // Indicates the call completed successfully (NFS3_OK)
	NFS3ErrPerm        uint32 = 1     // Not owner (NFS3ERR_PERM)
	NFS3ErrNoEnt       uint32 = 2     // No such file or directory (NFS3ERR_NOENT)
	NFS3ErrIO          uint32 = 5     // I/O error. A hard error (for example, a disk error) occurred (NFS3ERR_IO)
	NFS3ErrNXIO        uint32 = 6     // I/O error. No such device or address (NFS3ERR_NXIO)
	NFS3ErrAcces       uint32 = 13    // Permission denied (NFS3ERR_ACCES)
	NFS3ErrExist       uint32 = 17    // File exists (NFS3ERR_EXIST)
	NFS3ErrXDev        uint32 = 18    // Attempt to do a cross-device hard link (NFS3ERR_XDEV)
	NFS3ErrNoDev       uint32 = 19    // No such device (NFS3ERR_NODEV)
	NFS3ErrNotDir      uint32 = 20    // Not a directory (NFS3ERR_NOTDIR)
	NFS3ErrIsDir       uint32 = 21    // Is a directory (NFS3ERR_ISDIR)
	NFS3ErrInval       uint32 = 22    // Invalid argument or unsupported argument (NFS3ERR_INVAL)
	NFS3ErrFBig        uint32 = 27    // File too large (NFS3ERR_FBIG)
	NFS3ErrNoSpc       uint32 = 28    // No space left on device (NFS3ERR_NOSPC)
	NFS3ErrROFS        uint32 = 30    // Read-only file system (NFS3ERR_ROFS)
	NFS3ErrMLink       uint32 = 31    // Too many hard links (NFS3ERR_MLINK)
	NFS3ErrNameTooLong uint32 = 63    // The filename in an operation was too long (NFS3ERR_NAMETOOLONG)
	NFS3ErrNotEmpty    uint32 = 66    // An attempt was made to remove a directory that was not empty (NFS3ERR_NOTEMPTY)
	NFS3ErrDQuot       uint32 = 69    // Resource (quota) hard limit exceeded (NFS3ERR_DQUOT)
	NFS3ErrStale       uint32 = 70    // Invalid file handle (NFS3ERR_STALE)
	NFS3ErrRemote      uint32 = 71    // Too many levels of remote in path (NFS3ERR_REMOTE)
	NFS3ErrBadHandle   uint32 = 10001 // Illegal NFS file handle (NFS3ERR_BADHANDLE)
	NFS3ErrNotSync     uint32 = 10002 // Update synchronization mismatch was detected during a SETATTR operation (NFS3ERR_NOT_SYNC)
	NFS3ErrBadCookie   uint32 = 10003 // READDIR or READDIRPLUS cookie is stale (NFS3ERR_BAD_COOKIE)
	NFS3ErrNotSupp     uint32 = 10004 // Operation is not supported (NFS3ERR_NOTSUPP)
	NFS3ErrTooSmall    uint32 = 10005 // Buffer or request is too small (NFS3ERR_TOOSMALL)
	NFS3ErrServerFault uint32 = 10006 // An error occurred on the server which does not map to any of the legal NFS version 3 protocol error values (NFS3ERR_SERVERFAULT)
	NFS3ErrBadType     uint32 = 10007 // An attempt was made to create an object of a type not (NFS3ERR_BADTYPE)
	NFS3ErrJukeBox     uint32 = 10008 // The server initiated the request, but was not able to complete it in a timely fashion (NFS3ERR_JUKEBOX)
)

// Type of a file (enum ftype3)
const (
	NF3Reg  uint32 = 1 // regular file (NF3REG)
	NF3Dir  uint32 = 2 // directory (NF3DIR)
	NF3Blk  uint32 = 3 // block special device file (NF3BLK)
	NF3Chr  uint32 = 4 // character special device file (NF3CHR)
	NF3Lnk  uint32 = 5 // symbolic link (NF3LNK)
	NF3Sock uint32 = 6 // socket (NF3SOCK)
	NF3FIFO uint32 = 7 // named pipe (NF3FIFO)
)

// SpecData3 is returned as part of the FAttr3 structure (struct specdata3)
type SpecData3 struct {
	SpecData1 uint32
	SpecData2 uint32
}

// NFSFH3 describes a file handle which contains all the information
// the server needs to distuinguish an individual file (struct nfs_fh3)
type NFSFH3 struct {
	Data []byte
}

// NFSTime3 gives the number of seconds and nanoseconds since midnight
// January 1, 1970 Greenwich Mean Time (struct nfstime3)
type NFSTime3 struct {
	Seconds  uint32
	NSeconds uint32
}

// RPC procedure numbers
const (
	NFSProcedure3Null          uint32 = 0  // NFSPROC3_NULL
	NFSProcedure3GetAttributes uint32 = 1  // NFSPROC3_GETATTR
	NFSProcedure3SetAttributes uint32 = 2  // NFSPROC3_SETATTR
	NFSProcedure3Lookup        uint32 = 3  // NFSPROC3_LOOKUP
	NFSProcedure3Access        uint32 = 4  // NFSPROC3_ACCESS
	NFSProcedure3Readlink      uint32 = 5  // NFSPROC3_READLINK
	NFSProcedure3Read          uint32 = 6  // NFSPROC3_READ
	NFSProcedure3Write         uint32 = 7  // NFSPROC3_WRITE
	NFSProcedure3Create        uint32 = 8  // NFSPROC3_CREATE
	NFSProcedure3MkDir         uint32 = 9  // NFSPROC3_MKDIR
	NFSProcedure3Symlink       uint32 = 10 // NFSPROC3_SYMLINK
	NFSProcedure3MkNod         uint32 = 11 // NFSPROC3_MKNOD
	NFSProcedure3Remove        uint32 = 12 // NFSPROC3_REMOVE
	NFSProcedure3RmDir         uint32 = 13 // NFSPROC3_RMDIR
	NFSProcedure3Rename        uint32 = 14 // NFSPROC3_RENAME
	NFSProcedure3Link          uint32 = 15 // NFSPROC3_LINK
	NFSProcedure3ReadDir       uint32 = 16 // NFSPROC3_READDIR
	NFSProcedure3ReadDirPlus   uint32 = 17 // NFSPROC3_READDIRPLUS
	NFSProcedure3FSStat        uint32 = 18 // NFSPROC3_FSSTAT
	NFSProcedure3FSInfo        uint32 = 19 // NFSPROC3_FSINFO
	NFSProcedure3PathConf      uint32 = 20 // NFSPROC3_PATHCONF
	NFSProcedure3Commint       uint32 = 21 // NFSPROC3_COMMIT
)
