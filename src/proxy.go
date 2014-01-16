package main

import (
  "flag"
  "fmt"
  "net"
  "proxy"
  "log"
  "log/syslog"
)

var gPort int
var gPortStr string

func init() {
  flag.IntVar(&gPort, "port", 3100, "port to listen on")
  gPortStr = fmt.Sprintf(":%d",gPort)
}

func main() {
  // the lazy man solution: set the rsyslog instance as the log.std object so it
  // can be shared across all files. Another option is to follow the code under
  // http://golang.org/src/pkg/log/log.go to create a synched singleton rsyslog
  // writer
  slog, err := syslog.New(syslog.LOG_INFO, "goxy")
  log.SetOutput(slog)
  log.SetFlags(log.Lshortfile)
  log.Println("hello from goxy")
  
  flag.Parse()
  log.Println("listening on port ",gPortStr)
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
    log.Println("mainloop received: ",conn.LocalAddr(), conn.RemoteAddr())
    go proxy.Handle(conn)
  }
}
