package handlers

import (
	"log/slog"
	"strconv"

	"github.com/baeorg/buddy/pkg/helper"
	"github.com/baeorg/buddy/pkg/share"
	"github.com/baeorg/buddy/pkg/storage"
	"github.com/baeorg/buddy/pkg/taskpool"
	"github.com/baeorg/buddy/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

type MesgSendReq struct {
	FromID  string `json:"from_id" validate:"required"`
	ConvsID uint64 `json:"convs_id" validate:"required"`
	Payload []byte `json:"payload" validate:"required"`
}

type MesgSendRes struct {
	MesgID uint64 `json:"mesg_id"`
}

type MesgOnLine struct {
	MesgID uint64 `json:"mesg_id"`
	Body   []byte `json:"body"`
}

func EventMesgSend(req []byte) ([]byte, error) {

	var mreq MesgSendReq

	err := sonic.Unmarshal(req, &mreq)
	if err != nil {
		slog.Error("failed to unmarshal request", "error", err)
		return nil, share.ErrInvalidRequest
	}

	err = helper.Validate.Struct(mreq)
	if err != nil {
		return nil, err
	}

	mesgID := storage.GetNextMesgID()

	// save message into database
	mi := &types.MesgInfo{
		MsgType: types.EventMesgSend,
		Key:     mesgID,
		Content: req,
	}

	body, err := sonic.Marshal(&mi)
	if err != nil {
		slog.Error("failed to marshal message info", "error", err)
		return nil, err
	}
	storage.SaveDataIntoDB(body)

	// save message to conversation
	cmi := &types.MesgInfo{
		MsgType: types.EventConvAddMesg,
		Key:     mreq.ConvsID,
		Content: mesgID,
	}

	body, err = sonic.Marshal(&cmi)
	if err != nil {
		slog.Error("failed to marshal message info", "error", err)
		return nil, err
	}
	storage.SaveDataIntoDB(body)

	res, err := sonic.Marshal(&MesgSendRes{
		MesgID: mesgID,
	})

	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		return nil, err
	}

	// send message to users
	taskpool.Workers.AddTask(func() (any, error) {
		fromid, err := strconv.ParseUint(mreq.FromID, 10, 64)
		if err != nil {
			slog.Error("failed to parse from_id", "error", err)
			return nil, err
		}

		users, err := storage.GetUsersByConvID(mreq.ConvsID)
		if err != nil {
			slog.Error("failed to get users by convID", "error", err)
			return nil, err
		}

		mo := &MesgOnLine{
			MesgID: mesgID,
			Body:   req,
		}

		mos, err := sonic.Marshal(&mo)
		if err != nil {
			slog.Error("failed to marshal message", "error", err)
			return nil, err
		}

		count := 3

		for _, user := range users {
			if user == fromid {
				continue
			}
			toUser := strconv.FormatUint(user, 10)
			u, ok := OnLineUsers.Get(toUser)
			if !ok {
				continue
			}

			// if send message failed, retry 3 times
			count = 3
		REPEAT:
			err = u.Conn.WriteMessage(websocket.TextMessage, mos)
			if err != nil {
				count -= 1
				if count <= 0 {
					continue
				}
				goto REPEAT
			}
		}

		return nil, nil
	})

	return res, nil
}

func EventMesgEdit(req []byte) ([]byte, error) {
	return nil, nil
}

func EventMesgRemove(req []byte) ([]byte, error) {
	return nil, nil
}

func EventMesgGet(req []byte) ([]byte, error) {
	return nil, nil
}
