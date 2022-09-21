package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sync"
  "strings"
)

var url  string
var method string
var maxConnections int
var wg sync.WaitGroup
var testStatus string

func sendHttpReq(url string, method string, wg *sync.WaitGroup, connectionDataChannel *chan string) {
  client := &http.Client {
  }

  req, err := http.NewRequest(method, url, nil)
  if err != nil {
    *connectionDataChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + err.Error() + "\n\n"
    wg.Done()
    return
  }

  res, err := client.Do(req)
  if err != nil {
    *connectionDataChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + err.Error() + "\n\n"
    wg.Done()
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    *connectionDataChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + err.Error() + "\n\n"
    wg.Done()
    return
  }
  
  if strings.Contains(string(body), "EOF") || res.StatusCode > 400 { 
    *connectionDataChannel <- "url:" + url + "\nconnectionState: disconnected\nerror: " + string(body) + "\nstatusCode: " + string(res.StatusCode) + "\n\n"
	wg.Done()
    return
  }
  
  *connectionDataChannel <- "url: " + url + "\nconnectionState: active\n\n"
  wg.Done()
  return
}

func testConnections(url string, method string, maxConnections int) (connectionsCounter int){
  connectionDataChannel := make(chan string)

  for conn := 0; conn < maxConnections; conn++ {
	wg.Add(1)
    go sendHttpReq(url, method, &wg, &connectionDataChannel)
	  go log(&connectionDataChannel)
  }
  return 
}

func log(connectionDataChannel *chan string){
	for true {
		connectionData := <- *connectionDataChannel
		fmt.Println(connectionData)
		if !strings.Contains(connectionData, "disconnected") {
			fmt.Println("Successfully created a connection\n\n")
			testStatus = "success"
		} else {
			fmt.Println("Failed to create a connection\n\n")
			testStatus = "failed"
		}
	}
}

func main() {
  url = "http://10.155.0.113:22172/"
  method = "GET"
  maxConnections = 10
  testConnections(url, method, maxConnections)
  wg.Wait()

  fmt.Println("testStatus: " + testStatus)
}
  