package storage

import (
	"log/slog"

	"github.com/baeorg/buddy/pkg/types"
	"github.com/sunvim/gmdbx"
	"github.com/vmihailenco/msgpack/v5"
)

func UpdatePermission(userID string, token string) error {
	mi := &types.MesgInfo{
		MsgType: types.UserTokenUpdate,
		Key:     PermiPrefix + userID,
		Content: []byte(token),
	}
	ms, err := msgpack.Marshal(mi)
	if err != nil {
		return err
	}
	return dbHandler.Put(ms)
}

func IsPermit(userID string, token string) bool {

	rtx, err := dbHandler.Rtx()
	if err != nil {
		slog.Warn("failed to get permission for user ", "user id", userID, "err", err)
		return false
	}
	defer rtx.Commit()

	tokenKey := PermiPrefix + userID
	key := gmdbx.String(&tokenKey)
	val := gmdbx.Val{}
	xerr := rtx.Get(dbHandler.genv, &key, &val)
	if xerr != gmdbx.ErrSuccess {
		slog.Error("failed to get permission for user ", "user id", tokenKey, "err", xerr)
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
