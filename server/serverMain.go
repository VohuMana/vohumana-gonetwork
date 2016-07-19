// +build !client
package main

import
(
    "github.com/vohumana/vohumana-gonetwork/NetworkStack"
    "fmt"
)

func OnServerData(c *NetworkStack.ServerClient, message []byte) {
    if (len(message) == 0) {
        return
    }

    fmt.Println("OnData")
    fmt.Printf("%v\n", c.Connection.LocalAddr())
    fmt.Println(string(message))
}

func OnServerClose(c *NetworkStack.ServerClient, err error) {
    fmt.Println("Connection closed")
    fmt.Printf("%v\n", c.Connection.LocalAddr())
    fmt.Println(err)
}

func OnServerConnected(c *NetworkStack.ServerClient) {
    fmt.Println("Client connected")
    fmt.Printf("%v\n", c.Connection.LocalAddr())
    c.Connection.Write([]byte("Hello World!"))
}

func main() {
    server := NetworkStack.NewServer("localhost:9999", OnServerData, OnServerClose, OnServerConnected)

    server.OpenAndListenForConnections()
    
    fmt.Printf("Press any key to close...")
    fmt.Scanln()
}