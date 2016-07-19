package NetworkStack

import
(
    // "bufio"
    "net"
    "log"
)

type Client struct {
    Connection net.Conn
    Address string
    OnDataReceived func(message []byte)
    OnConnectionEstablished func()
    OnConnectionClosed func(err error)
}

func NewClient(address string, onData func(message []byte), onConnection func(), onClose func(err error)) *Client {
    return &Client {
        Address: address,
        OnDataReceived: onData,
        OnConnectionEstablished: onConnection,
        OnConnectionClosed: onClose }
}

// SendData
func (c *Client) SendData(message []byte) error {
    _, err := c.Connection.Write(message)
    return err
}

// Listen
func (c *Client) Listen() {
    conn, err := net.Dial("tcp", c.Address)
    if (err != nil) {
        log.Println(err)
        return
    }
    c.Connection = conn

    go c.listen()
}

func (c *Client) listen() {
    var buffer [512]byte
    for {
        numBytes, err := c.Connection.Read(buffer[0:])
        if (err != nil) {
            log.Println(err)
            c.Close()
            return
        }

        c.OnDataReceived(buffer[0:numBytes])
    }
}

// Close
func (c *Client) Close() error {
    err := c.Connection.Close()
    c.OnConnectionClosed(err)
    return err
}