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
)

//inhereting the error interface
type error interface {
  Error() string
}

//creating custom httpError struct 
type httpError struct {
  ResBody string
  ResStatusCode int
}

//implementing the Error() func of the httpError struct. raising connectionErrors for http packets.
// currently returning only the status code of the request without the body.
func (e *httpError) Error() string {
  return string(e.ResStatusCode)
}

//struct that holds the metadata about the connection. used for logging connection state and errors.
type connectionData struct {
  Url string `json:"url"`
  ConnectionState string `json:"connection_state"`
  Err string `json:"error"`
  TimeStamp string `json:"timestamp"`
}

//converting connectionData sturct to Json.
func (connectionData *connectionData) ToJson() string {
  connectionDataJson, err := json.Marshal(connectionData)
  if err != nil {
    fmt.Println(err)
  }
  return string(connectionDataJson)
}

//sends connectionData to connectionDataStream and sends Done() call to a given waitGroup.
func handleConnectionError(url string, err error, wgY *sync.WaitGroup, connectionDataChannel *chan connectionData){
  connectionData := connectionData{url, "disconnected", err.Error(), time.Now().String()}
  *connectionDataChannel <- connectionData
  wgY.Done()
}

//sends http req to a given url in a certain method and raises h handleConnectionError for any type of error 
//that may raise along the way.
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

// rises concurrent sendHttpReq functions as the ammount of maxTargetSessions.
// rises Add() function for the connection waitGroup and the Logging waitGroup in each iteration.
func testConnection(url string, method string, maxTragetSessions int, connectionDataChannel *chan connectionData, wgX *sync.WaitGroup, wgY *sync.WaitGroup, wgLogging *sync.WaitGroup){
  defer wgX.Done()

  for newTargetSessions := 0; newTargetSessions < maxTragetSessions; newTargetSessions++ {
    wgY.Add(1)
    wgLogging.Add(1)
    go sendHttpReq(url, method, wgY, connectionDataChannel)
  }
}

// logger function that reads the connectionDatas from the connectionDataStream and logs them to stdout.
// the logger also keeps a connectionReport that sums up all conection states.
// this func will be called logger in the future and will use the log module for better logging logic.
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

// used to create a slice of urls from a string typed enviorment variable.
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
  