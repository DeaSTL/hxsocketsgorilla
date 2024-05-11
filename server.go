package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func GenB64(length int) string {
	dembytes := make([]byte, length)
	_, err := rand.Read(dembytes)
	if err != nil {
		return ""
	}
	encoded := base64.URLEncoding.EncodeToString(dembytes)
	return encoded
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Printf("Error status: %d", status)
	},
}

type HXMessage struct {
	Request     string `json:"HX-Request"`
	Trigger     string `json:"HX-Trigger"`
	TriggerName string `json:"HX-Trigger-Name"`
	Target      string `json:"HX-Target"`
	CurrentURL  string `json:"HX-Current-URL"`
  Includes    map[string]string
}

type HXListener struct {
	Callback func(*HXConnectionCtx, *HXMessage)
}

type HXServer struct {
	Connections  *[]*HXConnectionCtx
	listeners    map[string]HXListener
	OnConnection func(*HXConnectionCtx)
	OnDisconnect func(*HXConnectionCtx)
}

type HXConnectionCtx struct {
	Client    *websocket.Conn
	SessionID string
}

// type Message struct {
// 	Type string          `json:"type"`
// 	Data HtmxMessage `json:"data"`
// }

func (ss *HXServer) LogConnections() {
	for _, client := range *ss.Connections {
		log.Printf("Client %v", client.SessionID)
	}
}

func (ss *HXServer) New() {
	ss.OnConnection = func(ctx *HXConnectionCtx) {}
	ss.OnDisconnect = func(ctx *HXConnectionCtx) {}
	ss.listeners = map[string]HXListener{}
	ss.Connections = &[]*HXConnectionCtx{}
}

func (ss *HXServer) Start(mux *http.ServeMux, endpoint string) {
	mux.Handle(endpoint, http.HandlerFunc(ss.handleNewConnection))
}

func (ss *HXServer) handleCloseConnection(ctx *HXConnectionCtx) func(int, string) error {

	return func(code int, text string) error {
		ss.OnDisconnect(ctx)
		var index int

		for c_index, conn := range *ss.Connections {
			if conn.SessionID == ctx.SessionID {
				index = c_index
			}
		}

		*ss.Connections = append((*ss.Connections)[:index], (*ss.Connections)[index+1:]...)

		return nil
	}
}
func (ss *HXServer) handleNewConnection(w http.ResponseWriter, r *http.Request) {
	log.Print("Attempting to connect new client")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Could not upgrade connection: %v", err)
	}

	new_ctx := ss.newConnection(conn)

	conn.SetCloseHandler(ss.handleCloseConnection(&new_ctx))

	ss.OnConnection(&new_ctx)

	ss.newMessageListener(&new_ctx)
}
/*
{
  HEADERS: {}, // HTMX stuff
  some_key: 'some value'
  ...
}
*/
func (ss *HXServer) newMessageListener(ctx *HXConnectionCtx) {
	go func(ctx *HXConnectionCtx) {
		for {
      rawMessage := map[string]any{}

			err := ctx.Client.ReadJSON(&rawMessage)

      

			if err != nil {
				log.Printf("Error reading json message: %v", err)
				break
			}
      
      /*
      We're pulling off the HEADERS and making it a HtmxMesssage
      */
      jsonMessage, err := json.Marshal(rawMessage["HEADERS"])

      if err != nil {

      }

      message := HXMessage{}

      err = json.Unmarshal(jsonMessage,&message)


      // Yoinkin' that header off the json blob
      delete(rawMessage,"HEADERS")
      

      message.Includes = make(map[string]string)

      for key, value := range rawMessage {

        switch v := value.(type) {
          case string:
            message.Includes[key] = v
          case int:
            message.Includes[key] = strconv.Itoa(v)
          case float64:
            message.Includes[key] = strconv.FormatFloat(v,'f',-1,64)
          default:
            message.Includes[key] = fmt.Sprintf("%v",v)

        }

      }

			err = ss.messageHandler(ctx, &message)

			if err != nil {
				log.Printf("Error handling message: %v ", err)
				continue
			}
		}
		defer ctx.Client.Close()
	}(ctx)
}

func (ss *HXServer) messageHandler(ctx *HXConnectionCtx, message *HXMessage) error {
	for event, listener := range ss.listeners {
		if message.Trigger == event {

			listener.Callback(ctx, message)
		}
	}
	return nil
}

func (ss *HXServer) newConnection(conn *websocket.Conn) HXConnectionCtx {
	new_ctx := HXConnectionCtx{
		SessionID: GenB64(32),
		Client:    conn,
	}

	*ss.Connections = append(*ss.Connections, &new_ctx)

	return new_ctx
}

func (ss *HXServer) Broadcast(event string, message []byte) error {
	for _, conn := range *ss.Connections {
		err := conn.Send(message)

		if err != nil {
			log.Printf("Error broadcasting: %v", err)
		}

		log.Printf("Broadcasting message: %v to: %v", message, conn.SessionID)
	}
	return nil
}

func (ss *HXServer) Listen(event string, listener func(*HXConnectionCtx, *HXMessage)) {
	ss.listeners[event] = HXListener{Callback: listener}
}

func (ctx *HXConnectionCtx) Send(message []byte) error {

  err := ctx.Client.WriteMessage(1,message)

	if err != nil {
		return fmt.Errorf("Could not send message err: %v ", err)
	}

	return nil
}

func (ctx *HXConnectionCtx) SendStr(message string) error {
  return ctx.Send([]byte(message))
}

func (ss *HXServer) SendFilter(event string, message []byte, check func(*HXConnectionCtx) bool) {
	for _, conn := range *ss.Connections {
		if check(conn) {
			err := conn.Send(message)

			if err != nil {
				log.Printf("Error broadcasting: %v", err)
			}

			log.Printf("Sending filtered message: %v to: %v", message, conn.SessionID)
		}
	}
}

