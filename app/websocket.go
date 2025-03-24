package app

import (
	"log/slog"
	"net/http"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = newUpgrader()

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.BlockingModAsyncWrite = true

	u.OnOpen(func(c *websocket.Conn) {
		slog.Info("OnOpen", "remoteAddr", c.RemoteAddr().String())
	})
	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		c.WriteMessage(messageType, data)
	})
	u.OnClose(func(c *websocket.Conn, err error) {
		slog.Info("OnClose", "remoteAddr", c.RemoteAddr().String(), "error", err)
	})
	return u
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	// TODO: need to check the user permission
	//
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade websocket connection", "error", err)
	}
}
