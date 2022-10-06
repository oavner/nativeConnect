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

// Inhereting the error interface
type error interface {
  Error() string
}

// Creating custom httpError struct 
type httpError struct {
  ResBody string
  ResStatusCode int
}

// Implementing the Error() func of the httpError struct. raising connectionErrors for http packets.
// Currently returning only the status code of the request without the body.
func (e *httpError) Error() string {
  return string(e.ResStatusCode)
}

// Struct that holds the metadata about the connection. used for logging connection state and errors.
type connectionData struct {
  Source string `json:"source"`
  Url string `json:"url"`
  ConnectionState string `json:"connection_state"`
  Err string `json:"error"`
  TimeStamp string `json:"timestamp"`
}

// Converting connectionData sturct to Json.
func (connectionData *connectionData) ToJson() string {
  connectionDataJson, err := json.Marshal(connectionData)
  if err != nil {
    fmt.Println(err)
  }
  return string(connectionDataJson)
}

// Sends connectionData to connectionDataStream and sends Done() call to a given waitGroup.
func handleConnectionError(sourceIP string, url string, err error, wgY *sync.WaitGroup, connectionDataChannel *chan connectionData){
  connectionData := connectionData{sourceIP, url, "disconnected", err.Error(), time.Now().String()}
  *connectionDataChannel <- connectionData
  wgY.Done()
}

func getSourceIP(r *http.Request) string {
  IPAddress := r.Header.Get("X-Real-Ip")
  if IPAddress == "" {
      IPAddress = r.Header.Get("X-Forwarded-For")
  }
  if IPAddress == "" {
      IPAddress = r.RemoteAddr
  }
  return IPAddress
}

// Sends http req to a given url in a certain method and raises h handleConnectionError for any type of error 
// That may raise along the way.
func sendHttpReq(url string, method string, wgY *sync.WaitGroup, connectionDataChannel *chan connectionData) {
  client := &http.Client {
  }

  req, err := http.NewRequest(method, url, nil)
  sourceIP := getSourceIP(req)
  if err != nil {
    handleConnectionError(sourceIP, url, err, wgY, connectionDataChannel)
    return
  }

  res, err := client.Do(req)
  if err != nil {
    handleConnectionError(sourceIP, url, err, wgY, connectionDataChannel)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    handleConnectionError(sourceIP, url, err, wgY, connectionDataChannel)
    return
  }
  
  if strings.Contains(string(body), "EOF") || res.StatusCode > 400 {
    handleConnectionError(sourceIP, url, &httpError{string(body), res.StatusCode}, wgY, connectionDataChannel)
    return
  }
  
  connectionData := connectionData{sourceIP, url, "active", "" , time.Now().String()}
  *connectionDataChannel <- connectionData
  wgY.Done()
  return
}

// Rises concurrent sendHttpReq functions as the ammount of maxTargetSessions.
// Rises Add() function for the connection waitGroup and the Logging waitGroup in each iteration.
func testConnection(url string, method string, maxTragetSessions int, connectionDataChannel *chan connectionData, wgX *sync.WaitGroup, wgY *sync.WaitGroup, wgLogging *sync.WaitGroup){
  defer wgX.Done()

  for newTargetSessions := 0; newTargetSessions < maxTragetSessions; newTargetSessions++ {
    wgY.Add(1)
    wgLogging.Add(1)
    go sendHttpReq(url, method, wgY, connectionDataChannel)
  }
}

// Logger function that reads the connectionDatas from the connectionDataStream and logs them to stdout.
// The logger also keeps a connectionReport that sums up all conection states.
// This func will be called logger in the future and will use the log module for better logging logic.
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

// Used to create a slice of urls from a string typed enviorment variable.
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