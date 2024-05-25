package siwebsocket

import (
	"context"
	"errors"
	"sync"
	"time"
)

// WsHub maintains the set of active clients and broadcasts messages to the
// clients.
type WsHub struct {
	// Connected clients.
	clients sync.Map

	// clients information and their location storage
	router Router

	// channel to broadcast message to connected clients
	broadcast chan []byte
	// used to wait until broadcasting completely finished.
	broadcastWait chan struct{}

	// channel to add clients to `clients` map
	register chan Client

	// channel to remove clients from `clients` map
	unregister chan Client

	// once runDone is closed, a client cannot be received from register or unregister channel.
	runDone chan struct{}
	// used to wait until receiving clients from register or unregister channel is completely finished.
	runWait chan struct{}
	// once stopClient was received, hub starts closing.
	stopClient chan string
	// once clientDone is closed, it prevents from sending into register, unregister and broadcast channel.
	clientDone chan struct{}
	// used to wait until the hub is completely closed.
	terminated chan struct{}

	// Time allowed to write a message to the peer.
	writeWait time.Duration
	// Time allowed to read the next pong message from the peer.
	readWait time.Duration
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod time.Duration
	// Maximum message size allowed from peer.
	maxMessageSize int
	// use ping/pong
	usePingPong bool

	// hub's address. eg. http://127.0.0.1:8080
	hubAddr string
	// hub's path. eg. /send/message/to/client_id
	hubPath string

	// handlers
	// called after deleting c from clients map.
	afterDeleteClient func(c Client, err error)
	// called after storing c into clients map.
	afterStoreClient func(c Client, err error)
}

// NewWsHub creates a hub
func NewWsHub(hubAddr, hubPath string, writeWait time.Duration, readWait time.Duration,
	maxMessageSize int, usePingPong bool, opts ...HubOption) *WsHub {

	pingPeriod := (readWait * 9) / 10

	h := &WsHub{
		broadcast:     make(chan []byte, 1024),
		broadcastWait: make(chan struct{}, 1),
		register:      make(chan Client),
		unregister:    make(chan Client),

		clients:    sync.Map{},
		runDone:    make(chan struct{}),
		runWait:    make(chan struct{}, 1),
		stopClient: make(chan string, 1),
		clientDone: make(chan struct{}),
		terminated: make(chan struct{}),
		router:     &NopRouter{},

		hubAddr:        hubAddr,
		hubPath:        hubPath,
		writeWait:      writeWait,
		readWait:       readWait,
		pingPeriod:     pingPeriod,
		maxMessageSize: maxMessageSize,
		usePingPong:    usePingPong,
	}

	for _, o := range opts {
		o.apply(h)
	}

	if h.afterDeleteClient == nil {
		h.afterDeleteClient = func(c Client, err error) {
			// nothing
		}
	}

	if h.afterStoreClient == nil {
		h.afterStoreClient = func(c Client, err error) {
			// nothing
		}
	}

	return h
}

var ErrClientNotExist = errors.New("client does not exist")

// runBroadcast starts receiving messages to broadcast to clients.
func (h *WsHub) runBroadcast() {
	defer close(h.broadcastWait)
	for message := range h.broadcast {
		h.clients.Range(func(key interface{}, value interface{}) bool {
			value.(Client).Send(message)
			return true
		})
	}
}

// stopAdnWaitBroadcast stops runBroadcast and wait for it to return.
func (h *WsHub) stopAndWaitBroadcast() {
	close(h.broadcast)
	<-h.broadcastWait
}

// runClient starts receiving new and disconnected clients.
func (h *WsHub) runClient() {
	defer close(h.runWait)
	for {
		select {
		case <-h.runDone:
			return
		case client := <-h.register:
			loadedClient, exist := h.clients.LoadOrStore(client.GetID(), client)
			if exist {
				// Stop will lead to removeClient method called.
				// Do not call removeClient method here.
				loadedClient.(Client).Stop()
				h.clients.Store(client.GetID(), client)
			}
			h.afterStoreClient(client, nil)
		case client := <-h.unregister:
			// stopped clients with connection closed are received here.
			// remove them from `clients` map
			_, exist := h.clients.LoadAndDelete(client.GetID())
			if !exist {
				h.afterDeleteClient(client, ErrClientNotExist)
			} else {
				h.afterDeleteClient(client, nil)
			}
		}
	}
}

// Run starts the hub.
func (h *WsHub) Run() {
	go h.waitStop()
	go h.runBroadcast()

	h.runClient()
}

