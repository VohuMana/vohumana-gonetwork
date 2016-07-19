
package NetworkStack

import
(
    "fmt"
    "net"
    "log"
)

type ServerClient struct {
    Connection net.Conn
    OnDataReceived func(c *ServerClient, message []byte)
    OnServerClientConnectionClosed func(c *ServerClient, err error)
}

func NewServerClient(connection net.Conn, onData func(c *ServerClient, message []byte), onClose func(c *ServerClient, err error)) *ServerClient {
    return &ServerClient {
        Connection: connection,
        OnDataReceived: onData,
        OnServerClientConnectionClosed: onClose }
}

func (c *ServerClient) Listen() {
    var buffer [512]byte
    for {
        numBytes, err := c.Connection.Read(buffer[0:])
        if (err != nil) {
            log.Println(err)
            c.Close()
            return
        }

        c.OnDataReceived(c, buffer[0:numBytes])
    }
}

func (c *ServerClient) Close() error {
    err := c.Connection.Close()
    c.OnServerClientConnectionClosed(c, err)
    return err
}

func (c *ServerClient) Send(message []byte) error {
    _, err := c.Connection.Write(message)
    return err
}

type Server struct {
    ServerClients []*ServerClient
    Address string
    OnDataReceived func(c *ServerClient, message []byte)
    OnServerClientConnectionClosed func(c *ServerClient, err error)
    OnNewServerClient func(c *ServerClient)
}

func NewServer(address string, onData func(c *ServerClient, message []byte), onServerClientClose func(c *ServerClient, err error), onServerClient func(c *ServerClient)) *Server {
    return &Server {
        Address: address,
        OnDataReceived: onData,
        OnServerClientConnectionClosed: onServerClientClose,
        OnNewServerClient: onServerClient }
}

func (s *Server) OpenAndListenForConnections() {
    listener, err := net.Listen("tcp", s.Address)
    if (err != nil) {
        log.Fatal(err)
    }
    defer listener.Close()

    for {
        connection, err := listener.Accept()
        if (err != nil) {
            log.Println(err)
        }

        fmt.Println("New ServerClient connected")
        
        ServerClient := NewServerClient(connection, s.OnDataReceived, s.OnServerClientConnectionClosed)
        s.ServerClients = append(s.ServerClients, ServerClient)

        // Tell the ServerClient to start listening for data
        go ServerClient.Listen()
        s.OnNewServerClient(ServerClient)
    }
}

