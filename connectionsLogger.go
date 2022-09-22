package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "sync"
  "strings"
  "time"
  "os"
  "encoding/json"
  "strconv"
  "reflect"
)

//var connectionsReport string 

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

func testConnection(url string, method string, maxTragetSessions int, connectionDataChannel *chan connectionData, wgX *sync.WaitGroup, wgY *sync.WaitGroup, wgLogging *sync.WaitGroup){
  defer wgX.Done()

  for newTargetSessions := 0; newTargetSessions < maxTragetSessions; newTargetSessions++ {
    wgY.Add(1)
    wgLogging.Add(1)
    go sendHttpReq(url, method, wgY, connectionDataChannel)
  }
}

func log(connectionDataChannel *chan connectionData, connectionsReport *string, wgLogging *sync.WaitGroup){
	for true {
		connectionData := <- *connectionDataChannel
		fmt.Println(connectionData.ToJson())
    
		if connectionData.ConnectionState == "active" {
			fmt.Println(connectionData.Url + " : " + "success\n\n")
			*connectionsReport += connectionData.Url + " : " + "success, timeStamp: " + connectionData.TimeStamp + "\n\n"
		} else {
      fmt.Println(connectionData.Url + " : " + "failed\n\n")
			*connectionsReport += connectionData.Url + " : " + "failed, timeStamp: " + connectionData.TimeStamp + "\n\n"
		}

    wgLogging.Done()
	}
}

func jsonStringToSlice(str string) (slc []string) {
  json.Unmarshal([]byte(str), &slc)
  return
}

func main() {
  var wgX sync.WaitGroup
  var wgY sync.WaitGroup
  var wgLogging sync.WaitGroup
  var connectionsReport string
  var connectionDataChannel = make(chan connectionData)

  urlsSlice := jsonStringToSlice(os.Getenv("URLS"))
  method := os.Getenv("METHOD")
  maxTragetSessions, _ := strconv.Atoi(os.Getenv("MAX_SESSIONS")) 
  
  go log(&connectionDataChannel, &connectionsReport, &wgLogging)
  for _, url := range urlsSlice{
    wgX.Add(1)
    go testConnection(url, method, maxTragetSessions, &connectionDataChannel, &wgX, &wgY, &wgLogging)
  }
  
  wgX.Wait()
  wgY.Wait()
  wgLogging.Wait()

  fmt.Println(connectionsReport)
}
  