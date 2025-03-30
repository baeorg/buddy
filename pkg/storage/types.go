package storage

import (
	"github.com/baeorg/buddy/pkg/types"
	"github.com/sunvim/gmdbx"
)

const (
	PermiPrefix string = "tokens:"
)

type MesgHandler func(mi *types.MesgInfo, wtx *gmdbx.Tx) error

var (
	mesgHandlers = map[types.MesgType]MesgHandler{
		types.EventTokenSet:   UserTokenUpdateHandler,
		types.EventConvCreate: ConvsCreateHandler,
		types.EventMesgSend:   MesgSendHandler,
	}
)
