# hx-sockets

hx-sockets is a small layer around the gorilla websockets *currently making it easier to get up and running with htmx websockets.
currently it is not ready for production use or anything outside of demoing because i'd like to make sure that the compatibility module system is the best that it can be.

here is some demo code
```go
func main() {
	mux := http.NewServeMux()
	server := compat.NewNetHttp(mux).(compat.NethttpServer)
	
	server.Listen("some_message", func(ctx *compat.NethttpClient, msg *hx.Message) {
	
		ctx.SendStr(`<a id="some_message">some message</a>`)
	})

	server.Start("/ws") // where the web socket mounts
  
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    	// send back index page
	})
  	http.ListenAndServe(":3000",mux)
}
```

```html
  // load in htmx and htmx-ws....
  <body hx-ext="ws" ws-connect="/ws">
    <div>
      <button id="some_message" hx-trigger="click" ws-send></button>
    </div>
  </body>
```

[Demo](https://github.com/DeaSTL/hx-sockets-demo)
