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

// Cookie3 (cookie3)
type Cookie3 uint64

// CookieVerifier3 (cookieverf3)
type CookieVerifier3 [NFS3CookieVerifierSize]byte

// FileAttr3 (struct fattr3)
type FileAttr3 struct {
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

// NFSFH3 (struct nfs_fh3)
type NFSFH3 struct {
	Data []byte
}

// PostOperationAttributes (union post_op_attr)
type PostOperationAttributes struct {
	AttributesFollow uint32    `xdr:"switch"` // TODO bool
	ObjectAttributes FileAttr3 `xdr:"case=1"`
}

// PostOperationFileHandle3 (union post_op_fh3)
type PostOperationFileHandle3 struct {
	HandleFollows uint32 // bool
	Handle        []byte // TODO struct nfs_fh3
}

// EntryPlus3 (struct entryplus3)
type EntryPlus3 struct {
	FileID         uint64
	FileName3      string
	Cookie         Cookie3
	NameAttributes PostOperationAttributes
	NameHandle     PostOperationFileHandle3
}

// DirListPlus3 (struct dirlistplus3)
type DirListPlus3 struct {
	Entries []EntryPlus3
	EOF     uint32 // bool
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
	NFS3CookieVerifierSize     uint32 = 8      // The size in bytes of the opaque cookie verifier passed by READDIR and READDIRPLUS (NFS3_COOKIEVERFSIZE)
)
