package app

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/baeorg/buddy/pkg/storage"
	"github.com/baeorg/buddy/pkg/types"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = newUpgrader()

type WsReq struct {
	Kind types.KindMesg `json:"kind" validate:"required"`
	Reqs []byte         `json:"reqs" validate:"required"`
}

type WsRsp struct {
	Kind types.KindMesg `json:"kind"`
	Rsp  []byte         `json:"rsp"`
}

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

const (
	XUserID    = "X-User-ID"
	XUserToken = "X-User-Token"
)

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	// check user permission
	userid := r.Header.Get(XUserID)
	userToken := r.Header.Get(XUserToken)
	if userid == "" || userToken == "" {
		http.Error(w, "Unauthorized: user id or user token should be provided", http.StatusUnauthorized)
		return
	}
	slog.Info("User connected", "userid", userid, "token", userToken, "remoteAddr", r.RemoteAddr)

	if !storage.IsPermit(userid, userToken) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade websocket connection", "error", err)
	}

	uid, err := strconv.ParseUint(userid, 10, 64)
	if err != nil {
		slog.Error("Failed to parse user id", "error", err)
		http.Error(w, "Failed to parse user id", http.StatusBadRequest)
		return
	}

	OnLineUsers.Set(userid, &User{
		ID:   uid,
		conn: conn,
	})
}
