package conduitServer

import (
    "fmt"
    "net"
    "log"
)

type baseServer struct {
    listenerSocket net.Listener
    clients map[uint]*Client
    createClient func (uint, net.Conn) Client
}

func newBaseServer() *baseServer {
    server := &baseServer {
        clients: make(map[uint]*Client),
    }
    server.createClient = server.baseClientCreation
    return server
}

func (s *baseServer) baseClientCreation(id uint, conn net.Conn) Client {
     return NewClient(id, conn, s.clientDisconnected, s.clientPacketReceieved)
}

func (s *baseServer) startListening() {
    // Start listener
    ln, err := net.Listen("tcp", ":8080")
    if nil != err {
        log.Fatal(err)
    }

    s.listenerSocket = ln
    defer s.listenerSocket.Close()

    for {
        // Accept incoming connection
        conn, err := ln.Accept()
        if nil != err {
            log.Fatal(err)
        }

        newClientId := uint(len(s.clients) + 1)

        client := s.createClient(newClientId, conn)

        go client.Listen()

        s.clients[newClientId] = &client
    }
}

func (s *baseServer) clientDisconnected(c Client) {
    fmt.Printf("Client (%v) with address %s disconnected\n", c.GetId(), c.GetConnection().RemoteAddr().String())
    s.clients[c.GetId()] = nil
}

func (s *baseServer) clientPacketReceieved(c Client, packet Packet) []byte {
    fmt.Printf("Message Received (%s): %s\n", c.GetConnection().RemoteAddr().String(), string(packet.GetData()[0 : packet.GetSize()]))

    return nil
}

type ConduitServer struct {
    server *baseServer
}

func NewConduitServer() *ConduitServer {
    conduitServer := &ConduitServer {
        server: newBaseServer(),
    }
    conduitServer.server.createClient = conduitServer.createConduit
    return conduitServer
}

func (s *ConduitServer) createConduit(id uint, conn net.Conn) Client {
    return NewConduitClient(id, conn, s.server.clientDisconnected)
}

func (s *ConduitServer) StartListening() {
    s.server.startListening()
}