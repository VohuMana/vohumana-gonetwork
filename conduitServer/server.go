package conduitServer

import (
    "fmt"
    "net"
    "log"
)

type baseServer struct {
    listenerSocket net.Listener
    clients map[uint]*client
    createClient func (uint, net.Conn) client
}

func newBaseServer() *baseServer {
    server := &baseServer {
        clients: make(map[uint]*client),
    }
    server.createClient = server.baseClientCreation
    return server
}

func (s *baseServer) baseClientCreation(id uint, conn net.Conn) client {
     return newClient(id, conn, s.clientDisconnected, s.clientPacketReceieved)
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

        go client.listen()

        s.clients[newClientId] = &client
    }
}

func (s *baseServer) clientDisconnected(c client) {
    fmt.Printf("Client (%v) with address %s disconnected\n", c.getId(), c.getConnection().RemoteAddr().String())
    s.clients[c.getId()] = nil
}

func (s *baseServer) clientPacketReceieved(c client, pkt packet) []byte {
    fmt.Printf("Message Received (%s): %s\n", c.getConnection().RemoteAddr().String(), string(pkt.getData()[0 : pkt.getSize()]))

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

func (s *ConduitServer) createConduit(id uint, conn net.Conn) client {
    return newConduitClient(id, conn, s.server.clientDisconnected)
}

func (s *ConduitServer) StartListening() {
    s.server.startListening()
}