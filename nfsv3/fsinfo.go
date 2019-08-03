package nfsv3

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dlorch/nfsv3/rpcv2"
)

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
