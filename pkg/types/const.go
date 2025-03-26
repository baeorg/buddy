package types

type KindMesg uint64

const (
	KindConvsCreate = 8000 + iota
	KindConvsUpdate
	KindConvsDelete
	KindConvsGet

	KindMesgCreate
	KindMesgUpdate
	KindMesgRemove
	KindMesgGet
)
