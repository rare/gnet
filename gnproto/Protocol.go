package gnproto

import (
	"bytes"
	"encoding/binary"
)

type Header struct {
	Cmd		uint16
	Ver		uint16
	Seq		uint32
	Len		uint32
}

const (
	HEADER_SIZE			=	uint32(12)
)

const (
	CMD_HEART_BEAT		=	uint16(0)	//no body
)

func (this *Header) Serialize() ([]byte, error) {
	buf	:= new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, *this); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *Header) Deserialize(b []byte) error {
	buf := bytes.NewReader(b)
	if err := binary.Read(buf, binary.BigEndian, this); err != nil {
		return err
	}
	return nil
}
