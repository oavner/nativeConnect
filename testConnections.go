package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sync"
)

func sendHttpReq(url string, method string, wg *sync.WaitGroup) {
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

  wg.Done()
}

func main() {
  var wg sync.WaitGroup
  wg.Add(2)

  url := "http://10.155.0.113:30606/"
  method := "GET"
  maxConnections := 1000

  for conn := 0; conn < maxConnections; conn++ {
    go sendHttpReq(url, method, &wg)    
  }
}
  