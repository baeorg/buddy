package types

type MesgType uint64

const (
	EventTokenSet MesgType = 6000 + iota
	EventConvCreate
	EventConvJoin
	EventConvLeave
	EventConvList
	EventConvAddMesg
	EventMesgSend
	EventMesgEdit
	EventMesgRemove
	EventMesgGet
)

type MesgInfo struct {
	MsgType MesgType
	Key     any
	Content any
}
