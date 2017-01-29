package conduitServer

import (
    "net"
)

type Client interface {
    Listen()
    GetConnection() net.Conn
    GetId() uint
}

type BasicClient struct {
    ClientSocket net.Conn
    OnClientDisconnected func (Client)
    OnClientPacketReceieved func (Client, Packet) []byte
    Id uint
}

func NewClient(id uint, conn net.Conn, disconnectHandler func (Client), packetHandler func (Client, Packet) []byte) Client {
    return BasicClient {
        Id: id,
        ClientSocket: conn,
        OnClientDisconnected: disconnectHandler,
        OnClientPacketReceieved: packetHandler,
    }
}

func (c BasicClient) Listen() {
    var packet [512]byte
    for {
        // will listen for message to process ending in newline (\n)
        bytes, err := c.ClientSocket.Read(packet[0:])
        if err != nil {
            c.OnClientDisconnected(c)
            c.ClientSocket.Close()
            break
        }

        response := c.OnClientPacketReceieved(c, NewRawPacket(packet[0 : ], uint(bytes)))
        if response != nil {
            _, err := c.ClientSocket.Write(response)
            if err != nil {
                c.OnClientDisconnected(c)
                c.ClientSocket.Close()
                break
            }
        }
    }
}

func (c BasicClient) GetConnection() net.Conn {
    return c.ClientSocket
}

func (c BasicClient) GetId() uint {
    return c.Id
}

func OnConduitPacketReceieved(socket Client, packet Packet) []byte {
    return []byte("READY FOR DATA")
}

func NewConduitClient(id uint, conn net.Conn, disconnectHandler func (c Client)) Client {
    return NewClient(id, conn, disconnectHandler, OnConduitPacketReceieved)
}