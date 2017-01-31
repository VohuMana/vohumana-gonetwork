package conduitServer

type packet interface {
    getData() []byte
    getSize() uint
}

type rawPacket struct {
    data []byte
    numBytes uint
}

func newRawPacket(data []byte, length uint) rawPacket {
    return rawPacket {
        data: data,
        numBytes: length,
    }
}

func (p rawPacket) getData() []byte {
    return p.data
}

func (p rawPacket) getSize() uint {
    return p.numBytes
}