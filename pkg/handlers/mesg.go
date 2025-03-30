package handlers

import (
	"log/slog"

	"github.com/baeorg/buddy/pkg/helper"
	"github.com/baeorg/buddy/pkg/storage"
	"github.com/baeorg/buddy/pkg/taskpool"
	"github.com/baeorg/buddy/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/vmihailenco/msgpack/v5"
)

type MesgSendReq struct {
	FromID  string `json:"from_id" validate:"required"`
	ConvsID uint64 `json:"convs_id" validate:"required"`
	Payload []byte `json:"payload" validate:"required"`
}

type MesgSendRes struct {
	MesgID uint64 `json:"mesg_id"`
}

func EventMesgSend(req []byte) ([]byte, error) {

	var mreq MesgSendReq
	if err := sonic.Unmarshal(req, &mreq); err != nil {
		return nil, err
	}

	err := helper.Validate.Struct(mreq)
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

	body, err := msgpack.Marshal(&mi)
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

	body, err = msgpack.Marshal(&cmi)
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
