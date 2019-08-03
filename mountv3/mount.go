package mountv3

import (
	"bytes"
	"encoding/binary"

	"github.com/dlorch/nfsv3/rpcv2"
)

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

// ToBytes serializes the MountResult3 to be sent back to the client
func (reply *MountResult3) ToBytes() ([]byte, error) {
	var responseBuffer = new(bytes.Buffer)

	err := binary.Write(responseBuffer, binary.BigEndian, Mount3OK)

	err = binary.Write(responseBuffer, binary.BigEndian, uint32(4))  // length of file handle
	err = binary.Write(responseBuffer, binary.BigEndian, uint32(42)) // file handle

	err = binary.Write(responseBuffer, binary.BigEndian, uint32(1))                // number of auth flavors
	err = binary.Write(responseBuffer, binary.BigEndian, rpcv2.AuthenticationUNIX) // allowed flavors

	responseBytes := make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())
	return responseBytes, err
}

func mountProcedure3mount(procedureArguments []byte) (rpcv2.Serializable, error) {
	// parse request
	requestBuffer := bytes.NewBuffer(procedureArguments)

	var dirPathLength uint32

	err := binary.Read(requestBuffer, binary.BigEndian, &dirPathLength)

	if err != nil {
		return nil, err
	}

	dirPathName := make([]byte, dirPathLength) // TODO check MNTPATHLEN

	err = binary.Read(requestBuffer, binary.BigEndian, &dirPathName)

	if err != nil {
		return nil, err
	}

	return &MountResult3{}, nil
}
