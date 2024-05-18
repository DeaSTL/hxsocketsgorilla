package compat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/deastl/hx-sockets"
	"github.com/deastl/hx-sockets/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  250000,
	WriteBufferSize: 250000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Printf("Error status: %d", status)
	},
}

type NethttpClient struct {
	SessionID string
	Conn      *websocket.Conn
}

type NethttpListener struct {
	Callback func(*NethttpClient, *hx.Message)
}
type NethttpServer struct {
	mux          *http.ServeMux
	OnConnection func(ctx *NethttpClient)
	OnDisconnect func(ctx *NethttpClient)
	listeners    map[string]NethttpListener
	Connections  *[]*NethttpClient
}

func NewNetHttp(mux *http.ServeMux) hx.IServer {
	return NethttpServer{
    mux: mux,
    OnConnection:func(ctx *NethttpClient) {},
    OnDisconnect : func(ctx *NethttpClient) {},
    listeners : map[string]NethttpListener{},
    Connections : &[]*NethttpClient{},
  }
}

// Broadcast implements hx.IServer.
func (s NethttpServer) Broadcast(event string, message []byte) error {
  return nil
}

// Listen implements hx.IServer.
func (s NethttpServer) Listen(event string, listener func(*NethttpClient, *hx.Message)) {
  s.listeners[event] = NethttpListener{Callback: listener}
}

// Start implements hx.IServer.
func (s NethttpServer) Start(mountpoint string) {
	s.mux.Handle(mountpoint, http.HandlerFunc(s.handleNewConnection))
}

func (s NethttpServer) newConnection(conn *websocket.Conn) NethttpClient {
	new_ctx := NethttpClient{
		SessionID: utils.GenB64(32),
		Conn:      conn,
	}

	*s.Connections = append(*s.Connections, &new_ctx)

	return new_ctx
}

func (s *NethttpServer) handleCloseConnection(client *NethttpClient) func(int, string) error {

	return func(code int, text string) error {
		s.OnDisconnect(client)
		var index int

		for c_index, conn := range *s.Connections {
			if conn.SessionID == client.SessionID {
				index = c_index
			}
		}

		*s.Connections = append((*s.Connections)[:index], (*s.Connections)[index+1:]...)

		return nil
	}
}

func (s *NethttpServer) newMessageListener(client *NethttpClient) {
	go func(client *NethttpClient) {
		for {
			rawMessage := map[string]any{}

			err := client.Conn.ReadJSON(&rawMessage)

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

			message := hx.Message{}

			err = json.Unmarshal(jsonMessage, &message)

			// Yoinkin' that header off the json blob
			delete(rawMessage, "HEADERS")

			message.Includes = rawMessage

			err = s.messageHandler(client, &message)

			if err != nil {
				log.Printf("Error handling message: %v ", err)
				continue
			}
		}
		defer client.Conn.Close()
	}(client)
}

func (s *NethttpServer) messageHandler(client *NethttpClient, message *hx.Message) error {
  log.Printf("client: %+v", *client)
  log.Printf("message: %+v", *message)
	for event, listener := range s.listeners {
		if message.Trigger == event {

			listener.Callback(client, message)
		}
	}
	return nil
}

func (s *NethttpServer) handleNewConnection(w http.ResponseWriter, r *http.Request) {
	log.Print("Attempting to connect new client")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Could not upgrade connection: %v", err)
	}

	new_client := s.newConnection(conn)

	conn.SetCloseHandler(s.handleCloseConnection(&new_client))

	s.OnConnection(&new_client)

	s.newMessageListener(&new_client)
}

func (ctx *NethttpClient) Send(message []byte) error {

	err := ctx.Conn.WriteMessage(1, message)

	if err != nil {
		return fmt.Errorf("Could not send message err: %v ", err)
	}

	return nil
}

func (ctx *NethttpClient) SendStr(message string) error {
	return ctx.Send([]byte(message))
}
