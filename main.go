package main

import (
	"log"
	"net/http"
)



func main() {

  mux := http.NewServeMux()

  server := HXServer{}
  server.New()

  server.OnConnection = func(ctx *HXConnectionCtx){
    log.Printf("New connection %+v", ctx.SessionID)
  }


  server.Listen("some_message",func(ctx *HXConnectionCtx,msg *HXMessage){
    log.Printf("htmx message: %+v", msg.Includes["c"])
    ctx.SendStr(`
    <div id="some_message">
      Godeem
    </div>
    `)
  })

  mux.Handle("/",http.FileServer(http.Dir("./static/")))
  server.Start(mux,"/ws")
  
  log.Printf("Listening for requests")
  http.ListenAndServe(":6900",mux)
}
