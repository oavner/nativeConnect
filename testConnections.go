package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sync"
  "strings"
)

func sendHttpReq(url string, method string, wg *sync.WaitGroup, success *chan bool) {
  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
  
  if strings.Contains(string(body), "EOF") { 
    *success <- false
  } else {
    *success <- true
  }

  wg.Done()
}

func main() {
  url := "http://10.155.0.113:30606/"
  method := "GET"
  maxConnections := 10000
  connectionsCounter := 0

  var wg sync.WaitGroup
  wg.Add(maxConnections)

  channel := make(chan bool)

  for conn := 0; conn < maxConnections; conn++ {
    go sendHttpReq(url, method, &wg, &channel)
    data := <- channel
    fmt.Println(data)
    if data {
      connectionsCounter++ 
    }
  }
  wg.Wait()

  fmt.Println(connectionsCounter)
}
  