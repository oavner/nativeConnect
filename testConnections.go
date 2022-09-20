package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://10.155.0.113:30606/"
  method := "GET"
  connectionLimit := 10000

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }

  for connection := 0; connection < connectionLimit; connection++ {
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
    }
  }