// Copyright 2022. Motty Cohen
//
// Default web socket client registry
//
package socket

import (
	"sync"
)

// DefaultClientRegistry is the basic implementation of web socket client registry
// It is used by web socket server to manage and track all the connected web socket clients
type DefaultClientRegistry struct {
	sync.RWMutex
	Connections map[string]IWSClient
}

// Initialize registry
func (r *DefaultClientRegistry) Start() {}

// Register new connected client
func (r *DefaultClientRegistry) RegisterClient(wsc IWSClient) {
	r.Lock()
	defer r.Unlock()
	r.Connections[wsc.ID()] = wsc
}

// Unregister disconnected client
func (r *DefaultClientRegistry) UnregisterClient(wsc IWSClient) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Connections[wsc.ID()]; ok {
		delete(r.Connections, wsc.ID())
	}
}

// Get number of current connected clients
func (r *DefaultClientRegistry) ConnectedClients() int {
	r.Lock()
	defer r.Unlock()
	return len(r.Connections)
}
