package nfsv3

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// FSInfo3Args (struct FSINFOargs)
type FSInfo3Args struct {
	FileHandle []byte
}

// FSInfo3ResultOK (struct FSINFO3resok)
type FSInfo3ResultOK struct {
	FSInfo3Result
	Objattributes        uint32 // TODO
	Rtmax                uint32
	Rtpref               uint32
	Rtmult               uint32
	Wtmax                uint32
	Wtpref               uint32
	Wtmult               uint32
	Dtpref               uint32
	Maxfilesize          uint64
	Timedeltaseconds     uint32
	Timedeltananoseconds uint32
	Properties           uint32
}

// FSInfo3ResultFail (struct FSINFO3resfail)
type FSInfo3ResultFail struct {
	FSInfo3Result
	// TODO post_op_attr obj_attributes
}

// FSInfo3Result (union FSINFO3res)
type FSInfo3Result struct {
	Status uint32
}

func nfsProcedure3FSInfo(procedureArguments []byte) (interface{}, error) {
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
			Status: NFS3OK,
		},
		Objattributes:        0,
		Rtmax:                131072,
		Rtpref:               131072,
		Rtmult:               4096,
		Wtmax:                131072,
		Wtpref:               131072,
		Wtmult:               4096,
		Dtpref:               4096,
		Maxfilesize:          8796093022207,
		Timedeltaseconds:     1,
		Timedeltananoseconds: 0,
		Properties:           0x0000001b,
	}

	return fsInfoResult, nil
}
