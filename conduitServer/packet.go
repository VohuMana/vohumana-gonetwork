package conduitServer

import (

)

type Packet interface {
    GetData() []byte
    GetSize() uint
}

type RawPacket struct {
    data []byte
    numBytes uint
}

func NewRawPacket(data []byte, length uint) RawPacket {
    return RawPacket {
        data: data,
        numBytes: length,
    }
}

func (p RawPacket) GetData() []byte {
    return p.data
}

func (p RawPacket) GetSize() uint {
    return p.numBytes
}