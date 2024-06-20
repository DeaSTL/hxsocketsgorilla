package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/deastl/hx-sockets/hx"
)

type data struct {
  State bool `json:"state"`
}



func main(){
  mux := http.NewServeMux()
  server := hx.NewServer(mux)

  server.Mount("/ws")

  server.Listen("some_message",func (client *hx.Client, msg []byte){
    d := data{} 
    err := json.Unmarshal(msg,&d)

    if err != nil {
      log.Printf("Error: %v", err)
    }

    d.State = !d.State

    stateStr,err := json.Marshal(d)

    if err != nil {
      log.Printf("err: %v", err)
    }
    
    strToSend := fmt.Sprintf(`<button hx-vals='%v' id="some_message" hx-trigger="click" ws-send>%t</button>`,string(stateStr),d.State)
    

    client.SendStr(strToSend)
    log.Printf("client: %v", client.ID) 
  })

  mux.Handle("/",http.HandlerFunc(func (res http.ResponseWriter,req *http.Request){
    res.Write([]byte(`
      <html>
        <head>
          <script src="https://unpkg.com/htmx.org@1.9.12"></script>
          <script src="https://unpkg.com/htmx.org@1.9.12/dist/ext/ws.js"></script>
        </head>
        <body hx-ext="ws" ws-connect="/ws">
          <button id="some_message" hx-trigger="load" ws-send>false</button>
        </body>
      </html>
    `))
  }))


  http.ListenAndServe(":3000",mux)


}
