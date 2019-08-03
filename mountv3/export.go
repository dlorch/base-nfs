package mountv3

import (
	"bytes"
	"encoding/binary"

	"github.com/dlorch/nfsv3/rpcv2"
)

// ToBytes serializes the ExportNode to be sent back to the client
func (reply *ExportNode) ToBytes() ([]byte, error) {
	var responseBuffer = new(bytes.Buffer)

	// --- mount service body

	var valueFollowsYes uint32
	valueFollowsYes = 1

	var valueFollowsNo uint32
	valueFollowsNo = 0

	directoryContents := "/volume1/Public"
	directoryLength := uint32(len(directoryContents))

	groupContents := "*"
	groupLength := uint32(len(groupContents))

	fillBytes := uint8(0)

	err := binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)
	err = binary.Write(responseBuffer, binary.BigEndian, &directoryLength)
	_, err = responseBuffer.Write([]byte(directoryContents))
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)

	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsYes)
	err = binary.Write(responseBuffer, binary.BigEndian, &groupLength)
	_, err = responseBuffer.Write([]byte(groupContents))
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &fillBytes)
	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsNo)

	err = binary.Write(responseBuffer, binary.BigEndian, &valueFollowsNo)

	responseBytes := make([]byte, responseBuffer.Len())
	copy(responseBytes, responseBuffer.Bytes())
	return responseBytes, err
}

func mountProcedure3export(procedureArguments []byte) (rpcv2.Serializable, error) {
	return &ExportNode{}, nil
}
