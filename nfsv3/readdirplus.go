// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nfsv3

// EntryPlus3 (struct entryplus3)
type EntryPlus3 struct {
	ValueFollows   uint32 `xdr:"switch"`
	FileID         uint64 `xdr:"case=1"`
	FileName3      string
	Cookie         Cookie3
	NameAttributes PostOpAttr
	NameHandle     PostOperationFileHandle3
	NextEntry      *EntryPlus3
}

// DirListPlus3 (struct dirlistplus3)
type DirListPlus3 struct {
	Entries *EntryPlus3
	EOF     uint32 // bool
}

// ReadDirPlus3ResultOK (struct READDIRPLUS3resok)
type ReadDirPlus3ResultOK struct {
	ReadDirPlus3Result
	DirectoryAttributes PostOpAttr
	CookieVerifier      CookieVerifier3
	Reply               DirListPlus3
}

// ReadDirPlus3Result (union READDIRPLUS3res)
type ReadDirPlus3Result struct {
	Status uint32
}

func nfsProcedure3ReadDirPlus(procedureArguments []byte) (interface{}, error) {
	// parse request
	// TODO

	// prepare result
	readDirPlusResult := &ReadDirPlus3ResultOK{
		ReadDirPlus3Result: ReadDirPlus3Result{
			Status: NFS3OK,
		},
		DirectoryAttributes: PostOpAttr{
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
		CookieVerifier: [NFS3CookieVerifierSize]byte{},
		Reply: DirListPlus3{
			Entries: &EntryPlus3{
				ValueFollows: 1,
				FileID:       2,
				FileName3:    "..",
				Cookie:       6457138716124813847,
				NameAttributes: PostOpAttr{
					AttributesFollow: 1,
					ObjectAttributes: FAttr3{
						Type:  2,
						Mode:  040777,
						Nlink: 15,
						UID:   0,
						GID:   0,
						Size:  4096,
						Used:  4096,
						RDev: SpecData3{
							SpecData1: 0,
							SpecData2: 0,
						},
						FSID:   0x388e4346cfc706a8,
						FileID: 2,
						ATime: NFSTime3{
							Seconds:  1562969613,
							NSeconds: 760001904,
						},
						MTime: NFSTime3{
							Seconds:  1562969597,
							NSeconds: 560001387,
						},
						CTime: NFSTime3{
							Seconds:  1562969597,
							NSeconds: 560001387,
						},
					},
				},
				NameHandle: PostOperationFileHandle3{
					HandleFollows: 1,
					Handle:        []byte{0x01, 0x00, 0x07, 0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x06, 0xc7, 0xcf, 0x46, 0x43, 0x8e, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				},
				NextEntry: &EntryPlus3{
					ValueFollows: 1,
					FileID:       16,
					FileName3:    ".",
					Cookie:       6684891493313481230,
					NameAttributes: PostOpAttr{
						AttributesFollow: 1,
						ObjectAttributes: FAttr3{
							Type:  2,
							Mode:  040755,
							Nlink: 15,
							UID:   0,
							GID:   0,
							Size:  4096,
							Used:  4096,
							RDev: SpecData3{
								SpecData1: 0,
								SpecData2: 0,
							},
							FSID:   0x388e4346cfc706a8,
							FileID: 2,
							ATime: NFSTime3{
								Seconds:  1562969613,
								NSeconds: 760001904,
							},
							MTime: NFSTime3{
								Seconds:  1562969597,
								NSeconds: 560001387,
							},
							CTime: NFSTime3{
								Seconds:  1562969597,
								NSeconds: 560001387,
							},
						},
					},
					NameHandle: PostOperationFileHandle3{
						HandleFollows: 1,
						Handle:        []byte{0x01, 0x00, 0x07, 0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x06, 0xc7, 0xcf, 0x46, 0x43, 0x8e, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					NextEntry: &EntryPlus3{
						ValueFollows: 1,
						FileID:       40243830,
						FileName3:    "gopher.go",
						Cookie:       3621999153351014942,
						NameAttributes: PostOpAttr{
							AttributesFollow: 1,
							ObjectAttributes: FAttr3{
								Type:  1,
								Mode:  0100666,
								Nlink: 1,
								UID:   1027,
								GID:   100,
								Size:  292,
								Used:  8192,
								RDev: SpecData3{
									SpecData1: 0,
									SpecData2: 0,
								},
								FSID:   0x388e4346cfc706a8,
								FileID: 40243830,
								ATime: NFSTime3{
									Seconds:  1456162928,
									NSeconds: 85375909,
								},
								MTime: NFSTime3{
									Seconds:  1389825403,
									NSeconds: 480233665,
								},
								CTime: NFSTime3{
									Seconds:  1419273932,
									NSeconds: 807093921,
								},
							},
						},
						NameHandle: PostOperationFileHandle3{
							HandleFollows: 1,
							Handle:        []byte{0x01, 0x00, 0x07, 0x02, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x06, 0xc7, 0xcf, 0x46, 0x43, 0x8e, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x76, 0x12, 0x66, 0x02, 0x6d, 0x85, 0xd2, 0x28, 0x10, 0x00, 0x00, 0x00, 0xd9, 0x3c, 0x6d, 0x78},
						},
						NextEntry: &EntryPlus3{
							ValueFollows: 0,
						},
					},
				},
			},
			EOF: 1,
		},
	}

	return readDirPlusResult, nil
}
