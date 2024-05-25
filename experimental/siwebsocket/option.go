package siwebsocket

import (
	"net/http"
	"time"

	"github.com/wonksing/si/v2/sio"
)

// ClientOption is an interface with apply method.
type ClientOption interface {
	apply(c *WsClient)
}

// ClientOptionFunc wraps a function to conforms to ClientOption interface
type ClientOptionFunc func(c *WsClient)

// apply implements ClientOption's apply method.
func (o ClientOptionFunc) apply(c *WsClient) {
	o(c)
}

// WithMessageHandler sets h to c.
func WithMessageHandler(h MessageHandler) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.SetMessageHandler(h)
	})
}

func WithHub(h Hub) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.SetHub(h)
	})
}

// WithReaderOpt sets ro to c.
func WithReaderOpt(ro sio.ReaderOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.appendReaderOpt(ro)
	})
}

// WithID sets id to c's ID.
func WithID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.SetID(id)
	})
}

// WithUserID sets id to c's userID
func WithUserID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.SetUserID(id)
	})
}

// WithUserGroupID sets id to c's userGroupID
func WithUserGroupID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.SetUserGroupID(id)
	})
}

func WithWriteWait(writeWait time.Duration) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.writeWait = writeWait
	})
}
func WithReadWait(readWait time.Duration) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.readWait = readWait
		c.pingPeriod = (readWait * 9) / 10
	})
}
func WithMaxMessageSize(maxMessageSize int) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.maxMessageSize = maxMessageSize
	})
}
func WithUsePingPong(usePingPong bool) ClientOptionFunc {
	return ClientOptionFunc(func(c *WsClient) {
		c.usePingPong = usePingPong
	})
}

// HubOption is an interface with apply method.
type HubOption interface {
	apply(h *WsHub)
}

// HubOptionFunc wraps a function to conforms to HubOption interface.
type HubOptionFunc func(h *WsHub)

// apply implements HubOption's apply method.
func (o HubOptionFunc) apply(h *WsHub) {
	o(h)
}

// WithRouter sets r to h's router.
func WithRouter(r Router) HubOptionFunc {
	return HubOptionFunc(func(h *WsHub) {
		h.SetRouter(r)
	})
}

// WithHubAddr sets addr to h's hubAddr.
func WithHubAddr(addr string) HubOptionFunc {
	return HubOptionFunc(func(h *WsHub) {
		h.SetHubAddr(addr)
	})
}

// WithHubPath sets path to h's hubPath.
func WithHubPath(path string) HubOptionFunc {
	return HubOptionFunc(func(h *WsHub) {
		h.SetHubPath(path)
	})
}

// WithAfterDeleteClient sets f to h's afterDeleteClient.
func WithAfterDeleteClient(f func(Client, error)) HubOptionFunc {
	return HubOptionFunc(func(h *WsHub) {
		h.afterDeleteClient = f
	})
}

// WithAfterStoreClient sets f to h's afterStoreClient.
func WithAfterStoreClient(f func(Client, error)) HubOptionFunc {
	return HubOptionFunc(func(h *WsHub) {
		h.afterStoreClient = f
	})
}

// UpgraderOption is an interface with apply method.
type UpgraderOption interface {
	apply(u *upgraderConfig)
}

// UpgraderOptionFunc wraps a function to conforms to ClientOption interface
type UpgraderOptionFunc func(u *upgraderConfig)

// apply implements UpgraderOption's apply method.
func (o UpgraderOptionFunc) apply(u *upgraderConfig) {
	o(u)
}

func WithUpgradeHandshakeTimeout(timeout time.Duration) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.handshakeTimeout = timeout
	})
}

func WithUpgradeReadBufferSize(bufferSize int) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.readBufferSize = bufferSize
	})
}
func WithUpgradeWriteBufferSize(bufferSize int) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.writeBufferSize = bufferSize
	})
}

func WithUpgradeSubprotocols(protocols []string) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.subprotocols = protocols
	})
}

func WithUpgradeError(f func(w http.ResponseWriter, r *http.Request, status int, reason error)) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.errorFunc = f
	})
}

// WithUpgradeCheckOrigin sets f to u's CheckOrigin.
func WithUpgradeCheckOrigin(f func(r *http.Request) bool) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.checkOrigin = f
	})
}

func WithUpgradeEnableCompression(enableCompression bool) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.enableCompression = enableCompression
	})
}
