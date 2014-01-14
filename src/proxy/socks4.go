package proxy

import (
  "fmt"
  "net"
)

const (
  VER = 0
  COMMAND = 1
  PORT = 2
  IPADDR = 4

  SUCCESS = 0x5A
  FAIL = 0x5B
)

type SocksHandler interface {
  Handle(conn net.Conn, handshake []byte)
}

type Socks4Handler struct {
}

func Handle(conn net.Conn) {
  var handler SocksHandler
  b := make([]byte,1024,1024)
  n,err := conn.Read(b)
  if(err!=nil) {
    // handle error
    conn.Close()
    return
  } else {
    switch ver := b[VER]; ver {
    case 0x4:
      handler = Socks4Handler{}
    case 0x5:
      fallthrough
    default:
      resp := []byte{0x4,FAIL,0,0,0,0,0,0}
      conn.Write(resp)
      conn.Close()
      return
    }
  }
  
  if(handler!=nil) {
    handler.Handle(conn, b[:n])
  }
}

func pipe(src net.Conn, tgt net.Conn) {
  b := make([]byte,1024,1024)
  n,err := src.Read(b)
  for ;err==nil; {
    fmt.Println(src.RemoteAddr()," received ",n)
    tgt.Write(b[:n])
    n,err = src.Read(b)
  }
  fmt.Println("Done with ",src.RemoteAddr(),err)
  src.Close()
  tgt.Close()
}

func (h Socks4Handler) Handle(conn net.Conn, handshake []byte) {
  fmt.Println("socks4 handle, handshake = ",handshake)

  // TODO: get the correct port from the request
  // TODO 2: SSL?
  // TODO 3: support 4a by doing DNS resolve on host
  ip_str := fmt.Sprintf("%d.%d.%d.%d:%d",
                handshake[IPADDR],handshake[IPADDR+1],handshake[IPADDR+2],handshake[IPADDR+3],
                uint(handshake[PORT])*0x100+uint(handshake[PORT+1]))
  fmt.Println("ipaddr = ",ip_str)
  resp := []byte{0,SUCCESS,0,0,0,0,0,0}
  conn.Write(resp)
  conn2,err := net.Dial("tcp",ip_str)
  fmt.Println(conn2.RemoteAddr(),err)
  go pipe(conn,conn2)
  go pipe(conn2,conn)
}

