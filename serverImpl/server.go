package main

import (
    "fmt"
    "github.com/vohumana/vohumana-gonetwork/conduitServer"
)

func main() {
    fmt.Println("Starting conduit server...")
    server := conduitServer.NewConduitServer()

    server.StartListening()
    
    fmt.Println("Server shutting down")
}