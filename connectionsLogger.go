package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sync"
  "strings"
  "time"
)

var urlsSlice []string
var method string
var maxTragetSessions int
var wg sync.WaitGroup
var connectionDataChannel = make(chan connectionData)
var connectionsReport string

type connectionData struct {
  url string
  connectionState string
  err string
  timeStamp string
}

func (connectionData *connectionData) getConnectionData() (connectionDataString string){
  connectionDataString = "url: " + connectionData.url + "\nconnectionState: " + connectionData.connectionState + "\nerror: " + connectionData.err + "\ntimeStamp: " + connectionData.timeStamp + "\n\n"
  return
}

func (connectionData *connectionData) Print() {
  connectionDataString := connectionData.getConnectionData()
  fmt.Println(connectionDataString)
}

func sendHttpReq(url string, method string, wg *sync.WaitGroup, connectionDataChannel *chan connectionData) {
  client := &http.Client {
  }

  req, err := http.NewRequest(method, url, nil)
  if err != nil {
    connectionData := connectionData{url, "disconnected", err.Error(), time.Now().String()}
    *connectionDataChannel <- connectionData
    wg.Done()
    return
  }

  res, err := client.Do(req)
  if err != nil {
    connectionData := connectionData{url, "disconnected", err.Error(), time.Now().String()}
    *connectionDataChannel <- connectionData
    wg.Done()
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    connectionData := connectionData{url, "disconnected", err.Error(), time.Now().String()}
    *connectionDataChannel <- connectionData
    wg.Done()
    return
  }
  
  if strings.Contains(string(body), "EOF") || res.StatusCode > 400 {
    connectionData := connectionData{url, "disconnected", string(body) + "\n" + string(res.StatusCode), time.Now().String()}
    *connectionDataChannel <- connectionData
	  wg.Done()
    return
  }
  
  connectionData := connectionData{url, "active", "" , time.Now().String()}
  *connectionDataChannel <- connectionData
  wg.Done()
  return
}

func testConnection(url string, method string, maxTragetSessions int, connectionDataChannel *chan connectionData){
  
  for newTargetSessions := 0; newTargetSessions < maxTragetSessions; newTargetSessions++ {
	  wg.Add(1)
    go sendHttpReq(url, method, &wg, connectionDataChannel)
  }
  wg.Done()
  return 
}

func log(connectionDataChannel *chan connectionData){

	for true {
		connectionData := <- *connectionDataChannel
		connectionData.Print()
    
		if connectionData.connectionState == "active" {
			fmt.Println(connectionData.url + " : " + "success\n\n")
			connectionsReport = connectionsReport + connectionData.url + " : " + "success, timeStamp: " + connectionData.timeStamp + "\n\n"
		} else {
      fmt.Println(connectionData.url + " : " + "failed\n\n")
			connectionsReport = connectionsReport + connectionData.url + " : " + "failed, timeStamp: " + connectionData.timeStamp + "\n\n"
		}
	}
}

func printReport(wg *sync.WaitGroup){
  wg.Wait()
  fmt.Println(connectionsReport)
}

func main() {
  urlsSlice := []string{"http://10.155.0.113:22172/", "http://10.155.0.113:20667/"}
  method = "GET"
  maxTragetSessions = 100
  
  go log(&connectionDataChannel)

  for _, url := range urlsSlice{
    wg.Add(1)
    go testConnection(url, method, maxTragetSessions, &connectionDataChannel)
  }

  printReport(&wg)
}
  