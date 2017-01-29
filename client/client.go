package main

import "net"
import "fmt"
import "log"

func main() {

  // connect to this socket
  conn, err := net.Dial("tcp", "127.0.0.1:8080")
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  var packet [512]byte
  conn.Write([]byte("Hello"))

  for {
    bytesRead, err := conn.Read(packet[0:])
    if err != nil {
      break
    }

    fmt.Printf("Data receieved from server (%v): %v\n", bytesRead, packet[0 : bytesRead])
  }
}