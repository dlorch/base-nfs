package nfsv3

import (
	"bytes"
	"encoding/binary"

	"github.com/dlorch/nfsv3/rpcv2"
)

// ReadDirPlus3ResultOK (struct READDIRPLUS3resok)
type ReadDirPlus3ResultOK struct {
	ReadDirPlus3Result
	DirectoryAttributes PostOperationAttributes
	CookieVerifier      CookieVerifier3
	Reply               DirListPlus3
}

// ReadDirPlus3Result (union READDIRPLUS3res)
type ReadDirPlus3Result struct {
	status uint32
}

// ToBytes serializes the PathConf3ResultOK to be sent back to the client
func (reply *ReadDirPlus3ResultOK) ToBytes() ([]byte, error) {
	var responseBuffer = new(bytes.Buffer)
	var responseBytes = []byte{}

	err := binary.Write(responseBuffer, binary.BigEndian, &reply.ReadDirPlus3Result)

	if err != nil {
		return responseBytes, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.DirectoryAttributes)

	if err != nil {
		return responseBytes, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.CookieVerifier)

	if err != nil {
		return responseBytes, err
	}

	for _, entry := range reply.Reply.Entries {
		valueFollowsYes := uint32(1)

		err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)

		if err != nil {
			return responseBytes, err
		}

		err = binary.Write(responseBuffer, binary.BigEndian, &entry.FileID)

		if err != nil {
			return responseBytes, err
		}

		fileNameLength := uint32(len(entry.FileName3))
		numFillBytes := int(4 - (fileNameLength % 4))
		fillByte := uint8(0)

		err = binary.Write(responseBuffer, binary.BigEndian, &fileNameLength)

		if err != nil {
			return responseBytes, err
		}

		_, err = responseBuffer.WriteString(entry.FileName3)

		if err != nil {
			return responseBytes, err
		}

		for i := 0; i < numFillBytes; i++ {
			err = binary.Write(responseBuffer, binary.BigEndian, &fillByte)

			if err != nil {
				return responseBytes, err
			}
		}

		err = binary.Write(responseBuffer, binary.BigEndian, &entry.Cookie)

		if err != nil {
			return responseBytes, err
		}

		err = binary.Write(responseBuffer, binary.BigEndian, &entry.NameAttributes)

		if err != nil {
			return responseBytes, err
		}

		err = binary.Write(responseBuffer, binary.BigEndian, &entry.NameHandle.HandleFollows)

		if err != nil {
			return responseBytes, err
		}

		handleLength := uint32(len(entry.NameHandle.Handle))

		err = binary.Write(responseBuffer, binary.BigEndian, &handleLength)

		if err != nil {
			return responseBytes, err
		}

		_, err = responseBuffer.Write(entry.NameHandle.Handle)

		if err != nil {
			return responseBytes, err
		}
	}

	valueFollowsNo := uint32(0)

	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsNo)

	if err != nil {
		return responseBytes, err
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &reply.Reply.EOF)

	if err != nil {
		return responseBytes, err
	}

	responseBytes = make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())

	return responseBytes, nil
}

func nfsProcedure3ReadDirPlus(procedureArguments []byte) (rpcv2.Serializable, error) {
	// parse request
	// TODO

	// prepare result
	readDirPlusResult := &ReadDirPlus3ResultOK{
		ReadDirPlus3Result: ReadDirPlus3Result{
			status: NFS3OK,
		},
		DirectoryAttributes: PostOperationAttributes{
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
		CookieVerifier: [NFS3CookieVerifierSize]byte{},
		Reply: DirListPlus3{
			Entries: []EntryPlus3{
				EntryPlus3{
					FileID:    2,
					FileName3: "..",
					Cookie:    6457138716124813847,
					NameAttributes: PostOperationAttributes{
						AttributesFollow: 1,
						ObjectAttributes: FileAttr3{
							Typ:              2,
							Mode:             040777,
							Nlink:            15,
							UID:              0,
							GID:              0,
							Size:             4096,
							Used:             4096,
							Specdata1:        0,
							Specdata2:        0,
							Fsid:             0x388e4346cfc706a8,
							Fileid:           2,
							Atimeseconds:     1562969613,
							Atimenanoseconds: 760001904,
							Mtimeseconds:     1562969597,
							Mtimenanoseconds: 560001387,
							Ctimeseconds:     1562969597,
							Ctimenanoseconds: 560001387,
						},
					},
					NameHandle: PostOperationFileHandle3{
						HandleFollows: 1,
						Handle:        []byte{0x01, 0x00, 0x07, 0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x06, 0xc7, 0xcf, 0x46, 0x43, 0x8e, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
				},
				EntryPlus3{
					FileID:    16,
					FileName3: ".",
					Cookie:    6684891493313481230,
					NameAttributes: PostOperationAttributes{
						AttributesFollow: 1,
						ObjectAttributes: FileAttr3{
							Typ:              2,
							Mode:             040755,
							Nlink:            15,
							UID:              0,
							GID:              0,
							Size:             4096,
							Used:             4096,
							Specdata1:        0,
							Specdata2:        0,
							Fsid:             0x388e4346cfc706a8,
							Fileid:           2,
							Atimeseconds:     1562969613,
							Atimenanoseconds: 760001904,
							Mtimeseconds:     1562969597,
							Mtimenanoseconds: 560001387,
							Ctimeseconds:     1562969597,
							Ctimenanoseconds: 560001387,
						},
					},
					NameHandle: PostOperationFileHandle3{
						HandleFollows: 1,
						Handle:        []byte{0x01, 0x00, 0x07, 0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x06, 0xc7, 0xcf, 0x46, 0x43, 0x8e, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
				},
				EntryPlus3{
					FileID:    40243830,
					FileName3: "gopher.go",
					Cookie:    3621999153351014942,
					NameAttributes: PostOperationAttributes{
						AttributesFollow: 1,
						ObjectAttributes: FileAttr3{
							Typ:              1,
							Mode:             0100666,
							Nlink:            1,
							UID:              1027,
							GID:              100,
							Size:             292,
							Used:             8192,
							Specdata1:        0,
							Specdata2:        0,
							Fsid:             0x388e4346cfc706a8,
							Fileid:           40243830,
							Atimeseconds:     1456162928,
							Atimenanoseconds: 85375909,
							Mtimeseconds:     1389825403,
							Mtimenanoseconds: 480233665,
							Ctimeseconds:     1419273932,
							Ctimenanoseconds: 807093921,
						},
					},
					NameHandle: PostOperationFileHandle3{
						HandleFollows: 1,
						Handle:        []byte{0x01, 0x00, 0x07, 0x02, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x06, 0xc7, 0xcf, 0x46, 0x43, 0x8e, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x76, 0x12, 0x66, 0x02, 0x6d, 0x85, 0xd2, 0x28, 0x10, 0x00, 0x00, 0x00, 0xd9, 0x3c, 0x6d, 0x78},
					},
				},
			},
			EOF: 1,
		},
	}

	return readDirPlusResult, nil
}
