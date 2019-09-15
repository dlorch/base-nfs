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
	ObjAttributes PostOpAttr
	DirAttributes PostOpAttr
}

// Lookup3ResFail ...
type Lookup3ResFail struct {
	DirAttributes PostOpAttr
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
			ObjAttributes: PostOpAttr{
				AttributesFollow: 1,
				ObjectAttributes: FAttr3{
					Type:  2,
					Mode:  040777,
					Nlink: 4,
					UID:   0,
					GID:   0,
					Size:  4096,
					Used:  8192,
					RDev: SpecData3{
						SpecData1: 0,
						SpecData2: 0,
					},
					FSID:   0x388e4346cfc706a8,
					FileID: 16,
					ATime: NFSTime3{
						Seconds:  1563137262,
						NSeconds: 460002975,
					},
					MTime: NFSTime3{
						Seconds:  1537128120,
						NSeconds: 839607220,
					},
					CTime: NFSTime3{
						Seconds:  1537128120,
						NSeconds: 839607220,
					},
				},
			},
			DirAttributes: PostOpAttr{
				AttributesFollow: 1,
				ObjectAttributes: FAttr3{
					Type:  2,
					Mode:  040777,
					Nlink: 4,
					UID:   0,
					GID:   0,
					Size:  4096,
					Used:  8192,
					RDev: SpecData3{
						SpecData1: 0,
						SpecData2: 0,
					},
					FSID:   0x388e4346cfc706a8,
					FileID: 16,
					ATime: NFSTime3{
						Seconds:  1563137262,
						NSeconds: 460002975,
					},
					MTime: NFSTime3{
						Seconds:  1537128120,
						NSeconds: 839607220,
					},
					CTime: NFSTime3{
						Seconds:  1537128120,
						NSeconds: 839607220,
					},
				},
			},
		},
	}
	return res, nil
}
