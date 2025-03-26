package app

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	OnLineUsers = cmap.New[*User]()
)

type User struct {
	ID   uint64
	conn *websocket.Conn
}
