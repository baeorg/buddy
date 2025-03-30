package app

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/baeorg/buddy/pkg/handlers"
	"github.com/baeorg/buddy/pkg/share"
	"github.com/baeorg/buddy/pkg/storage"
	"github.com/baeorg/buddy/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = newUpgrader()

type WsReq struct {
	Kind types.Kind `json:"kind" validate:"required"`
	Reqs []byte     `json:"reqs" validate:"required"`
}

type WsRsp struct {
	Kind types.Kind `json:"kind"`
	Code uint64     `json:"code"`
	Mesg string     `json:"mesg"`
	Rsp  []byte     `json:"rsp"`
}

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.BlockingModAsyncWrite = true

	u.OnOpen(func(c *websocket.Conn) {
		slog.Info("OnOpen", "remoteAddr", c.RemoteAddr().String())
	})

	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		var (
			wsreq WsReq
			wsrsp = &WsRsp{
				Code: share.SuccessCode,
				Mesg: share.Success.Error(),
			}
		)
		if err := sonic.Unmarshal(data, &wsreq); err != nil {
			slog.Error("Failed to unmarshal WsReq", "error", err)
			wsrsp.Code = share.ErrInvalidRequestCode
			wsrsp.Mesg = share.ErrInvalidRequest.Error()
			goto END
		}

		if handler, ok := handlers.HandlerMap[wsreq.Kind]; ok {
			rs, err := handler(wsreq.Reqs)
			if err != nil {
				slog.Error("Failed to handle request", "error", err)
				wsrsp.Code = share.ErrInternalErrorCode
				wsrsp.Mesg = share.ErrInternalError.Error()
				goto END
			}
			wsrsp.Kind = wsreq.Kind
			wsrsp.Rsp = rs
		} else {
			slog.Error("Handler not found", "kind", wsreq.Kind)
			wsrsp.Code = share.ErrNotSupportedCode
			wsrsp.Mesg = share.ErrNotSupported.Error()
		}
	END:
		rsp, _ := sonic.Marshal(wsrsp)
		c.WriteMessage(messageType, rsp)
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

	handlers.OnLineUsers.Set(userid, &handlers.User{
		ID:   uid,
		Conn: conn,
	})
}
