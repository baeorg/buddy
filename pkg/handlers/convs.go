package handlers

import (
	"log/slog"

	"github.com/baeorg/buddy/pkg/share"
	"github.com/baeorg/buddy/pkg/storage"
	"github.com/baeorg/buddy/pkg/types"
	"github.com/bytedance/sonic"
)

type ConvsReq struct {
	Title   string   `json:"title" validate:"required"`
	UserIDs []uint64 `json:"user_ids" validate:"required"`
}

type ConvsRes struct {
	ID uint64 `json:"id"`
}

func ConvsCreate(req []byte) (rsp []byte, err error) {
	var (
		reqt ConvsReq
	)

	err = sonic.Unmarshal(req, &reqt)
	if err != nil {
		slog.Error("failed to convert request to ConvsReq")
		return nil, share.ErrInvalidRequest
	}

	if len(reqt.UserIDs) < 2 {
		slog.Error("at least two users are required")
		return nil, share.ErrAtLeastTwoUsers
	}

	convsID := storage.GetNextConvID()

	mi := types.MesgInfo{
		MsgType: types.EventConvCreate,
		Key:     convsID,
		Content: req,
	}

	body, err := sonic.Marshal(&mi)
	if err != nil {
		slog.Error("failed to marshal message info", "error", err)
		return nil, err
	}

	storage.SaveDataIntoDB(body)

	res := ConvsRes{
		ID: convsID,
	}

	rsp, err = sonic.Marshal(res)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		return nil, err
	}

	return rsp, nil
}

func ConvsUpdate(req []byte) (rsp []byte, err error) {
	return nil, nil
}

func ConvsDelete(req []byte) (rsp []byte, err error) {
	return nil, nil
}

func ConvsGet(req []byte) (rsp []byte, err error) {
	return nil, nil
}
