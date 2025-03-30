package handlers

import "github.com/baeorg/buddy/pkg/types"

type BusiHandler func(req []byte) (rsp []byte, err error)

var (
	HandlerMap = map[types.Kind]BusiHandler{
		types.KindConvsCreate: ConvsCreate,
		types.KindConvsUpdate: ConvsUpdate,
		types.KindConvsGet:    ConvsGet,
		types.KindConvsDelete: ConvsDelete,
		types.KindMesgCreate:  EventMesgSend,
		types.KindMesgUpdate:  EventMesgEdit,
		types.KindMesgRemove:  EventMesgRemove,
		types.KindMesgGet:     EventMesgGet,
	}
)
