package storage

import (
	"context"

	"github.com/spf13/viper"
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
