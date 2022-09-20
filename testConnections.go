package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sync"
  "strings"
)

func sendHttpReq(url string, method string, wg *sync.WaitGroup, connectionStateChannel *chan string) {
  client := &http.Client {
  }

  req, err := http.NewRequest(method, url, nil)
  if err != nil {
    fmt.Println(err)
    *connectionStateChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + err.Error() + "\n\n"
    wg.Done()
    return
  }

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    *connectionStateChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + err.Error() + "\n\n"
    wg.Done()
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    *connectionStateChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + err.Error() + "\n\n"
    wg.Done()
    return
  }
  
  if strings.Contains(string(body), "EOF") || res.StatusCode > 400 { 
    *connectionStateChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + string(body) + "\n\n"
  } else {
    *connectionStateChannel <- "url: " + url + "\nconnectionState: connected\n\n"
    wg.Done()
    return
  }

  wg.Done()
  return
}

func testConnections(url string, method string, maxConnections int) (connectionsCounter int){
  connectionsCounter = 0

  var wg sync.WaitGroup
  wg.Add(maxConnections)

  connectionStateChannel := make(chan string)

  for conn := 0; conn < maxConnections; conn++ {
    go sendHttpReq(url, method, &wg, &connectionStateChannel)
    connectionState := <- connectionStateChannel
    fmt.Println(connectionState)
    if !strings.Contains(connectionState, "disconnected") {
      connectionsCounter++ 
    } else {
      fmt.Println("Failed to create new connections stopping session creations")
      break
    }
  }
  return 
}

var url  string
var method string
var maxConnections int

func main() {
  url = "http://10.155.0.113:22172/"
  method = "GET"
  maxConnections = 10
  
  if testConnections(url, method, maxConnections) == maxConnections {
    fmt.Println("Successfully reached the max ammount of connections")
  } else {
    fmt.Println("Failed to reached the max ammount of connections. Check the SDN state")
  }
}
  