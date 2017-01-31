package conduitServer

import (
    "net"
)

type client interface {
    listen()
    getConnection() net.Conn
    getId() uint
}

type basicClient struct {
    clientSocket net.Conn
    onClientDisconnected func (client)
    onClientPacketReceieved func (client, packet) []byte
    id uint
}

func newClient(id uint, conn net.Conn, disconnectHandler func (client), packetHandler func (client, packet) []byte) client {
    return basicClient {
        id: id,
        clientSocket: conn,
        onClientDisconnected: disconnectHandler,
        onClientPacketReceieved: packetHandler,
    }
}

func (c basicClient) listen() {
    defer c.clientSocket.Close()

    var packet [512]byte
    for {
        // will listen for message to process ending in newline (\n)
        bytes, err := c.clientSocket.Read(packet[0:])
        if err != nil {
            c.onClientDisconnected(c)
            break
        }

        response := c.onClientPacketReceieved(c, newRawPacket(packet[0 : ], uint(bytes)))
        if response != nil {
            _, err := c.clientSocket.Write(response)
            if err != nil {
                c.onClientDisconnected(c)
                break
            }
        }
    }
}

func (c basicClient) getConnection() net.Conn {
    return c.clientSocket
}

func (c basicClient) getId() uint {
    return c.id
}

func onConduitPacketReceieved(socket client, packet packet) []byte {
    return []byte("READY FOR DATA")
}

func newConduitClient(id uint, conn net.Conn, disconnectHandler func (c client)) client {
    return newClient(id, conn, disconnectHandler, onConduitPacketReceieved)
}