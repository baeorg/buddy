package app

import (
	"io"
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

		slog.Info("Received WsReq", "kind", wsreq.Kind, "req", wsreq.Reqs)

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

type Request struct {
	ID    uint64 `json:"id"`
	Token string `json:"token"`
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("read request body failed:", "err", err)
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	var req Request
	err = sonic.Unmarshal(body, &req)
	if err != nil {
		slog.Error("parse body failed:", "err", err)
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	// check user permission
	userid := strconv.FormatUint(req.ID, 10)
	userToken := req.Token
	if userid == "" || userToken == "" {
		http.Error(w, "Unauthorized: user id or user token should be provided", http.StatusUnauthorized)
		return
	}
	slog.Info("User connected", "userid", userid, "token", userToken, "remoteAddr", r.RemoteAddr)

	if !storage.IsPermit(req.ID, userToken) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade websocket connection", "error", err)
	}

	handlers.OnLineUsers.Set(userid, &handlers.User{
		ID:   req.ID,
		Conn: conn,
	})
}
