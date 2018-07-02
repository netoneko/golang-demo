package main

import (
  "net"
  "os"
  "fmt"
  "time"
  "encoding/json"
  "io/ioutil"
  "strings"
)

func handleClient(conn net.Conn) {
   conn.SetReadDeadline(time.Now().Add(30 * time.Second))
   
   defer conn.Close()

   fmt.Println("Reading request")
   request, err := ioutil.ReadAll(conn)

   fmt.Println("Message received", string(request))

   if err != nil {
     fmt.Println(err)
   }

   response := map[string]string{
     "timestamp": time.Now().Format(time.RFC3339),
     "node": os.Args[2],
   }

   jsonResponse, err := json.Marshal(response)
   conn.Write(jsonResponse)
}

func pingPeer(peer string) (string, error) {
  peerAddr, err := net.ResolveTCPAddr("tcp", peer)

  if err != nil {
    return "", err
  }

  conn, err := net.DialTCP("tcp", nil, peerAddr)

  if err != nil {
    return "", err
  } else {
    conn.Write([]byte("{\"ping\":\"" + os.Args[2] + "\"}"));
    conn.CloseWrite()
  }

  result, err := ioutil.ReadAll(conn)

  if err != nil {
    return "", err
  }

  return string(result), nil
}

func pingPeerContinuously(peer string) {
  for {
     fmt.Println("Trying to ping", peer);

     result, err := pingPeer(peer)
     fmt.Println("Ping result", result, err)

     time.Sleep(time.Second)
   }
}

func main() {
  fmt.Println(os.Args)

  for _, peer := range strings.Split(os.Args[3], ",") {
    go pingPeerContinuously(peer)
  }

  listener, err := net.Listen("tcp", os.Args[1])

  if err == nil {
    fmt.Println("Listening on port", os.Args[1])
    
    for {
      conn, err := listener.Accept()

      if err != nil {
        fmt.Println("ERROR OCCURED", err)
        conn.Close()
      }

      go handleClient(conn)
    }
  }
}
