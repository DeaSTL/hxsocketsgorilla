package hx

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	Request     string `json:"HX-Request"`
	Trigger     string `json:"HX-Trigger"`
	TriggerName string `json:"HX-Trigger-Name"`
	Target      string `json:"HX-Target"`
	CurrentURL  string `json:"HX-Current-URL"`
	Includes    map[string]any
}


type IServer interface { 
  Start(mountpoint string);
  //Listen(event string, listener func(*IClient, *Message))
  Broadcast(event string, message []byte) error
}

type IClient interface {
  Send(message []byte) error
  SendStr(message string) error
}

// type Server struct {
// 	Connections  *[]*ConnectionCtx
// 	listeners    map[string]Listener
// 	OnConnection func(*ConnectionCtx)
// 	OnDisconnect func(*ConnectionCtx)
// }

type ConnectionCtx struct {
	Client    *websocket.Conn
	SessionID string
}

// func (ss *Server) LogConnections() {
// 	for _, client := range *ss.Connections {
// 		log.Printf("Client %v", client.SessionID)
// 	}
// }

// func (ss *Server) New() {
// 	ss.OnConnection = func(ctx *ConnectionCtx) {}
// 	ss.OnDisconnect = func(ctx *ConnectionCtx) {}
// 	ss.listeners = map[string]Listener{}
// 	ss.Connections = &[]*ConnectionCtx{}
// }

// func (ss *IServer) Start(endpoint string) {
//   switch ss.(type) {
//     case compat.FiberServer:
//       ss.Start(endpoint) 
//     case compat.NethttpServer:
//       ss.Start(endpoint)
//   }
// 	mux.Handle(endpoint, http.HandlerFunc(ss.handleNewConnection))
// }

// func (ss *Server) handleCloseConnection(ctx *IClient) func(int, string) error {
// 
// 	return func(code int, text string) error {
// 		ss.OnDisconnect(ctx)
// 		var index int
// 
// 		for c_index, conn := range *ss.Connections {
// 			if conn.SessionID == ctx.SessionID {
// 				index = c_index
// 			}
// 		}
// 
// 		*ss.Connections = append((*ss.Connections)[:index], (*ss.Connections)[index+1:]...)
// 
// 		return nil
// 	}
// }
// func (ss *Server) handleNewConnection(w http.ResponseWriter, r *http.Request) {
// 	log.Print("Attempting to connect new client")
// 
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 
// 	if err != nil {
// 		log.Printf("Could not upgrade connection: %v", err)
// 	}
// 
// 	new_ctx := ss.newConnection(conn)
// 
// 	//conn.SetCloseHandler(ss.handleCloseConnection(&new_ctx))
// 
// 	ss.OnConnection(&new_ctx)
// 
// 	ss.newMessageListener(&new_ctx)
// }

/*
	{
	  HEADERS: {}, // HTMX stuff
	  some_key: 'some value'
	  ...
	}
*/
// func (ss *Server) newMessageListener(ctx *ConnectionCtx) {
// 	go func(ctx *ConnectionCtx) {
// 		for {
// 			rawMessage := map[string]any{}
// 
// 			err := ctx.Client.ReadJSON(&rawMessage)
// 
// 			if err != nil {
// 				log.Printf("Error reading json message: %v", err)
// 				break
// 			}
// 
// 			/*
// 			   We're pulling off the HEADERS and making it a HtmxMesssage
// 			*/
// 			jsonMessage, err := json.Marshal(rawMessage["HEADERS"])
// 
// 			if err != nil {
// 
// 			}
// 
// 			message := Message{}
// 
// 			err = json.Unmarshal(jsonMessage, &message)
// 
// 			// Yoinkin' that header off the json blob
// 			delete(rawMessage, "HEADERS")
// 
// 			message.Includes = rawMessage
// 
// 			err = ss.messageHandler(ctx, &message)
// 
// 			if err != nil {
// 				log.Printf("Error handling message: %v ", err)
// 				continue
// 			}
// 		}
// 		defer ctx.Client.Close()
// 	}(ctx)
// }

// func (ss *Server) messageHandler(ctx *ConnectionCtx, message *Message) error {
// 	for event, listener := range ss.listeners {
// 		if message.Trigger == event {
// 
// 			listener.Callback(ctx, message)
// 		}
// 	}
// 	return nil
// }

// func (ss *Server) newConnection(conn *websocket.Conn) ConnectionCtx {
// 	new_ctx := ConnectionCtx{
// 		SessionID: GenB64(32),
// 		Client:    conn,
// 	}
// 
// 	*ss.Connections = append(*ss.Connections, &new_ctx)
// 
// 	return new_ctx
// }

// func (ss *Server) Broadcast(event string, message []byte) error {
// 	for _, conn := range *ss.Connections {
// 		err := conn.Send(message)
// 
// 		if err != nil {
// 			log.Printf("Error broadcasting: %v", err)
// 		}
// 
// 		log.Printf("Broadcasting message: %v to: %v", message, conn.SessionID)
// 	}
// 	return nil
// }
// 
// func (ss *Server) Listen(event string, listener func(*IClient, *Message)) {
// 	ss.listeners[event] = Listener{Callback: listener}
// }

// func (ctx *ConnectionCtx) Send(message []byte) error {
// 
// 	err := ctx.Client.WriteMessage(1, message)
// 
// 	if err != nil {
// 		return fmt.Errorf("Could not send message err: %v ", err)
// 	}
// 
// 	return nil
// }
// 
// func (ctx *ConnectionCtx) SendStr(message string) error {
// 	return ctx.Send([]byte(message))
// }
// 
// func (ss *IServer) SendFilter(event string, message []byte, check func(*ConnectionCtx) bool) {
// 	for _, conn := range *ss.Connections {
// 		if check(conn) {
// 			err := conn.Send(message)
// 
// 			if err != nil {
// 				log.Printf("Error broadcasting: %v", err)
// 			}
// 
// 			log.Printf("Sending filtered message: %v to: %v", message, conn.SessionID)
// 		}
// 	}
// }
