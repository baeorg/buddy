package storage

import (
	"fmt"
	"log/slog"

	"github.com/baeorg/buddy/pkg/types"
	"github.com/sunvim/gmdbx"
)

func UserTokenUpdateHandler(mi *types.MesgInfo, wtx *gmdbx.Tx, dbi gmdbx.DBI) error {
	keys := mi.Key.(string)
	key := gmdbx.String(&keys)
	val := gmdbx.Bytes(&mi.Content)
	gerr := wtx.Put(dbi, &key, &val, gmdbx.PutUpsert)
	if gerr != gmdbx.ErrSuccess {
		slog.Error("put user token failed: ", "err", gerr)
		return fmt.Errorf("failed to update user token: %v", gerr)
	}
	slog.Info("user token updated", "user id", keys, "token", mi.Content)
	return nil
}

func ConvsCreateHandler(mi *types.MesgInfo, wtx *gmdbx.Tx, dbi gmdbx.DBI) error {
	convID := mi.Key.(uint64)
	conv := gmdbx.U64(&convID)
	val := gmdbx.Bytes(&mi.Content)
	gerr := wtx.Put(dbi, &conv, &val, gmdbx.PutUpsert)
	if gerr != gmdbx.ErrSuccess {
		slog.Error("put conversation failed: ", "err", gerr)
		return fmt.Errorf("failed to create conversation: %v", gerr)
	}
	slog.Info("conversation created", "conversation id", convID)
	return nil
}
