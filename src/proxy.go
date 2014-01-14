package main

import (
  "flag"
  "fmt"
  "net"
  "proxy"
)

var gPort int
var gPortStr string

func init() {
  flag.IntVar(&gPort, "port", 3100, "port to listen on")
  gPortStr = fmt.Sprintf(":%d",gPort)
}

func main() {
  flag.Parse()
  fmt.Println("listening on port ",gPortStr)
  ln, err := net.Listen("tcp",gPortStr)
  if err!=nil {
    // handle error
  }
  for {
    conn,err := ln.Accept()
    if err!=nil {
      // handle error
      continue
    }
    fmt.Println("mainloop received: ",conn.LocalAddr(), conn.RemoteAddr())
    go proxy.Handle(conn)
  }
}
