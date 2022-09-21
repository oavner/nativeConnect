package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sync"
  "strings"
  "time"
  "encoding/json"
)

var urlsSlice []string
var method string
var maxTragetSessions int
var wgX sync.WaitGroup
var wgY sync.WaitGroup
var wgLogging sync.WaitGroup
var connectionDataChannel = make(chan connectionData)
var connectionsReport string

type error interface {
  Error() string
}

type httpError struct {
  ResBody string
  ResStatusCode int
}

func (e *httpError) Error() string {
  return string(e.ResStatusCode)
}

type connectionData struct {
  Url string `json:"url"`
  ConnectionState string `json:"connection_state"`
  Err string `json:"error"`
  TimeStamp string `json:"timestamp"`
}

func (connectionData *connectionData) ToJson() string {
  connectionDataJson, err := json.Marshal(connectionData)
  if err != nil {
    fmt.Println(err)
  }
  return string(connectionDataJson)
}

func handleConnectionError(url string, err error, wgY *sync.WaitGroup, connectionDataChannel *chan connectionData){
  connectionData := connectionData{url, "disconnected", err.Error(), time.Now().String()}
  *connectionDataChannel <- connectionData
  wgY.Done()
}

func sendHttpReq(url string, method string, wgY *sync.WaitGroup, connectionDataChannel *chan connectionData) {
  client := &http.Client {
  }

  req, err := http.NewRequest(method, url, nil)
  if err != nil {
    handleConnectionError(url, err, wgY, connectionDataChannel)
    return
  }

  res, err := client.Do(req)
  if err != nil {
    handleConnectionError(url, err, wgY, connectionDataChannel)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    handleConnectionError(url, err, wgY, connectionDataChannel)
    return
  }
  
  if strings.Contains(string(body), "EOF") || res.StatusCode > 400 {
    handleConnectionError(url, &httpError{string(body), res.StatusCode}, wgY, connectionDataChannel)
    return
  }
  
  connectionData := connectionData{url, "active", "" , time.Now().String()}
  *connectionDataChannel <- connectionData
  wgY.Done()
  return
}

func testConnection(url string, method string, maxTragetSessions int, connectionDataChannel *chan connectionData){
  defer wgX.Done()

  for newTargetSessions := 0; newTargetSessions < maxTragetSessions; newTargetSessions++ {
    wgY.Add(1)
    wgLogging.Add(1)
    go sendHttpReq(url, method, &wgY, connectionDataChannel)
  }
}

func log(connectionDataChannel *chan connectionData){

	for true {
		connectionData := <- *connectionDataChannel
		fmt.Println(connectionData.ToJson())
    
		if connectionData.ConnectionState == "active" {
			fmt.Println(connectionData.Url + " : " + "success\n\n")
			connectionsReport = connectionsReport + connectionData.Url + " : " + "success, timeStamp: " + connectionData.TimeStamp + "\n\n"
		} else {
      fmt.Println(connectionData.Url + " : " + "failed\n\n")
			connectionsReport = connectionsReport + connectionData.Url + " : " + "failed, timeStamp: " + connectionData.TimeStamp + "\n\n"
		}

    wgLogging.Done()
	}
}

func main() {
  urlsSlice = []string{"http://10.155.0.113:22172/", "http://10.155.0.113:20667/"}
  method = "GET"
  maxTragetSessions = 4
  
  go log(&connectionDataChannel)
  for _, url := range urlsSlice{
    wgX.Add(1)
    go testConnection(url, method, maxTragetSessions, &connectionDataChannel)
  }
  
  wgX.Wait()
  wgY.Wait()
  wgLogging.Wait()

  fmt.Println(connectionsReport)
}
  