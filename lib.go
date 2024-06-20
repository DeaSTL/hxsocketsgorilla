package hx

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ListenerFunc func(*Client, []byte)
type ClientConnectFunc func(*Client)
type ClientDisconnectFunc func(*Client)

type Server struct {
	mux                *http.ServeMux
	clients            map[string]*Client
	listeners          map[string]ListenerFunc
	OnClientConnect    ClientConnectFunc
	OnClientDisconnect ClientDisconnectFunc
	mtex               sync.Mutex
}

func NewServer(mux *http.ServeMux) Server {
	return Server{
		mux:                mux,
		clients:            map[string]*Client{},
		listeners:          map[string]ListenerFunc{},
		OnClientConnect:    func(*Client) {},
		OnClientDisconnect: func(*Client) {},
		mtex:               sync.Mutex{},
	}
}

func genB64(length int) string {
	dembytes := make([]byte, length)
	_, err := rand.Read(dembytes)
	if err != nil {
		return ""
	}
	encoded := base64.URLEncoding.EncodeToString(dembytes)
	return encoded
}

type Client struct {
	Conn *websocket.Conn
	ID   string
}

// WriteMessage wraps the underlying websocket WriteMessage() function for convenience.
// Example:
//
//	<div id="event-name">some text</div>
//
// This will receive this message and handle it based on the state of hx-swap
func (c *Client) WriteMessage(msg []byte) error {
	return c.Conn.WriteMessage(1, msg)
}

// HXHeaders is part of the [HXWSHeadersWrapper]
// This can be used in conjuction with yor message struct (under the json key of "HEADERS") if these attributes are needed
type HXHeaders struct {
	HXRequest     string  `json:"HX-Request"`
	HXTrigger     string  `json:"HX-Trigger"`
	HXTriggerName *string `json:"HX-Trigger-Name"`
	HXTarget      string  `json:"HX-Target"`
	HXCurrentURL  string  `json:"HX-Current-URL"`
}

type hXWSHeadersWrapper struct {
	Headers HXHeaders `json:"HEADERS"`
}



var Upgrader = websocket.Upgrader{
	ReadBufferSize:  250000,
	WriteBufferSize: 250000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Printf("Error status: %d", status)
	},
}


func (s *Server) Listen(event string, listener ListenerFunc) {
	s.listeners[event] = listener
}

// Start implements hx.IServer.
func (s *Server) Mount(mountpoint string) {
	s.mux.Handle(mountpoint,http.HandlerFunc(
    func (res http.ResponseWriter,req *http.Request){
      s.handle(res,req)
    },
  ))
}

func (s *Server) newConnection(conn *websocket.Conn) Client {
	newClient := Client{
		ID: genB64(32),
		Conn:      conn,
	}

  s.clients[newClient.ID] = &newClient
  
	return newClient
}

func (s *Server) handleCloseConnection(client *Client) func(int, string) error {

	return func(code int, text string) error {
		s.OnClientDisconnect(client)
    
    delete(s.clients,client.ID)

		return nil
	}
}

func (s *Server) newMessageListener(client *Client) {
	go func(client *Client) {
		for {

      msgType,msg,err  := client.Conn.ReadMessage()

			if err != nil {
				// log.Printf("Error reading message: %v", err)
        client.Conn.Close()
        break
			}
      
      if msgType != websocket.TextMessage {
        log.Printf("Error accepting message of type %v, expected message type 1", msgType)
        continue
      }

      data := hXWSHeadersWrapper{}

      err = json.Unmarshal(msg,&data)

      if err != nil {
        log.Printf("Error reading ws headers wrapper: %v", err)
        continue
      }

			s.listeners[data.Headers.HXTrigger](client, msg)
		}
	}(client)
}


func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	log.Print("Attempting to connect new client")

	conn, err := Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Could not upgrade connection: %v", err)
	}

	new_client := s.newConnection(conn)

	conn.SetCloseHandler(s.handleCloseConnection(&new_client))

	s.OnClientConnect(&new_client)

	s.newMessageListener(&new_client)
}

func (ctx *Client) Send(message []byte) error {

	err := ctx.Conn.WriteMessage(1, message)

	if err != nil {
		return fmt.Errorf("could not send message err: %v ", err)
	}

	return nil
}

func (ctx *Client) SendStr(message string) error {
	return ctx.Send([]byte(message))
}
