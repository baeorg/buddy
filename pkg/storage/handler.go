package storage

import (
	"fmt"
	"log/slog"

	"github.com/baeorg/buddy/pkg/types"
	"github.com/sunvim/gmdbx"
)

func UserTokenUpdateHandler(mi *types.MesgInfo, wtx *gmdbx.Tx) error {
	keys := uint64(mi.Key.(float64))
	key := gmdbx.U64(&keys)
	msgVal := mi.Content.(string)
	val := gmdbx.String(&msgVal)
	gerr := wtx.Put(dbHandler.permits, &key, &val, gmdbx.PutUpsert)
	if gerr != gmdbx.ErrSuccess {
		slog.Error("put user token failed: ", "err", gerr)
		return fmt.Errorf("failed to update user token: %v", gerr)
	}
	slog.Info("user token updated", "user id", keys, "token", mi.Content)
	return nil
}

func ConvsCreateHandler(mi *types.MesgInfo, wtx *gmdbx.Tx) error {
	slog.Info("conversation created", "mi", mi)

	convID := uint64(mi.Key.(float64))
	conv := gmdbx.U64(&convID)

	msgVal := mi.Content.(string)
	val := gmdbx.String(&msgVal)

	gerr := wtx.Put(dbHandler.convsUsers, &conv, &val, gmdbx.PutUpsert)
	if gerr != gmdbx.ErrSuccess {
		slog.Error("put conversation failed: ", "err", gerr)
		return fmt.Errorf("failed to create conversation: %v", gerr)
	}
	slog.Info("conversation created", "conversation id", convID)
	return nil
}

func MesgSendHandler(mi *types.MesgInfo, wtx *gmdbx.Tx) error {
	msgID := uint64(mi.Key.(float64))
	msg := gmdbx.U64(&msgID)

	msgVal := mi.Content.([]byte)
	val := gmdbx.Bytes(&msgVal)
	gerr := wtx.Put(dbHandler.mesgs, &msg, &val, gmdbx.PutUpsert)
	if gerr != gmdbx.ErrSuccess {
		slog.Error("put message failed: ", "err", gerr)
		return fmt.Errorf("failed to send message: %v", gerr)
	}
	slog.Info("message sent", "message id", msgID)
	return nil
}

func MesgConvsSaveHandler(mi *types.MesgInfo, wtx *gmdbx.Tx) error {
	msgID := uint64(mi.Key.(float64))
	msg := gmdbx.U64(&msgID)
	msgVal := mi.Content.(uint64)
	val := gmdbx.U64(&msgVal)
	gerr := wtx.Put(dbHandler.convsMesgs, &msg, &val, gmdbx.PutUpsert)
	if gerr != gmdbx.ErrSuccess {
		slog.Error("put message failed: ", "err", gerr)
		return fmt.Errorf("failed to send message: %v", gerr)
	}
	slog.Info("message sent", "message id", msgID)
	return nil
}
