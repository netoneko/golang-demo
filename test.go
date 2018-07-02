package main

import(
  "fmt"
  "time"
  "strconv"
)

type Payload struct {
  timestamp time.Time
  message string
}

func generateMessages(messages chan Payload, f func(i int) string) {
  for i:=0 ; true ; i++ {
    messages <- Payload{time.Now(), f(i)}
    time.Sleep(1 * 1e9)
  }
}

func listenToMessages(messages chan Payload) {
  for msg := range messages {
    fmt.Println("Received message:", msg.timestamp, msg.message)
  }
}

func main() {
  const test = "Hello"
  var timestamp = time.Now()

  var messages = make(chan Payload)

  fmt.Println(test, timestamp)

  go generateMessages(messages, func(i int) string {
    return "Hello " + strconv.Itoa(i)
  })

  time.Sleep(0.5 * 1e9)

  go generateMessages(messages, func(i int) string {
    return "World " + strconv.Itoa(i)
  })

  go listenToMessages(messages)

  time.Sleep(3 * 1e9)

  select {

  }
}
