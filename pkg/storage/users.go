package storage

import (
	"log/slog"

	"github.com/baeorg/buddy/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/sunvim/gmdbx"
)

func UpdatePermission(userID uint64, token string) error {
	mi := &types.MesgInfo{
		MsgType: types.EventTokenSet,
		Key:     userID,
		Content: token,
	}
	ms, err := sonic.Marshal(mi)
	if err != nil {
		return err
	}
	return dbHandler.Put(ms)
}

func IsPermit(userID uint64, token string) bool {

	rtx, err := dbHandler.Rtx()
	if err != nil {
		slog.Warn("failed to get permission for user ", "user id", userID, "err", err)
		return false
	}
	defer rtx.Commit()

	key := gmdbx.U64(&userID)
	val := gmdbx.Val{}
	xerr := rtx.Get(dbHandler.permits, &key, &val)
	if xerr != gmdbx.ErrSuccess {
		slog.Error("failed to get permission for user ", "user id", key, "err", xerr)
		return false
	}
	if token == val.String() {
		return true
	}

	return false
}

func CreateUser(user *types.User) error {
	return nil
}

func DeleteUser(userID string) error {
	return nil
}

func UpdateUser(user *types.User) error {
	return nil
}

func GetUser(userID string) (*types.User, error) {
	return nil, nil
}
