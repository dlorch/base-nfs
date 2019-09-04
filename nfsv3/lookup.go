// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Procedure 3: LOOKUP -  Lookup filename
// https://tools.ietf.org/html/rfc1813#page-37

package nfsv3

// Lookup3Args ...
type Lookup3Args struct {
	// what diropargs3
}

// Lookup3ResOK ...
type Lookup3ResOK struct {
	Object        NFSFH3
	ObjAttributes PostOperationAttributes
	DirAttributes PostOperationAttributes
}

// Lookup3ResFail ...
type Lookup3ResFail struct {
	DirAttributes PostOperationAttributes
}

// Lookup3Res ...
type Lookup3Res struct {
	Status  uint32         `xdr:"switch"`
	ResOK   Lookup3ResOK   `xdr:"case=0"`
	ResFail Lookup3ResFail `xdr:"default"`
}

// Lookup3 (NFSPROC3_LOOKUP) searches a directory for a specific name
// and returns the file handle for the corresponding file system object.
func Lookup3(arg []byte) (interface{}, error) {
	res := &Lookup3Res{
		Status: NFS3OK,
		ResOK: Lookup3ResOK{
			Object: NFSFH3{
				Data: []byte{1, 2, 3},
			},
			ObjAttributes: PostOperationAttributes{
				AttributesFollow: 1,
				ObjectAttributes: FileAttr3{
					Typ:              2,
					Mode:             040777,
					Nlink:            4,
					UID:              0,
					GID:              0,
					Size:             4096,
					Used:             8192,
					Specdata1:        0,
					Specdata2:        0,
					Fsid:             0x388e4346cfc706a8,
					Fileid:           16,
					Atimeseconds:     1563137262,
					Atimenanoseconds: 460002975,
					Mtimeseconds:     1537128120,
					Mtimenanoseconds: 839607220,
					Ctimeseconds:     1537128120,
					Ctimenanoseconds: 839607220,
				},
			},
			DirAttributes: PostOperationAttributes{
				AttributesFollow: 1,
				ObjectAttributes: FileAttr3{
					Typ:              2,
					Mode:             040777,
					Nlink:            4,
					UID:              0,
					GID:              0,
					Size:             4096,
					Used:             8192,
					Specdata1:        0,
					Specdata2:        0,
					Fsid:             0x388e4346cfc706a8,
					Fileid:           16,
					Atimeseconds:     1563137262,
					Atimenanoseconds: 460002975,
					Mtimeseconds:     1537128120,
					Mtimenanoseconds: 839607220,
					Ctimeseconds:     1537128120,
					Ctimenanoseconds: 839607220,
				},
			},
		},
	}
	return res, nil
}
