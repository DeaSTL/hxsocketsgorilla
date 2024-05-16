package compat

import "github.com/deastl/hx-sockets"

type FiberServer struct {
}

// Broadcast implements hx.IServer.
func (f FiberServer) Broadcast(event string, message []byte) error {
	panic("unimplemented")
}

// Start implements hx.IServer.
func (f FiberServer) Start(mountpoint string) {
	panic("unimplemented")
}

type FiberClient struct {
}

func NewFiberClient() hx.IServer {
	return FiberServer{}
}
