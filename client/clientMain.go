// +build client
package main

import
(
    "github.com/vohumana/vohumana-gonetwork/NetworkStack"
    "fmt"
)

func OnData(message []byte) {
    fmt.Println("OnData")
    fmt.Println(string(message))
    fmt.Printf("Length: %v\n", len(message))
}

func OnClose(err error) {
    fmt.Println("Connection closed")
    fmt.Println(err)
}

func OnConnected() {
    fmt.Println("Client connected")
}

func main() {
    var command string
    fmt.Printf("Enter address to connect to: ")
    fmt.Scanln(&command)

    client := NetworkStack.NewClient(command, OnData, OnConnected, OnClose)
    defer client.Close()

    client.Listen()

    fmt.Printf("Press any key to close...")
    fmt.Scanln(&command)
}