package storage

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/sunvim/gmdbx"
	"github.com/sunvim/mq"
)

var (
	dbHandler *DB
)

func InitDB(ctx context.Context) {
	dbHandler = New(ctx, viper.GetString("db.path"), viper.GetString("mq.path"))
}

func (d *DB) Put(data []byte) error {
	return d.mesgeq.Push(&mq.Message{Data: data})
}

func GetNextConvID() uint64 {
	return dbHandler.seqm[SeqmConvs].Add(1)
}

func GetNextMesgID() uint64 {
	return dbHandler.seqm[SeqmMesgs].Add(1)
}

func SaveDataIntoDB(data []byte) error {
	return dbHandler.Put(data)
}

func GetUsersByConvID(convID uint64) ([]uint64, error) {
	rtx, err := dbHandler.Rtx()
	if err != nil {
		return nil, err
	}
	defer rtx.Commit()

	cur, xerr := rtx.OpenCursor(dbHandler.convsUsers)
	if xerr != gmdbx.ErrSuccess {
		return nil, fmt.Errorf("failed to open cursor: %v", xerr)
	}
	defer cur.Close()
	users := make([]uint64, 0)

	key := gmdbx.U64(&convID)
	val := gmdbx.Val{}
	for cur.Get(&key, &val, gmdbx.CursorNextDup) == gmdbx.ErrSuccess {
		users = append(users, val.U64())
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no users found for convID %d", convID)
	}

	return users, nil
}