// stopAndWaitClient stops runClient method and wait for it to return
func (h *WsHub) stopAndWaitClient() {
	close(h.runDone) // stops Run method, closing this channel doesn't mean the select loop returns
	<-h.runWait      // wait Run method to return
}

func (h *WsHub) waitStop() {
	<-h.stopClient      // wait until Stop method is called
	close(h.clientDone) // prevent from sending into register/unregister/broadcast channel
	h.stopAndWaitClient()
	h.stopAndWaitBroadcast()
	h.removeAllClients() // stops and closes all clients and remove from clients map
	close(h.terminated)
}

// Stop initiates the hub to stop.
func (h *WsHub) Stop() error {
	select {
	case h.stopClient <- "stop":
	default:
		return ErrStopChannelFull
	}
	return nil
}

// Wait waits until h is completely finished.
func (h *WsHub) Wait() {
	<-h.terminated
}

var ErrHubClosed = errors.New("hub is closed")

func (h *WsHub) Add(client Client) error {
	select {
	case <-h.clientDone:
		return ErrHubClosed
	default:
	}

	select {
	case <-h.clientDone:
		h.afterStoreClient(client, ErrHubClosed)
		return ErrHubClosed
	case h.register <- client:
		return h.router.Store(context.Background(), client.GetID(), client.GetUserID(),
			client.GetUserGroupID(), h.hubAddr, h.hubPath)
	}
}

func (h *WsHub) Remove(client Client) error {
	select {
	case <-h.clientDone:
		return ErrHubClosed
	default:
	}

	select {
	case <-h.clientDone:
		h.afterDeleteClient(client, ErrHubClosed)
		return ErrHubClosed
	case h.unregister <- client:
		return h.router.Delete(context.Background(), client.GetID())
	}
}

func (h *WsHub) Broadcast(message []byte) error {
	select {
	case <-h.clientDone:
		return ErrHubClosed
	default:
	}

	select {
	case <-h.clientDone:
		return ErrHubClosed
	case h.broadcast <- message:
	}

	return nil
}

func (h *WsHub) removeAllClients() error {

	// stops gracefully
	h.clients.Range(func(key interface{}, value interface{}) bool {
		value.(Client).Stop()
		value.(Client).Wait()
		h.clients.Delete(value.(Client).GetID())
		h.router.Delete(context.Background(), value.(Client).GetID())
		return true
	})

	return nil
}

func (h *WsHub) RemoveRandomClient() error {
	h.clients.Range(func(key interface{}, value interface{}) bool {
		value.(Client).Stop()
		return false
	})

	return nil
}

func (h *WsHub) LenClients() int {
	lenClients := 0
	h.clients.Range(func(key interface{}, value interface{}) bool {
		lenClients++
		return true
	})

	return lenClients
}

func (h *WsHub) SendMessage(id string, msg []byte) error {
	if c, ok := h.clients.Load(id); !ok {
		return errors.New("client not found, " + id)
	} else {
		err := c.(Client).Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendMessageWithIDAndUserGroupID sends msg to a client with id and userGroupID.
func (h *WsHub) SendMessageWithIDAndUserGroupID(id, userGroupID string, msg []byte) error {
	if c, ok := h.clients.Load(id); !ok {
		return errors.New("client not found, " + id)
	} else {
		if c.(Client).GetUserGroupID() != userGroupID {
			return errors.New("client with, " + userGroupID + ", not found")
		}
		err := c.(Client).Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendMessageAndWait send msg to a client with id and wait until the msg was
// successfully written to the client.
func (h *WsHub) SendMessageAndWait(id string, msg []byte) error {
	if c, ok := h.clients.Load(id); !ok {
		return errors.New("client not found, " + id)
	} else {
		err := c.(Client).SendAndWait(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *WsHub) SetRouter(r Router) {
	h.router = r
}
func (h *WsHub) SetHubAddr(hubAddr string) {
	h.hubAddr = hubAddr
}
func (h *WsHub) SetHubPath(hubPath string) {
	h.hubPath = hubPath
}

func (h *WsHub) GetWriteWait() time.Duration {
	return h.writeWait
}

func (h *WsHub) GetReadWait() time.Duration {
	return h.readWait
}

func (h *WsHub) GetUsePingPong() bool {
	return h.usePingPong
}

func (h *WsHub) GetMaxMessageSize() int {
	return h.maxMessageSize
}
